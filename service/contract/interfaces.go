package contract

import (
	"context"
	"subscriber-service/cluster/task"
	"subscriber-service/service/subscriber"

	"github.com/sunshineOfficial/golib/goctx"
)

type Repository interface {
	AddContract(ctx context.Context, request AddContractRequest) (Contract, error)
	GetLastContractByObjectID(ctx context.Context, objectID int) (Contract, error)
}

type SubscriberRepository interface {
	UpdateSubscriberStatus(ctx context.Context, subscriberID int, newStatus subscriber.Status) error
}

type TaskService interface {
	GetTaskByID(ctx goctx.Context, id int) (task.Task, error)
}
