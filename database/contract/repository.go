package contract

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"subscriber-service/service/contract"
	"subscriber-service/service/object"
	"subscriber-service/service/subscriber"

	"github.com/jmoiron/sqlx"
	"github.com/sunshineOfficial/golib/db"
)

var (
	//go:embed sql/add_contract.sql
	addContractSQL string

	//go:embed sql/get_all_contracts.sql
	getAllContractsSQL string

	//go:embed sql/get_last_contract_by_object_id.sql
	getLastContractByObjectIDSQL string

	//go:embed sql/upsert_contract.sql
	upsertContractSQL string
)

type Repository struct {
	db                   *sqlx.DB
	subscriberRepository SubscriberRepository
	objectRepository     ObjectRepository
}

func NewRepository(db *sqlx.DB, subscriberRepository SubscriberRepository, objectRepository ObjectRepository) *Repository {
	return &Repository{
		db:                   db,
		subscriberRepository: subscriberRepository,
		objectRepository:     objectRepository,
	}
}

func (r *Repository) AddContract(ctx context.Context, request contract.AddContractRequest) (contract.Contract, error) {
	sub, err := r.subscriberRepository.GetSubscriberByID(ctx, request.SubscriberID)
	if err != nil {
		return contract.Contract{}, fmt.Errorf("get subscriber: %w", err)
	}

	obj, err := r.objectRepository.GetObjectByID(ctx, request.ObjectID)
	if err != nil {
		return contract.Contract{}, fmt.Errorf("get object: %w", err)
	}

	dbContract, err := MapAddContractRequestToDB(request)
	if err != nil {
		return contract.Contract{}, fmt.Errorf("map contract: %w", err)
	}

	err = db.NamedGetWithDB(ctx, r.db, &dbContract, addContractSQL, dbContract)
	if err != nil {
		return contract.Contract{}, fmt.Errorf("add contract: %w", err)
	}

	newContract := MapContractFromDB(dbContract)
	newContract.Subscriber = sub
	newContract.Object = obj

	return newContract, nil
}

func (r *Repository) GetAllContracts(ctx context.Context) ([]contract.Contract, error) {
	var dbContracts []Contract
	err := r.db.SelectContext(ctx, &dbContracts, getAllContractsSQL)
	if err != nil {
		return nil, fmt.Errorf("get all contracts: %w", err)
	}

	subscribers, err := r.subscriberRepository.GetAllSubscribers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all subscribers: %w", err)
	}

	objects, err := r.objectRepository.GetAllObjects(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all objects: %w", err)
	}

	subscriberMap := make(map[int]subscriber.Subscriber)
	for _, sub := range subscribers {
		subscriberMap[sub.ID] = sub
	}

	objectMap := make(map[int]object.Object)
	for _, obj := range objects {
		objectMap[obj.ID] = obj
	}

	contracts := make([]contract.Contract, 0, len(dbContracts))
	for _, dbContract := range dbContracts {
		sub, ok := subscriberMap[dbContract.SubscriberID]
		if !ok {
			return nil, fmt.Errorf("subscriber %d not found for contract %d", dbContract.SubscriberID, dbContract.ID)
		}

		obj, ok := objectMap[dbContract.ObjectID]
		if !ok {
			return nil, fmt.Errorf("object %d not found for contract %d", dbContract.ObjectID, dbContract.ID)
		}

		newContract := MapContractFromDB(dbContract)
		newContract.Subscriber = sub
		newContract.Object = obj

		contracts = append(contracts, newContract)
	}

	return contracts, nil
}

func (r *Repository) GetLastContractByObjectID(ctx context.Context, objectID int) (contract.Contract, error) {
	var dbContract Contract
	err := r.db.GetContext(ctx, &dbContract, getLastContractByObjectIDSQL, objectID)
	if err != nil {
		return contract.Contract{}, fmt.Errorf("get last contract by object id: %w", err)
	}

	sub, err := r.subscriberRepository.GetSubscriberByID(ctx, dbContract.SubscriberID)
	if err != nil {
		return contract.Contract{}, fmt.Errorf("get subscriber: %w", err)
	}

	obj, err := r.objectRepository.GetObjectByID(ctx, dbContract.ObjectID)
	if err != nil {
		return contract.Contract{}, fmt.Errorf("get object: %w", err)
	}

	newContract := MapContractFromDB(dbContract)
	newContract.Subscriber = sub
	newContract.Object = obj

	return newContract, nil
}

func (r *Repository) UpsertContracts(ctx context.Context, contracts []contract.UpsertContractRequest) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, fmt.Errorf("transaction rollback: %w", tx.Rollback()))
		}
	}()

	dbContracts := MapUpsertContractRequestsToDB(contracts)
	for _, dbContract := range dbContracts {
		result, sqlErr := tx.NamedExecContext(ctx, upsertContractSQL, dbContract)
		if sqlErr != nil {
			err = fmt.Errorf("upsert contract: %w", sqlErr)
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
