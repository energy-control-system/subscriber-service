package api

import (
	"context"
	"fmt"
	"subscriber-service/api/handler"
	"subscriber-service/config"
	"subscriber-service/service/contract"
	"subscriber-service/service/object"
	"subscriber-service/service/registry"
	"subscriber-service/service/subscriber"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
	"github.com/sunshineOfficial/golib/gohttp/gorouter/middleware"
	"github.com/sunshineOfficial/golib/gohttp/gorouter/plugin"
	"github.com/sunshineOfficial/golib/gohttp/goserver"
	"github.com/sunshineOfficial/golib/golog"
)

type ServerBuilder struct {
	server goserver.Server
	router *gorouter.Router
}

func NewServerBuilder(ctx context.Context, log golog.Logger, settings config.Settings) *ServerBuilder {
	return &ServerBuilder{
		server: goserver.NewHTTPServer(ctx, log, fmt.Sprintf(":%d", settings.Port)),
		router: gorouter.NewRouter(log).Use(
			middleware.Metrics(),
			middleware.Recover,
			middleware.LogError,
		),
	}
}

func (s *ServerBuilder) AddDebug() {
	s.router.Install(plugin.NewPProf(), plugin.NewMetrics())
}

func (s *ServerBuilder) AddSubscribers(service *subscriber.Service) {
	r := s.router.SubRouter("/subscribers")
	r.HandlePost("", handler.AddSubscriber(service))
	r.HandleGet("/{id}", handler.GetSubscriberByID(service))
	r.HandleGet("", handler.GetAllSubscribers(service))
}

func (s *ServerBuilder) AddObjects(service *object.Service) {
	r := s.router.SubRouter("/objects")
	r.HandlePost("", handler.AddObject(service))
	r.HandleGet("/{id}", handler.GetObjectByID(service))
	r.HandleGet("/devices/{deviceID}", handler.GetObjectByDeviceID(service))
	r.HandleGet("/seals/{sealID}", handler.GetObjectBySealID(service))
	r.HandleGet("", handler.GetAllObjects(service))
}

func (s *ServerBuilder) AddContracts(service *contract.Service) {
	r := s.router.SubRouter("/contracts")
	r.HandlePost("", handler.AddContract(service))
	r.HandleGet("", handler.GetAllContracts(service))
	r.HandleGet("/objects/{objectID}/last", handler.GetLastContractByObjectID(service))
}

func (s *ServerBuilder) AddRegistry(service *registry.Service) {
	r := s.router.SubRouter("/registry")
	r.HandlePost("/parse", handler.ParseRegistry(service))
}

func (s *ServerBuilder) Build() goserver.Server {
	s.server.UseHandler(s.router)

	return s.server
}
