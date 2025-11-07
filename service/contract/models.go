package contract

import (
	"subscriber-service/service/object"
	"subscriber-service/service/subscriber"
	"time"
)

type Contract struct {
	ID         int                   `json:"ID"`
	Number     string                `json:"Number"`
	Subscriber subscriber.Subscriber `json:"Subscriber"`
	Object     object.Object         `json:"Object"`
	SignDate   time.Time             `json:"SignDate"`
	CreatedAt  time.Time             `json:"CreatedAt"`
	UpdatedAt  time.Time             `json:"UpdatedAt"`
}

type AddContractRequest struct {
	Number       string `json:"Number"`
	SubscriberID int    `json:"SubscriberID"`
	ObjectID     int    `json:"ObjectID"`
	SignDate     string `json:"SignDate"`
}

type UpsertContractRequest struct {
	Number                  string `json:"Number"`
	SubscriberAccountNumber string `json:"SubscriberAccountNumber"`
	ObjectAddress           string `json:"ObjectAddress"`
	SignDate                string `json:"SignDate"`
}
