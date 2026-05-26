package contract

import (
	"context"
	"reflect"
	"testing"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/pagination"
)

type mockRepository struct {
	gotObjectIDs []int
	gotUpdateID  int
	gotUpdate    UpdateContractRequest
	gotDeleteID  int
	contracts    []Contract
	updated      Contract
	deleted      Contract
}

func (m *mockRepository) AddContract(context.Context, AddContractRequest) (Contract, error) {
	return Contract{}, nil
}

func (m *mockRepository) GetAllContracts(context.Context, pagination.Pagination) ([]Contract, error) {
	return nil, nil
}

func (m *mockRepository) GetLastContractByObjectID(context.Context, int) (Contract, error) {
	return Contract{}, nil
}

func (m *mockRepository) GetLastContractsByObjectIDs(_ context.Context, objectIDs []int) ([]Contract, error) {
	m.gotObjectIDs = append([]int(nil), objectIDs...)
	return m.contracts, nil
}

func (m *mockRepository) UpdateContract(_ context.Context, id int, request UpdateContractRequest) (Contract, error) {
	m.gotUpdateID = id
	m.gotUpdate = request
	return m.updated, nil
}

func (m *mockRepository) DeleteContract(_ context.Context, id int) (Contract, error) {
	m.gotDeleteID = id
	return m.deleted, nil
}

func TestGetLastContractsByObjectIDsDelegatesToRepository(t *testing.T) {
	repository := &mockRepository{contracts: []Contract{{ID: 1}, {ID: 2}}}
	service := NewService(repository, nil, nil)

	got, err := service.GetLastContractsByObjectIDs(goctx.Wrap(context.Background()), []int{10, 20})
	if err != nil {
		t.Fatalf("GetLastContractsByObjectIDs returned error: %v", err)
	}

	if !reflect.DeepEqual(repository.gotObjectIDs, []int{10, 20}) {
		t.Fatalf("repository object ids = %v, want [10 20]", repository.gotObjectIDs)
	}
	if !reflect.DeepEqual(got, repository.contracts) {
		t.Fatalf("contracts = %+v, want %+v", got, repository.contracts)
	}
}

func TestUpdateContractDelegatesToRepository(t *testing.T) {
	request := UpdateContractRequest{
		Number:       "C-100",
		SubscriberID: 1,
		ObjectID:     2,
		SignDate:     "2024-01-02",
	}
	repository := &mockRepository{updated: Contract{ID: 31, Number: request.Number}}
	service := NewService(repository, nil, nil)

	got, err := service.UpdateContract(goctx.Wrap(context.Background()), 31, request)
	if err != nil {
		t.Fatalf("UpdateContract returned error: %v", err)
	}

	if repository.gotUpdateID != 31 {
		t.Fatalf("repository id = %d, want 31", repository.gotUpdateID)
	}
	if !reflect.DeepEqual(repository.gotUpdate, request) {
		t.Fatalf("repository request = %+v, want %+v", repository.gotUpdate, request)
	}
	if !reflect.DeepEqual(got, repository.updated) {
		t.Fatalf("contract = %+v, want %+v", got, repository.updated)
	}
}

func TestDeleteContractDelegatesToRepository(t *testing.T) {
	repository := &mockRepository{deleted: Contract{ID: 32, Number: "C-101"}}
	service := NewService(repository, nil, nil)

	got, err := service.DeleteContract(goctx.Wrap(context.Background()), 32)
	if err != nil {
		t.Fatalf("DeleteContract returned error: %v", err)
	}

	if repository.gotDeleteID != 32 {
		t.Fatalf("repository id = %d, want 32", repository.gotDeleteID)
	}
	if !reflect.DeepEqual(got, repository.deleted) {
		t.Fatalf("contract = %+v, want %+v", got, repository.deleted)
	}
}
