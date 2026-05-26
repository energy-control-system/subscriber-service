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
	extended      ExtendedSubscriber
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

	filter := ListFilter{Search: "Иван"}
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
