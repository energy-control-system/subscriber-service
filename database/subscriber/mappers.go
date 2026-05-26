package subscriber

import (
	"fmt"
	"subscriber-service/service/subscriber"
	"time"
)

func MapAddSubscriberRequestToDB(request subscriber.AddSubscriberRequest) (Subscriber, Passport, error) {
	birthDate, err := time.ParseInLocation(time.DateOnly, request.BirthDate, time.UTC)
	if err != nil {
		return Subscriber{}, Passport{}, fmt.Errorf("parse birth date: %w", err)
	}

	return Subscriber{
			AccountNumber: request.AccountNumber,
			Surname:       request.Surname,
			Name:          request.Name,
			Patronymic:    request.Patronymic,
			PhoneNumber:   request.PhoneNumber,
			Email:         request.Email,
			INN:           request.INN,
			BirthDate:     birthDate,
			Status:        int(subscriber.StatusActive),
		}, Passport{
			Series:    request.Passport.Series,
			Number:    request.Passport.Number,
			IssuedBy:  request.Passport.IssuedBy,
			IssueDate: request.Passport.IssueDate,
		}, nil
}

func MapUpsertSubscriberRequestToDB(request subscriber.UpsertSubscriberRequest) (SubscriberWithPassport, error) {
	birthDate, err := time.ParseInLocation(time.DateOnly, request.BirthDate, time.UTC)
	if err != nil {
		return SubscriberWithPassport{}, fmt.Errorf("parse birth date: %w", err)
	}

	return SubscriberWithPassport{
		AccountNumber:     request.AccountNumber,
		Surname:           request.Surname,
		Name:              request.Name,
		Patronymic:        request.Patronymic,
		PhoneNumber:       request.PhoneNumber,
		Email:             request.Email,
		INN:               request.INN,
		BirthDate:         birthDate,
		PassportSeries:    request.Passport.Series,
		PassportNumber:    request.Passport.Number,
		PassportIssuedBy:  request.Passport.IssuedBy,
		PassportIssueDate: request.Passport.IssueDate,
	}, nil
}

func MapUpsertSubscriberRequestsToDB(requests []subscriber.UpsertSubscriberRequest) ([]SubscriberWithPassport, error) {
	result := make([]SubscriberWithPassport, 0, len(requests))
	for _, request := range requests {
		s, err := MapUpsertSubscriberRequestToDB(request)
		if err != nil {
			return nil, fmt.Errorf("map request: %w", err)
		}

		result = append(result, s)
	}

	return result, nil
}

func MapUpdateSubscriberRequestToDB(id int, request subscriber.UpdateSubscriberRequest) (UpdateSubscriberRequest, error) {
	birthDate, err := time.ParseInLocation(time.DateOnly, request.BirthDate, time.UTC)
	if err != nil {
		return UpdateSubscriberRequest{}, fmt.Errorf("parse birth date: %w", err)
	}

	return UpdateSubscriberRequest{
		ID:                id,
		AccountNumber:     request.AccountNumber,
		Surname:           request.Surname,
		Name:              request.Name,
		Patronymic:        request.Patronymic,
		PhoneNumber:       request.PhoneNumber,
		Email:             request.Email,
		INN:               request.INN,
		BirthDate:         birthDate,
		Status:            int(request.Status),
		PassportSeries:    request.Passport.Series,
		PassportNumber:    request.Passport.Number,
		PassportIssuedBy:  request.Passport.IssuedBy,
		PassportIssueDate: request.Passport.IssueDate,
	}, nil
}

func MapSubscriberFromDB(s Subscriber, p Passport) subscriber.Subscriber {
	return subscriber.Subscriber{
		ID:            s.ID,
		AccountNumber: s.AccountNumber,
		Surname:       s.Surname,
		Name:          s.Name,
		Patronymic:    s.Patronymic,
		PhoneNumber:   s.PhoneNumber,
		Email:         s.Email,
		INN:           s.INN,
		BirthDate:     s.BirthDate,
		Status:        subscriber.Status(s.Status),
		Passport:      MapPassportFromDB(p),
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

func MapPassportFromDB(p Passport) subscriber.Passport {
	return subscriber.Passport{
		ID:        p.ID,
		Series:    p.Series,
		Number:    p.Number,
		IssuedBy:  p.IssuedBy,
		IssueDate: p.IssueDate,
	}
}

func MapContractFromDB(c Contract) subscriber.Contract {
	return subscriber.Contract{
		ID:        c.ID,
		Number:    c.Number,
		ObjectID:  c.ObjectID,
		SignDate:  c.SignDate,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func MapContractsFromDB(contracts []Contract) []subscriber.Contract {
	result := make([]subscriber.Contract, 0, len(contracts))
	for _, c := range contracts {
		result = append(result, MapContractFromDB(c))
	}

	return result
}
