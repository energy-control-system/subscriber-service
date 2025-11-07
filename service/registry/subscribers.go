package registry

import (
	"errors"
	"fmt"
	"subscriber-service/service/subscriber"

	"github.com/xuri/excelize/v2"
)

const subscribersSheetName = "Абоненты"

const (
	subscribersRowNumberIdx = iota
	subscribersAccountNumberIdx
	subscribersSurnameIdx
	subscribersNameIdx
	subscribersPatronymicIdx
	subscribersPhoneNumberIdx
	subscribersEmailIdx
	subscribersINNIdx
	subscribersBirthDateIdx
	subscribersPassportSeriesIdx
	subscribersPassportNumberIdx
	subscribersPassportIssuedByIdx
	subscribersPassportIssueDateIdx

	subscribersRowLength
)

func parseSubscribers(registry *excelize.File) ([]subscriber.UpsertSubscriberRequest, error) {
	rows, err := registry.Rows(subscribersSheetName)
	if err != nil {
		return nil, fmt.Errorf("get rows: %w", err)
	}
	defer func() {
		err = errors.Join(err, fmt.Errorf("rows error: %w", rows.Error()), fmt.Errorf("close rows: %w", rows.Close()))
	}()

	if !rows.Next() {
		return nil, ErrNoRows
	}

	var subs []subscriber.UpsertSubscriberRequest
	for i := 0; rows.Next(); i++ {
		row, columnsErr := rows.Columns()
		if columnsErr != nil {
			err = fmt.Errorf("get columns: %w", columnsErr)
			return nil, err
		}

		if len(row) != subscribersRowLength {
			err = fmt.Errorf("got %d cells in row %d, want %d", len(row), i, subscribersRowLength)
			return nil, err
		}

		subs = append(subs, subscriber.UpsertSubscriberRequest{
			AccountNumber: row[subscribersAccountNumberIdx],
			Surname:       row[subscribersSurnameIdx],
			Name:          row[subscribersNameIdx],
			Patronymic:    row[subscribersPatronymicIdx],
			PhoneNumber:   row[subscribersPhoneNumberIdx],
			Email:         row[subscribersEmailIdx],
			INN:           row[subscribersINNIdx],
			BirthDate:     row[subscribersBirthDateIdx],
			Passport: subscriber.AddSubscriberRequestPassport{
				Series:    row[subscribersPassportSeriesIdx],
				Number:    row[subscribersPassportNumberIdx],
				IssuedBy:  row[subscribersPassportIssuedByIdx],
				IssueDate: row[subscribersPassportIssueDateIdx],
			},
		})
	}

	return subs, nil
}
