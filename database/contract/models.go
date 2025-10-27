package contract

import "time"

type Contract struct {
	ID           int       `db:"id"`
	Number       string    `db:"number"`
	SubscriberID int       `db:"subscriber_id"`
	ObjectID     int       `db:"object_id"`
	SignDate     string    `db:"sign_date"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
