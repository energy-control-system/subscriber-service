package contract

import (
	"context"
	"subscriber-service/service/object"
	"subscriber-service/service/subscriber"
)

type SubscriberRepository interface {
	GetSubscriberByID(ctx context.Context, id int) (subscriber.Subscriber, error)
}

type ObjectRepository interface {
	GetObjectByID(ctx context.Context, id int) (object.Object, error)
}
