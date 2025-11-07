package registry

import (
	"errors"
	"fmt"
	"subscriber-service/service/object"

	"github.com/xuri/excelize/v2"
)

const sealsSheetName = "Пломбы"

const (
	sealsRowNumberIdx = iota
	sealsDeviceNumberIdx
	sealsNumberIdx
	sealsPlaceIdx

	sealsRowLength
)

func parseSeals(registry *excelize.File) ([]object.UpsertSealRequest, error) {
	rows, err := registry.Rows(sealsSheetName)
	if err != nil {
		return nil, fmt.Errorf("get rows: %w", err)
	}
	defer func() {
		err = errors.Join(err, fmt.Errorf("rows error: %w", rows.Error()), fmt.Errorf("close rows: %w", rows.Close()))
	}()

	if !rows.Next() {
		return nil, ErrNoRows
	}

	var seals []object.UpsertSealRequest
	for i := 0; rows.Next(); i++ {
		row, columnsErr := rows.Columns()
		if columnsErr != nil {
			err = fmt.Errorf("get columns: %w", columnsErr)
			return nil, err
		}

		if len(row) != sealsRowLength {
			err = fmt.Errorf("got %d cells in row %d, want %d", len(row), i, sealsRowLength)
			return nil, err
		}

		seals = append(seals, object.UpsertSealRequest{
			DeviceNumber: row[sealsDeviceNumberIdx],
			Number:       row[sealsNumberIdx],
			Place:        row[sealsPlaceIdx],
		})
	}

	return seals, nil
}
