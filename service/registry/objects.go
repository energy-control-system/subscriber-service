package registry

import (
	"errors"
	"fmt"
	"subscriber-service/service/object"

	"github.com/xuri/excelize/v2"
)

const objectsSheetName = "Объекты"

const (
	objectsRowNumberIdx = iota
	objectsAddressIdx
	objectsHaveAutomatonIdx

	objectsRowLength
)

func parseObjects(registry *excelize.File) ([]object.UpsertObjectRequest, error) {
	rows, err := registry.Rows(objectsSheetName)
	if err != nil {
		return nil, fmt.Errorf("get rows: %w", err)
	}
	defer func() {
		err = errors.Join(err, fmt.Errorf("rows error: %w", rows.Error()), fmt.Errorf("close rows: %w", rows.Close()))
	}()

	if !rows.Next() {
		return nil, ErrNoRows
	}

	var objects []object.UpsertObjectRequest
	for i := 0; rows.Next(); i++ {
		row, columnsErr := rows.Columns()
		if columnsErr != nil {
			err = fmt.Errorf("get columns: %w", columnsErr)
			return nil, err
		}

		if len(row) != objectsRowLength {
			err = fmt.Errorf("got %d cells in row %d, want %d", len(row), i, objectsRowLength)
			return nil, err
		}

		haveAutomaton := false
		if row[objectsHaveAutomatonIdx] == "Да" {
			haveAutomaton = true
		}

		objects = append(objects, object.UpsertObjectRequest{
			Address:       row[objectsAddressIdx],
			HaveAutomaton: haveAutomaton,
		})
	}

	return objects, nil
}
