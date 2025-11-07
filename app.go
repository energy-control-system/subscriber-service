package main

import (
	"context"
	"fmt"
	"io/fs"
	"subscriber-service/api"
	"subscriber-service/cluster/task"
	"subscriber-service/config"
	dbcontract "subscriber-service/database/contract"
	dbobject "subscriber-service/database/object"
	dbsubscriber "subscriber-service/database/subscriber"
	"subscriber-service/service/contract"
	"subscriber-service/service/object"
	"subscriber-service/service/registry"
	"subscriber-service/service/subscriber"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/sunshineOfficial/golib/db"
	"github.com/sunshineOfficial/golib/gohttp"
	"github.com/sunshineOfficial/golib/gohttp/goserver"
	"github.com/sunshineOfficial/golib/gokafka"
	"github.com/sunshineOfficial/golib/golog"
)

const (
	serviceName = "subscriber-service"
	dbTimeout   = 15 * time.Second
)

type App struct {
	/* main */
	mainCtx  context.Context
	log      golog.Logger
	settings config.Settings

	/* http */
	server goserver.Server

	/* db */
	postgres           *sqlx.DB
	kafka              gokafka.Kafka
	inspectionConsumer gokafka.Consumer

	/* services */
	subscriberService *subscriber.Service
	objectService     *object.Service
	contractService   *contract.Service
	registryService   *registry.Service
}

func NewApp(mainCtx context.Context, log golog.Logger, settings config.Settings) *App {
	return &App{
		mainCtx:  mainCtx,
		log:      log,
		settings: settings,
	}
}

func (a *App) InitDatabases(fs fs.FS, path string) (err error) {
	postgresCtx, cancelPostgresCtx := context.WithTimeout(a.mainCtx, dbTimeout)
	defer cancelPostgresCtx()

	a.postgres, err = db.NewPgx(postgresCtx, a.settings.Databases.Postgres)
	if err != nil {
		return fmt.Errorf("init postgres: %w", err)
	}

	err = db.Migrate(fs, a.log, a.postgres, path, "postgres")
	if err != nil {
		return fmt.Errorf("migrate postgres: %w", err)
	}

	a.kafka = gokafka.NewKafka(a.settings.Databases.Kafka.Brokers)

	a.inspectionConsumer, err = a.kafka.Consumer(a.log.WithTags("inspectionConsumer"), func() (context.Context, context.CancelFunc) {
		return context.WithCancel(a.mainCtx)
	}, gokafka.WithTopic(a.settings.Databases.Kafka.Topics.Inspections), gokafka.WithConsumerGroup(serviceName))
	if err != nil {
		return fmt.Errorf("init inspection consumer: %w", err)
	}

	return nil
}

func (a *App) InitServices() error {
	subscriberRepository := dbsubscriber.NewRepository(a.postgres)
	objectRepository := dbobject.NewRepository(a.postgres)
	contractRepository := dbcontract.NewRepository(a.postgres, subscriberRepository, objectRepository)

	httpClient := gohttp.NewClient(gohttp.WithTimeout(1 * time.Minute))

	taskClient := task.NewClient(httpClient, a.settings.Cluster.TaskService)

	a.subscriberService = subscriber.NewService(subscriberRepository)
	a.objectService = object.NewService(objectRepository)
	a.contractService = contract.NewService(contractRepository, subscriberRepository, taskClient)
	a.registryService = registry.NewService(subscriberRepository, objectRepository, contractRepository)

	return nil
}

func (a *App) InitServer() {
	sb := api.NewServerBuilder(a.mainCtx, a.log, a.settings)
	sb.AddDebug()
	sb.AddSubscribers(a.subscriberService)
	sb.AddObjects(a.objectService)
	sb.AddContracts(a.contractService)
	sb.AddRegistry(a.registryService)

	a.server = sb.Build()
}

func (a *App) Start() {
	a.server.Start()
	a.inspectionConsumer.Subscribe(a.contractService.SubscriberOnInspectionEvent(a.mainCtx, a.log.WithTags("inspectionSubscriber")))
}

func (a *App) Stop(ctx context.Context) {
	consumerCtx, cancelConsumerCtx := context.WithTimeout(ctx, dbTimeout)
	defer cancelConsumerCtx()

	err := a.inspectionConsumer.Close(consumerCtx)
	if err != nil {
		a.log.Errorf("failed to close inspection consumer: %v", err)
	}

	a.server.Stop()

	err = a.postgres.Close()
	if err != nil {
		a.log.Errorf("failed to close postgres connection: %v", err)
	}
}
