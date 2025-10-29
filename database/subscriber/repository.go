package subscriber

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"subscriber-service/service/subscriber"

	"github.com/jmoiron/sqlx"
	"github.com/sunshineOfficial/golib/db"
)

var (
	//go:embed sql/add_passport.sql
	addPassportSQL string

	//go:embed sql/add_subscriber.sql
	addSubscriberSQL string

	//go:embed sql/get_all_passports.sql
	getAllPassportsSQL string

	//go:embed sql/get_all_subscribers.sql
	getAllSubscribersSQL string

	//go:embed sql/get_passport_by_subscriber_id.sql
	getPassportBySubscriberIDSQL string

	//go:embed sql/get_subscriber_by_id.sql
	getSubscriberByIDSQL string

	//go:embed sql/update_subscriber_status.sql
	updateSubscriberStatusSQL string
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) AddSubscriber(ctx context.Context, request subscriber.AddSubscriberRequest) (subscriber.Subscriber, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return subscriber.Subscriber{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, fmt.Errorf("transaction rollback: %w", tx.Rollback()))
		}
	}()

	dbSubscriber, dbPassport, err := MapAddSubscriberRequestToDB(request)
	if err != nil {
		return subscriber.Subscriber{}, fmt.Errorf("map add subscriber request: %w", err)
	}

	err = db.NamedGet(tx, &dbSubscriber, addSubscriberSQL, dbSubscriber)
	if err != nil {
		err = fmt.Errorf("add subscriber: %w", err)
		return subscriber.Subscriber{}, err
	}

	dbPassport.SubscriberID = dbSubscriber.ID
	err = db.NamedGet(tx, &dbPassport, addPassportSQL, dbPassport)
	if err != nil {
		err = fmt.Errorf("add passport: %w", err)
		return subscriber.Subscriber{}, err
	}

	newSubscriber := MapSubscriberFromDB(dbSubscriber, dbPassport)

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("commit transaction: %w", err)
		return subscriber.Subscriber{}, err
	}

	return newSubscriber, err
}

func (r *Repository) GetSubscriberByID(ctx context.Context, id int) (subscriber.Subscriber, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return subscriber.Subscriber{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, fmt.Errorf("transaction rollback: %w", tx.Rollback()))
		}
	}()

	var dbSubscriber Subscriber
	err = tx.GetContext(ctx, &dbSubscriber, getSubscriberByIDSQL, id)
	if err != nil {
		err = fmt.Errorf("get subscriber: %w", err)
		return subscriber.Subscriber{}, err
	}

	var dbPassport Passport
	err = tx.GetContext(ctx, &dbPassport, getPassportBySubscriberIDSQL, id)
	if err != nil {
		err = fmt.Errorf("get passport: %w", err)
		return subscriber.Subscriber{}, err
	}

	newSubscriber := MapSubscriberFromDB(dbSubscriber, dbPassport)

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("commit transaction: %w", err)
		return subscriber.Subscriber{}, err
	}

	return newSubscriber, err
}

func (r *Repository) GetAllSubscribers(ctx context.Context) ([]subscriber.Subscriber, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, fmt.Errorf("transaction rollback: %w", tx.Rollback()))
		}
	}()

	var dbSubscribers []Subscriber
	err = tx.SelectContext(ctx, &dbSubscribers, getAllSubscribersSQL)
	if err != nil {
		err = fmt.Errorf("get all subscribers: %w", err)
		return nil, err
	}

	var dbPassports []Passport
	err = tx.SelectContext(ctx, &dbPassports, getAllPassportsSQL)
	if err != nil {
		err = fmt.Errorf("get all passports: %w", err)
		return nil, err
	}

	passportMap := make(map[int]Passport, len(dbPassports))
	for _, dbPassport := range dbPassports {
		passportMap[dbPassport.SubscriberID] = dbPassport
	}

	subscribers := make([]subscriber.Subscriber, 0, len(dbSubscribers))
	for _, dbSubscriber := range dbSubscribers {
		dbPassport, ok := passportMap[dbSubscriber.ID]
		if !ok {
			err = fmt.Errorf("passport not found for subscriber %d", dbSubscriber.ID)
			return nil, err
		}

		subscribers = append(subscribers, MapSubscriberFromDB(dbSubscriber, dbPassport))
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("commit transaction: %w", err)
		return nil, err
	}

	return subscribers, err
}

func (r *Repository) UpdateSubscriberStatus(ctx context.Context, subscriberID int, newStatus subscriber.Status) error {
	_, err := r.db.ExecContext(ctx, updateSubscriberStatusSQL, subscriberID, newStatus)
	if err != nil {
		return fmt.Errorf("update subscriber: %w", err)
	}

	return nil
}
