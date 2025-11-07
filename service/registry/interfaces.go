package registry

import (
	"context"
	"subscriber-service/service/contract"
	"subscriber-service/service/object"
	"subscriber-service/service/subscriber"
)

type SubscriberRepository interface {
	UpsertSubscribers(ctx context.Context, subs []subscriber.UpsertSubscriberRequest) error
}

type ObjectRepository interface {
	UpsertObjects(ctx context.Context, objects []object.UpsertObjectRequest) error
	UpsertDevices(ctx context.Context, devices []object.UpsertDeviceRequest) error
	UpsertSeals(ctx context.Context, seals []object.UpsertSealRequest) error
}

type ContractRepository interface {
	UpsertContracts(ctx context.Context, contracts []contract.UpsertContractRequest) error
}
