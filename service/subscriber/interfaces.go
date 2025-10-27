package subscriber

import "context"

type Repository interface {
	AddSubscriber(ctx context.Context, request AddSubscriberRequest) (Subscriber, error)
	GetSubscriberByID(ctx context.Context, id int) (Subscriber, error)
}
