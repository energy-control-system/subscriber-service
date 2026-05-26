package subscriber

import (
	"context"
	"reflect"
	"testing"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/pagination"
)

type extendedMockRepository struct {
	gotExtendedID int
	gotPage       pagination.Pagination
	gotFilter     ListFilter
	gotUpdateID   int
	gotUpdate     UpdateSubscriberRequest
	gotDeleteID   int
	extended      ExtendedSubscriber
	updated       Subscriber
	deleted       Subscriber
}

func (m *extendedMockRepository) AddSubscriber(context.Context, AddSubscriberRequest) (Subscriber, error) {
	return Subscriber{}, nil
}

func (m *extendedMockRepository) GetSubscriberByID(context.Context, int) (Subscriber, error) {
	return Subscriber{}, nil
}

func (m *extendedMockRepository) GetAllSubscribers(_ context.Context, page pagination.Pagination, filter ListFilter) ([]Subscriber, error) {
	m.gotPage = page
	m.gotFilter = filter
	return []Subscriber{{ID: 7, AccountNumber: "A-100"}}, nil
}

func (m *extendedMockRepository) UpdateSubscriber(_ context.Context, id int, request UpdateSubscriberRequest) (Subscriber, error) {
	m.gotUpdateID = id
	m.gotUpdate = request
	return m.updated, nil
}

func (m *extendedMockRepository) DeleteSubscriber(_ context.Context, id int) (Subscriber, error) {
	m.gotDeleteID = id
	return m.deleted, nil
}

func (m *extendedMockRepository) GetSubscriberExtendedByID(_ context.Context, id int) (ExtendedSubscriber, error) {
	m.gotExtendedID = id
	return m.extended, nil
}

func TestGetSubscriberExtendedByIDDelegatesToRepository(t *testing.T) {
	repository := &extendedMockRepository{
		extended: ExtendedSubscriber{
			Subscriber: Subscriber{ID: 5, AccountNumber: "1234567890"},
		},
	}
	service := NewService(repository)

	got, err := service.GetSubscriberExtendedByID(goctx.Wrap(context.Background()), 5)
	if err != nil {
		t.Fatalf("GetSubscriberExtendedByID returned error: %v", err)
	}

	if repository.gotExtendedID != 5 {
		t.Fatalf("repository id = %d, want 5", repository.gotExtendedID)
	}
	if !reflect.DeepEqual(got, repository.extended) {
		t.Fatalf("extended subscriber = %+v, want %+v", got, repository.extended)
	}
}

func TestGetAllSubscribersDelegatesPaginationAndSearchFilter(t *testing.T) {
	repository := &extendedMockRepository{}
	service := NewService(repository)

	filter := ListFilter{Search: "Иван", Status: StatusActive}
	page := pagination.Pagination{Limit: 20, Offset: 40}

	got, err := service.GetAllSubscribers(goctx.Wrap(context.Background()), page, filter)
	if err != nil {
		t.Fatalf("GetAllSubscribers returned error: %v", err)
	}

	if !reflect.DeepEqual(repository.gotPage, page) {
		t.Fatalf("repository page = %+v, want %+v", repository.gotPage, page)
	}
	if !reflect.DeepEqual(repository.gotFilter, filter) {
		t.Fatalf("repository filter = %+v, want %+v", repository.gotFilter, filter)
	}
	if len(got) != 1 || got[0].ID != 7 {
		t.Fatalf("subscribers = %+v, want subscriber ID 7", got)
	}
}

func TestUpdateSubscriberDelegatesToRepository(t *testing.T) {
	request := UpdateSubscriberRequest{
		AccountNumber: "АБВ-ГДЕ-1",
		Surname:       "Ivanov",
		Name:          "Ivan",
		Patronymic:    "Ivanovich",
		PhoneNumber:   "+79991234567",
		Email:         "ivan@example.com",
		INN:           "1234567890",
		BirthDate:     "1990-01-02",
		Status:        StatusActive,
		Passport: UpdateSubscriberRequestPassport{
			Series:    "1234",
			Number:    "123456",
			IssuedBy:  "UFMS",
			IssueDate: "2010-01-02",
		},
	}
	repository := &extendedMockRepository{updated: Subscriber{ID: 11, AccountNumber: request.AccountNumber}}
	service := NewService(repository)

	got, err := service.UpdateSubscriber(goctx.Wrap(context.Background()), 11, request)
	if err != nil {
		t.Fatalf("UpdateSubscriber returned error: %v", err)
	}

	if repository.gotUpdateID != 11 {
		t.Fatalf("repository id = %d, want 11", repository.gotUpdateID)
	}
	if !reflect.DeepEqual(repository.gotUpdate, request) {
		t.Fatalf("repository request = %+v, want %+v", repository.gotUpdate, request)
	}
	if !reflect.DeepEqual(got, repository.updated) {
		t.Fatalf("subscriber = %+v, want %+v", got, repository.updated)
	}
}

func TestDeleteSubscriberDelegatesToRepository(t *testing.T) {
	repository := &extendedMockRepository{deleted: Subscriber{ID: 12, AccountNumber: "АБВ-ГДЕ-1"}}
	service := NewService(repository)

	got, err := service.DeleteSubscriber(goctx.Wrap(context.Background()), 12)
	if err != nil {
		t.Fatalf("DeleteSubscriber returned error: %v", err)
	}

	if repository.gotDeleteID != 12 {
		t.Fatalf("repository id = %d, want 12", repository.gotDeleteID)
	}
	if !reflect.DeepEqual(got, repository.deleted) {
		t.Fatalf("subscriber = %+v, want %+v", got, repository.deleted)
	}
}
