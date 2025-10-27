package contract

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"subscriber-service/cluster/inspection"
	"subscriber-service/service/subscriber"
	"time"

	"github.com/sunshineOfficial/golib/goctx"
	"github.com/sunshineOfficial/golib/gokafka"
	"github.com/sunshineOfficial/golib/golog"
)

const kafkaSubscribeTimeout = 2 * time.Minute

type Service struct {
	repository           Repository
	subscriberRepository SubscriberRepository
	taskService          TaskService
}

func NewService(repository Repository, subscriberRepository SubscriberRepository, taskService TaskService) *Service {
	return &Service{
		repository:           repository,
		subscriberRepository: subscriberRepository,
		taskService:          taskService,
	}
}

func (s *Service) AddContract(ctx goctx.Context, request AddContractRequest) (Contract, error) {
	c, err := s.repository.AddContract(ctx, request)
	if err != nil {
		return Contract{}, fmt.Errorf("add contract to repository: %w", err)
	}

	return c, nil
}

func (s *Service) GetLastContractByObjectID(ctx goctx.Context, objectID int) (Contract, error) {
	c, err := s.repository.GetLastContractByObjectID(ctx, objectID)
	if err != nil {
		return Contract{}, fmt.Errorf("get last contract from repository: %w", err)
	}

	return c, nil
}

func (s *Service) SubscriberOnInspectionEvent(mainCtx context.Context, log golog.Logger) gokafka.Subscriber {
	return func(message gokafka.Message, err error) {
		ctx, cancel := context.WithTimeout(mainCtx, kafkaSubscribeTimeout)
		defer cancel()

		if err != nil {
			log.Errorf("got error on inspection event: %v", err)
			return
		}

		var event inspection.Event
		err = json.Unmarshal(message.Value, &event)
		if err != nil {
			log.Errorf("failed to unmarshal inspection event: %v", err)
			return
		}

		switch event.Type {
		case inspection.EventTypeStart:
			err = s.handleStartedInspection(ctx, event.Inspection)
		case inspection.EventTypeFinish:
			err = s.handleFinishedInspection(ctx, event.Inspection)
		default:
			err = fmt.Errorf("unknown event type: %v", event.Type)
		}

		if err != nil {
			log.Errorf("failed to handle inspection event (type = %d): %v", event.Type, err)
			return
		}
	}
}

func (s *Service) handleStartedInspection(ctx context.Context, ins inspection.Inspection) error {
	return nil
}

func (s *Service) handleFinishedInspection(ctx context.Context, ins inspection.Inspection) error {
	if ins.Status != inspection.StatusDone {
		return fmt.Errorf("invalid inspection status: %v", ins.Status)
	}

	if ins.Type == nil {
		return errors.New("inspection type is nil")
	}

	insType := *ins.Type
	if insType != inspection.TypeVerification && insType != inspection.TypeUnauthorizedConnection {
		return nil
	}

	goCtx := goctx.Wrap(ctx)

	t, err := s.taskService.GetTaskByID(goCtx, ins.TaskID)
	if err != nil {
		return fmt.Errorf("get task by id: %w", err)
	}

	c, err := s.repository.GetLastContractByObjectID(goCtx, t.ObjectID)
	if err != nil {
		return fmt.Errorf("get last contract by object id: %w", err)
	}

	newStatus := getSubscriberStatusByInspection(ins)
	if err = s.subscriberRepository.UpdateSubscriberStatus(goCtx, c.Subscriber.ID, newStatus); err != nil {
		return fmt.Errorf("update subscriber status: %w", err)
	}

	return nil
}

func getSubscriberStatusByInspection(ins inspection.Inspection) subscriber.Status {
	if ins.IsViolationDetected != nil && *ins.IsViolationDetected {
		return subscriber.StatusViolator
	}

	return subscriber.StatusActive
}
