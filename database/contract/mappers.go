package contract

import (
	"subscriber-service/service/contract"
)

func MapAddContractRequestToDB(request contract.AddContractRequest) Contract {
	return Contract{
		Number:       request.Number,
		SubscriberID: request.SubscriberID,
		ObjectID:     request.ObjectID,
		SignDate:     request.SignDate,
	}
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
