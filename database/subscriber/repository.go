package subscriber

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	dbobject "subscriber-service/database/object"
	"subscriber-service/service/object"
	"subscriber-service/service/subscriber"

	"github.com/jmoiron/sqlx"
	"github.com/sunshineOfficial/golib/db"
	"github.com/sunshineOfficial/golib/pagination"
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

	//go:embed sql/get_contracts_by_subscriber_id.sql
	getContractsBySubscriberIDSQL string

	//go:embed sql/get_devices_by_object_ids.sql
	getDevicesByObjectIDsSQL string

	//go:embed sql/get_objects_by_ids.sql
	getObjectsByIDsSQL string

	//go:embed sql/get_passport_by_subscriber_id.sql
	getPassportBySubscriberIDSQL string

	//go:embed sql/get_seals_by_object_ids.sql
	getSealsByObjectIDsSQL string

	//go:embed sql/get_subscriber_by_id.sql
	getSubscriberByIDSQL string

	//go:embed sql/update_subscriber_status.sql
	updateSubscriberStatusSQL string

	//go:embed sql/upsert_subscriber.sql
	upsertSubscriberSQL string
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
	dbSubscriber, dbPassport, err := MapAddSubscriberRequestToDB(request)
	if err != nil {
		return subscriber.Subscriber{}, fmt.Errorf("map add subscriber request: %w", err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return subscriber.Subscriber{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, fmt.Errorf("transaction rollback: %w", tx.Rollback()))
		}
	}()

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

func (r *Repository) GetSubscriberExtendedByID(ctx context.Context, id int) (subscriber.ExtendedSubscriber, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return subscriber.ExtendedSubscriber{}, fmt.Errorf("begin transaction: %w", err)
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
		return subscriber.ExtendedSubscriber{}, err
	}

	var dbPassport Passport
	err = tx.GetContext(ctx, &dbPassport, getPassportBySubscriberIDSQL, id)
	if err != nil {
		err = fmt.Errorf("get passport: %w", err)
		return subscriber.ExtendedSubscriber{}, err
	}

	var dbContracts []Contract
	err = tx.SelectContext(ctx, &dbContracts, getContractsBySubscriberIDSQL, id)
	if err != nil {
		err = fmt.Errorf("get contracts: %w", err)
		return subscriber.ExtendedSubscriber{}, err
	}

	extendedSubscriber := subscriber.ExtendedSubscriber{
		Subscriber: MapSubscriberFromDB(dbSubscriber, dbPassport),
		Contracts:  MapContractsFromDB(dbContracts),
		Objects:    []object.Object{},
	}

	objectIDs := uniqueContractObjectIDs(dbContracts)
	if len(objectIDs) > 0 {
		var dbObjects []dbobject.Object
		err = tx.SelectContext(ctx, &dbObjects, getObjectsByIDsSQL, objectIDs)
		if err != nil {
			err = fmt.Errorf("get objects: %w", err)
			return subscriber.ExtendedSubscriber{}, err
		}

		var dbDevices []dbobject.Device
		err = tx.SelectContext(ctx, &dbDevices, getDevicesByObjectIDsSQL, objectIDs)
		if err != nil {
			err = fmt.Errorf("get object devices: %w", err)
			return subscriber.ExtendedSubscriber{}, err
		}

		var dbSeals []dbobject.Seal
		err = tx.SelectContext(ctx, &dbSeals, getSealsByObjectIDsSQL, objectIDs)
		if err != nil {
			err = fmt.Errorf("get object seals: %w", err)
			return subscriber.ExtendedSubscriber{}, err
		}

		deviceMap := make(map[int][]dbobject.Device, len(dbObjects))
		for _, dbDevice := range dbDevices {
			deviceMap[dbDevice.ObjectID] = append(deviceMap[dbDevice.ObjectID], dbDevice)
		}

		extendedSubscriber.Objects = make([]object.Object, 0, len(dbObjects))
		for _, dbObject := range dbObjects {
			obj, mapErr := dbobject.MapObjectFullFromDB(dbObject, deviceMap[dbObject.ID], dbSeals)
			if mapErr != nil {
				err = fmt.Errorf("map object %d: %w", dbObject.ID, mapErr)
				return subscriber.ExtendedSubscriber{}, err
			}

			extendedSubscriber.Objects = append(extendedSubscriber.Objects, obj)
		}
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("commit transaction: %w", err)
		return subscriber.ExtendedSubscriber{}, err
	}

	return extendedSubscriber, err
}

func (r *Repository) GetAllSubscribers(ctx context.Context, page pagination.Pagination, filter subscriber.ListFilter) ([]subscriber.Subscriber, error) {
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
	err = tx.SelectContext(
		ctx,
		&dbSubscribers,
		getAllSubscribersSQL,
		filter.Surname,
		filter.Name,
		filter.Patronymic,
		filter.AccountNumber,
		filter.PhoneNumber,
		filter.Address,
		page.LimitArg(),
		page.Offset,
	)
	if err != nil {
		err = fmt.Errorf("get all subscribers: %w", err)
		return nil, err
	}
	if len(dbSubscribers) == 0 {
		if err = tx.Commit(); err != nil {
			err = fmt.Errorf("commit transaction: %w", err)
			return nil, err
		}
		return []subscriber.Subscriber{}, nil
	}

	subscriberIDs := make([]int, 0, len(dbSubscribers))
	for _, dbSubscriber := range dbSubscribers {
		subscriberIDs = append(subscriberIDs, dbSubscriber.ID)
	}

	var dbPassports []Passport
	err = tx.SelectContext(ctx, &dbPassports, getAllPassportsSQL, subscriberIDs)
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

func uniqueContractObjectIDs(contracts []Contract) []int {
	objectIDs := make([]int, 0, len(contracts))
	seen := make(map[int]struct{}, len(contracts))
	for _, c := range contracts {
		if _, ok := seen[c.ObjectID]; ok {
			continue
		}

		seen[c.ObjectID] = struct{}{}
		objectIDs = append(objectIDs, c.ObjectID)
	}

	return objectIDs
}

func (r *Repository) UpdateSubscriberStatus(ctx context.Context, subscriberID int, newStatus subscriber.Status) error {
	_, err := r.db.ExecContext(ctx, updateSubscriberStatusSQL, subscriberID, newStatus)
	if err != nil {
		return fmt.Errorf("update subscriber: %w", err)
	}

	return nil
}

func (r *Repository) UpsertSubscribers(ctx context.Context, subscribers []subscriber.UpsertSubscriberRequest) error {
	dbSubscribers, err := MapUpsertSubscriberRequestsToDB(subscribers)
	if err != nil {
		return fmt.Errorf("map requests: %w", err)
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, fmt.Errorf("transaction rollback: %w", tx.Rollback()))
		}
	}()

	for _, dbSubscriber := range dbSubscribers {
		result, sqlErr := tx.NamedExecContext(ctx, upsertSubscriberSQL, dbSubscriber)
		if sqlErr != nil {
			err = fmt.Errorf("upsert subscriber: %w", sqlErr)
			return err
		}

		rows, sqlErr := result.RowsAffected()
		if sqlErr != nil {
			err = fmt.Errorf("get rows affected: %w", sqlErr)
			return err
		}

		if rows != 1 {
			err = fmt.Errorf("rows affected = %d, expected 1", rows)
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("commit transaction: %w", err)
		return err
	}

	return err
}
