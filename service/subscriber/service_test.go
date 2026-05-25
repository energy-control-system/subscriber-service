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
	extended      ExtendedSubscriber
}

func (m *extendedMockRepository) AddSubscriber(context.Context, AddSubscriberRequest) (Subscriber, error) {
	return Subscriber{}, nil
}

func (m *extendedMockRepository) GetSubscriberByID(context.Context, int) (Subscriber, error) {
	return Subscriber{}, nil
}

func (m *extendedMockRepository) GetAllSubscribers(context.Context, pagination.Pagination) ([]Subscriber, error) {
	return nil, nil
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
