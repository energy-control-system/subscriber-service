package object

import (
	"context"

	"github.com/sunshineOfficial/golib/pagination"
)

type Repository interface {
	AddObject(ctx context.Context, request AddObjectRequest) (Object, error)
	GetObjectByID(ctx context.Context, id int) (Object, error)
	GetObjectByDeviceID(ctx context.Context, deviceID int) (Object, error)
	GetObjectBySealID(ctx context.Context, sealID int) (Object, error)
	GetAllObjects(ctx context.Context, page pagination.Pagination) ([]Object, error)
	UpdateObject(ctx context.Context, id int, request UpdateObjectRequest) (Object, error)
	DeleteObject(ctx context.Context, id int) (Object, error)
}
