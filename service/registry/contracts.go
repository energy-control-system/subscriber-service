package registry

import (
	"errors"
	"fmt"
	"subscriber-service/service/contract"

	"github.com/xuri/excelize/v2"
)

const contractsSheetName = "Договора"

const (
	contractsRowNumberIdx = iota
	contractsNumberIdx
	contractsSignDateIdx
	contractsSubscriberAccountNumberIdx
	contractsObjectAddressIdx

	contractsRowLength
)

func parseContracts(registry *excelize.File) ([]contract.UpsertContractRequest, error) {
	rows, err := registry.Rows(contractsSheetName)
	if err != nil {
		return nil, fmt.Errorf("get rows: %w", err)
	}
	defer func() {
		err = errors.Join(err, fmt.Errorf("rows error: %w", rows.Error()), fmt.Errorf("close rows: %w", rows.Close()))
	}()

	if !rows.Next() {
		return nil, ErrNoRows
	}

	var contracts []contract.UpsertContractRequest
	for i := 0; rows.Next(); i++ {
		row, columnsErr := rows.Columns()
		if columnsErr != nil {
			err = fmt.Errorf("get columns: %w", columnsErr)
			return nil, err
		}

		if len(row) != contractsRowLength {
			err = fmt.Errorf("got %d cells in row %d, want %d", len(row), i, contractsRowLength)
			return nil, err
		}

		contracts = append(contracts, contract.UpsertContractRequest{
			Number:                  row[contractsNumberIdx],
			SubscriberAccountNumber: row[contractsSubscriberAccountNumberIdx],
			ObjectAddress:           row[contractsObjectAddressIdx],
			SignDate:                row[contractsSignDateIdx],
		})
	}

	return contracts, nil
}
