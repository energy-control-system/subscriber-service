package object

import (
	"context"
	"reflect"
	"testing"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/pagination"
)

type mockRepository struct {
	gotUpdateID int
	gotUpdate   UpdateObjectRequest
	gotDeleteID int
	updated     Object
	deleted     Object
}

func (m *mockRepository) AddObject(context.Context, AddObjectRequest) (Object, error) {
	return Object{}, nil
}

func (m *mockRepository) GetObjectByID(context.Context, int) (Object, error) {
	return Object{}, nil
}

func (m *mockRepository) GetObjectByDeviceID(context.Context, int) (Object, error) {
	return Object{}, nil
}

func (m *mockRepository) GetObjectBySealID(context.Context, int) (Object, error) {
	return Object{}, nil
}

func (m *mockRepository) GetAllObjects(context.Context, pagination.Pagination) ([]Object, error) {
	return nil, nil
}

func (m *mockRepository) UpdateObject(_ context.Context, id int, request UpdateObjectRequest) (Object, error) {
	m.gotUpdateID = id
	m.gotUpdate = request
	return m.updated, nil
}

func (m *mockRepository) DeleteObject(_ context.Context, id int) (Object, error) {
	m.gotDeleteID = id
	return m.deleted, nil
}

func TestUpdateObjectDelegatesToRepository(t *testing.T) {
	request := UpdateObjectRequest{Address: "Main st. 10", HaveAutomaton: true}
	repository := &mockRepository{updated: Object{ID: 21, Address: request.Address, HaveAutomaton: true}}
	service := NewService(repository)

	got, err := service.UpdateObject(goctx.Wrap(context.Background()), 21, request)
	if err != nil {
		t.Fatalf("UpdateObject returned error: %v", err)
	}

	if repository.gotUpdateID != 21 {
		t.Fatalf("repository id = %d, want 21", repository.gotUpdateID)
	}
	if !reflect.DeepEqual(repository.gotUpdate, request) {
		t.Fatalf("repository request = %+v, want %+v", repository.gotUpdate, request)
	}
	if !reflect.DeepEqual(got, repository.updated) {
		t.Fatalf("object = %+v, want %+v", got, repository.updated)
	}
}

func TestDeleteObjectDelegatesToRepository(t *testing.T) {
	repository := &mockRepository{deleted: Object{ID: 22, Address: "Main st. 11"}}
	service := NewService(repository)

	got, err := service.DeleteObject(goctx.Wrap(context.Background()), 22)
	if err != nil {
		t.Fatalf("DeleteObject returned error: %v", err)
	}

	if repository.gotDeleteID != 22 {
		t.Fatalf("repository id = %d, want 22", repository.gotDeleteID)
	}
	if !reflect.DeepEqual(got, repository.deleted) {
		t.Fatalf("object = %+v, want %+v", got, repository.deleted)
	}
}
