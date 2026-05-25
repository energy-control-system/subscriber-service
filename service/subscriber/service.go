package subscriber

import (
	"fmt"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/pagination"
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
	if err := ValidateAccountNumber(request.AccountNumber); err != nil {
		return Subscriber{}, fmt.Errorf("validate account number: %w", err)
	}

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

func (s *Service) GetSubscriberExtendedByID(ctx goctx.Context, id int) (ExtendedSubscriber, error) {
	subscriber, err := s.repository.GetSubscriberExtendedByID(ctx, id)
	if err != nil {
		return ExtendedSubscriber{}, fmt.Errorf("get extended subscriber by id from repository: %w", err)
	}

	return subscriber, nil
}

func (s *Service) GetAllSubscribers(ctx goctx.Context, page pagination.Pagination) ([]Subscriber, error) {
	if err := page.Validate(); err != nil {
		return nil, fmt.Errorf("validate pagination: %w", err)
	}

	subscribers, err := s.repository.GetAllSubscribers(ctx, page)
	if err != nil {
		return nil, fmt.Errorf("get all subscribers from repository: %w", err)
	}

	return subscribers, nil
}
