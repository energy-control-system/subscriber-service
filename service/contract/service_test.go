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
	contracts    []Contract
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
