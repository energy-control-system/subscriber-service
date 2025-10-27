package contract

import (
	"context"
	_ "embed"
	"fmt"
	"subscriber-service/service/contract"

	"github.com/jmoiron/sqlx"
	"github.com/sunshineOfficial/golib/db"
)

var (
	//go:embed sql/add_contract.sql
	addContractSQL string

	//go:embed sql/get_last_contract_by_object_id.sql
	getLastContractByObjectIDSQL string
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

	dbContract := MapAddContractRequestToDB(request)

	err = db.NamedGetWithDB(ctx, r.db, &dbContract, addContractSQL, dbContract)
	if err != nil {
		return contract.Contract{}, fmt.Errorf("add contract: %w", err)
	}

	newContract := MapContractFromDB(dbContract)
	newContract.Subscriber = sub
	newContract.Object = obj

	return newContract, nil
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
