package contract

import (
	"fmt"
	"subscriber-service/service/contract"
	"time"
)

func MapAddContractRequestToDB(request contract.AddContractRequest) (Contract, error) {
	signDate, err := time.ParseInLocation(time.DateOnly, request.SignDate, time.UTC)
	if err != nil {
		return Contract{}, fmt.Errorf("parse sign date: %w", err)
	}

	return Contract{
		Number:       request.Number,
		SubscriberID: request.SubscriberID,
		ObjectID:     request.ObjectID,
		SignDate:     signDate,
	}, nil
}

func MapContractFromDB(c Contract) contract.Contract {
	return contract.Contract{
		ID:        c.ID,
		Number:    c.Number,
		SignDate:  c.SignDate,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func MapUpsertContractRequestToDB(request contract.UpsertContractRequest) UpsertContractRequest {
	return UpsertContractRequest{
		Number:                  request.Number,
		SubscriberAccountNumber: request.SubscriberAccountNumber,
		ObjectAddress:           request.ObjectAddress,
		SignDate:                request.SignDate,
	}
}

func MapUpsertContractRequestsToDB(requests []contract.UpsertContractRequest) []UpsertContractRequest {
	result := make([]UpsertContractRequest, 0, len(requests))
	for _, request := range requests {
		result = append(result, MapUpsertContractRequestToDB(request))
	}

	return result
}
