package subscriber

import (
	"fmt"

	"github.com/sunshineOfficial/golib/goctx"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) AddSubscriber(ctx goctx.Context, request AddSubscriberRequest) (Subscriber, error) {
	sub, err := s.repository.AddSubscriber(ctx, request)
	if err != nil {
		return Subscriber{}, fmt.Errorf("add subscriber to repository: %w", err)
	}

	return sub, nil
}

func (s *Service) GetSubscriberByID(ctx goctx.Context, id int) (Subscriber, error) {
	subscriber, err := s.repository.GetSubscriberByID(ctx, id)
	if err != nil {
		return Subscriber{}, fmt.Errorf("get subscriber by id from repository: %w", err)
	}

	return subscriber, nil
}
