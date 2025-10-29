package object

import "context"

type Repository interface {
	AddObject(ctx context.Context, request AddObjectRequest) (Object, error)
	GetObjectByID(ctx context.Context, id int) (Object, error)
	GetObjectByDeviceID(ctx context.Context, deviceID int) (Object, error)
	GetObjectBySealID(ctx context.Context, sealID int) (Object, error)
	GetAllObjects(ctx context.Context) ([]Object, error)
}
