package registry

import (
	"errors"
	"fmt"
	"subscriber-service/service/object"

	"github.com/xuri/excelize/v2"
)

const devicesSheetName = "Приборы учета"

const (
	devicesRowNumberIdx = iota
	devicesObjectAddressIdx
	devicesTypeIdx
	devicesNumberIdx
	devicesPlaceTypeIdx
	devicesPlaceDescriptionIdx

	devicesRowLength
)

func parseDevices(registry *excelize.File) ([]object.UpsertDeviceRequest, error) {
	rows, err := registry.Rows(devicesSheetName)
	if err != nil {
		return nil, fmt.Errorf("get rows: %w", err)
	}
	defer func() {
		err = errors.Join(err, fmt.Errorf("rows error: %w", rows.Error()), fmt.Errorf("close rows: %w", rows.Close()))
	}()

	if !rows.Next() {
		return nil, ErrNoRows
	}

	var devices []object.UpsertDeviceRequest
	for i := 0; rows.Next(); i++ {
		row, columnsErr := rows.Columns()
		if columnsErr != nil {
			err = fmt.Errorf("get columns: %w", columnsErr)
			return nil, err
		}

		if len(row) != devicesRowLength {
			err = fmt.Errorf("got %d cells in row %d, want %d", len(row), i, devicesRowLength)
			return nil, err
		}

		devices = append(devices, object.UpsertDeviceRequest{
			ObjectAddress:    row[devicesObjectAddressIdx],
			Type:             row[devicesTypeIdx],
			Number:           row[devicesNumberIdx],
			PlaceType:        object.ParseDevicePlaceType(row[devicesPlaceTypeIdx]),
			PlaceDescription: row[devicesPlaceDescriptionIdx],
		})
	}

	return devices, nil
}
