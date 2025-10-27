package subscriber

import "time"

type Subscriber struct {
	ID            int       `db:"id"`
	AccountNumber string    `db:"account_number"`
	Surname       string    `db:"surname"`
	Name          string    `db:"name"`
	Patronymic    string    `db:"patronymic"`
	PhoneNumber   string    `db:"phone_number"`
	Email         string    `db:"email"`
	INN           string    `db:"inn"`
	BirthDate     time.Time `db:"birth_date"`
	Status        int       `db:"status"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type Passport struct {
	ID           int    `db:"id"`
	SubscriberID int    `db:"subscriber_id"`
	Series       string `db:"series"`
	Number       string `db:"number"`
	IssuedBy     string `db:"issued_by"`
	IssueDate    string `db:"issue_date"`
}
