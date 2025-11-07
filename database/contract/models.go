package contract

import "time"

type Contract struct {
	ID           int       `db:"id"`
	Number       string    `db:"number"`
	SubscriberID int       `db:"subscriber_id"`
	ObjectID     int       `db:"object_id"`
	SignDate     time.Time `db:"sign_date"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type UpsertContractRequest struct {
	Number                  string `db:"number"`
	SubscriberAccountNumber string `db:"subscriber_account_number"`
	ObjectAddress           string `db:"object_address"`
	SignDate                string `db:"sign_date"`
}
