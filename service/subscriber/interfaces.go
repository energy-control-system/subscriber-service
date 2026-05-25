package subscriber

import (
	"context"

	"github.com/sunshineOfficial/golib/pagination"
)

type Repository interface {
	AddSubscriber(ctx context.Context, request AddSubscriberRequest) (Subscriber, error)
	GetSubscriberByID(ctx context.Context, id int) (Subscriber, error)
	GetSubscriberExtendedByID(ctx context.Context, id int) (ExtendedSubscriber, error)
	GetAllSubscribers(ctx context.Context, page pagination.Pagination) ([]Subscriber, error)
}
