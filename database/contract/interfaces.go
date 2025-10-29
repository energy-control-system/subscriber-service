package contract

import (
	"context"
	"subscriber-service/service/object"
	"subscriber-service/service/subscriber"
)

type SubscriberRepository interface {
	GetSubscriberByID(ctx context.Context, id int) (subscriber.Subscriber, error)
	GetAllSubscribers(ctx context.Context) ([]subscriber.Subscriber, error)
}

type ObjectRepository interface {
	GetObjectByID(ctx context.Context, id int) (object.Object, error)
	GetAllObjects(ctx context.Context) ([]object.Object, error)
}
