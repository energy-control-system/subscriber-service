package object

import (
	"context"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"subscriber-service/service/object"

	"github.com/jmoiron/sqlx"
	"github.com/sunshineOfficial/golib/db"
)

var (
	//go:embed sql/add_device.sql
	addDeviceSQL string

	//go:embed sql/add_object.sql
	addObjectSQL string

	//go:embed sql/add_seal.sql
	addSealSQL string

	//go:embed sql/get_object_by_device_id.sql
	getObjectByDeviceIDSQL string

	//go:embed sql/get_object_by_id.sql
	getObjectByIDSQL string

	//go:embed sql/get_object_by_seal_id.sql
	getObjectBySealIDSQL string

	//go:embed sql/get_object_devices.sql
	getObjectDevicesSQL string

	//go:embed sql/get_object_seals.sql
	getObjectSealsSQL string
)

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) AddObject(ctx context.Context, request object.AddObjectRequest) (object.Object, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return object.Object{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, fmt.Errorf("transaction rollback: %w", tx.Rollback()))
		}
	}()

	dbObject := MapAddObjectRequestToDB(request)
	err = db.NamedGet(tx, &dbObject, addObjectSQL, dbObject)
	if err != nil {
		err = fmt.Errorf("add object: %w", err)
		return object.Object{}, err
	}

	dbDevices := MapAddObjectRequestDevicesToDB(request.Devices, dbObject.ID)
	_, err = tx.NamedExecContext(ctx, addDeviceSQL, dbDevices)
	if err != nil {
		err = fmt.Errorf("add devices: %w", err)
		return object.Object{}, err
	}

	err = tx.SelectContext(ctx, &dbDevices, getObjectDevicesSQL, dbObject.ID)
	if err != nil {
		err = fmt.Errorf("get object devices: %w", err)
		return object.Object{}, err
	}

	deviceIDMap := make(map[string]int, len(dbDevices))
	for _, dbDevice := range dbDevices {
		deviceIDMap[fmt.Sprintf("%s/%s", dbDevice.Type, dbDevice.Number)] = dbDevice.ID
	}

	dbSeals := make([]Seal, 0, len(dbDevices))
	for _, device := range request.Devices {
		deviceKey := fmt.Sprintf("%s/%s", device.Type, device.Number)

		deviceID, ok := deviceIDMap[deviceKey]
		if !ok {
			err = fmt.Errorf("device %s not found", deviceKey)
			return object.Object{}, err
		}

		dbSeals = append(dbSeals, MapAddObjectRequestSealsToDB(device.Seals, deviceID)...)
	}

	_, err = tx.NamedExecContext(ctx, addSealSQL, dbSeals)
	if err != nil {
		err = fmt.Errorf("add seals: %w", err)
		return object.Object{}, err
	}

	err = tx.SelectContext(ctx, &dbSeals, getObjectSealsSQL, dbObject.ID)
	if err != nil {
		err = fmt.Errorf("get object seals: %w", err)
		return object.Object{}, err
	}

	newObject, err := MapObjectFullFromDB(dbObject, dbDevices, dbSeals)
	if err != nil {
		err = fmt.Errorf("map object: %w", err)
		return object.Object{}, err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("commit transaction: %w", err)
		return object.Object{}, err
	}

	return newObject, err
}

func (r *Repository) GetObjectByID(ctx context.Context, id int) (object.Object, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return object.Object{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, fmt.Errorf("transaction rollback: %w", tx.Rollback()))
		}
	}()

	var dbObject Object
	err = tx.GetContext(ctx, &dbObject, getObjectByIDSQL, id)
	if err != nil {
		err = fmt.Errorf("get object by id: %w", err)
		return object.Object{}, err
	}

	var dbDevices []Device
	err = tx.SelectContext(ctx, &dbDevices, getObjectDevicesSQL, id)
	if err != nil {
		err = fmt.Errorf("get object devices: %w", err)
		return object.Object{}, err
	}

	var dbSeals []Seal
	err = tx.SelectContext(ctx, &dbSeals, getObjectSealsSQL, id)
	if err != nil {
		err = fmt.Errorf("get object seals: %w", err)
		return object.Object{}, err
	}

	obj, err := MapObjectFullFromDB(dbObject, dbDevices, dbSeals)
	if err != nil {
		err = fmt.Errorf("map object: %w", err)
		return object.Object{}, err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("commit transaction: %w", err)
		return object.Object{}, err
	}

	return obj, err
}

func (r *Repository) GetObjectByDeviceID(ctx context.Context, deviceID int) (object.Object, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return object.Object{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, fmt.Errorf("transaction rollback: %w", tx.Rollback()))
		}
	}()

	var dbObject Object
	err = tx.GetContext(ctx, &dbObject, getObjectByDeviceIDSQL, deviceID)
	if err != nil {
		err = fmt.Errorf("get object by device id: %w", err)
		return object.Object{}, err
	}

	var dbDevices []Device
	err = tx.SelectContext(ctx, &dbDevices, getObjectDevicesSQL, dbObject.ID)
	if err != nil {
		err = fmt.Errorf("get object devices: %w", err)
		return object.Object{}, err
	}

	var dbSeals []Seal
	err = tx.SelectContext(ctx, &dbSeals, getObjectSealsSQL, dbObject.ID)
	if err != nil {
		err = fmt.Errorf("get object seals: %w", err)
		return object.Object{}, err
	}

	obj, err := MapObjectFullFromDB(dbObject, dbDevices, dbSeals)
	if err != nil {
		err = fmt.Errorf("map object: %w", err)
		return object.Object{}, err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("commit transaction: %w", err)
		return object.Object{}, err
	}

	return obj, err
}

func (r *Repository) GetObjectBySealID(ctx context.Context, sealID int) (object.Object, error) {
	tx, err := r.db.BeginTxx(ctx, &sql.TxOptions{ReadOnly: true})
	if err != nil {
		return object.Object{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if err != nil {
			err = errors.Join(err, fmt.Errorf("transaction rollback: %w", tx.Rollback()))
		}
	}()

	var dbObject Object
	err = tx.GetContext(ctx, &dbObject, getObjectBySealIDSQL, sealID)
	if err != nil {
		err = fmt.Errorf("get object by seal id: %w", err)
		return object.Object{}, err
	}

	var dbDevices []Device
	err = tx.SelectContext(ctx, &dbDevices, getObjectDevicesSQL, dbObject.ID)
	if err != nil {
		err = fmt.Errorf("get object devices: %w", err)
		return object.Object{}, err
	}

	var dbSeals []Seal
	err = tx.SelectContext(ctx, &dbSeals, getObjectSealsSQL, dbObject.ID)
	if err != nil {
		err = fmt.Errorf("get object seals: %w", err)
		return object.Object{}, err
	}

	obj, err := MapObjectFullFromDB(dbObject, dbDevices, dbSeals)
	if err != nil {
		err = fmt.Errorf("map object: %w", err)
		return object.Object{}, err
	}

	err = tx.Commit()
	if err != nil {
		err = fmt.Errorf("commit transaction: %w", err)
		return object.Object{}, err
	}

	return obj, err
}
