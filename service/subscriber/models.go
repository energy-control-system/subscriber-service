package subscriber

import "time"

type Status int

const (
	StatusUnknown Status = iota
	StatusActive
	StatusViolator
	StatusArchived
)

type Subscriber struct {
	ID            int       `json:"ID"`
	AccountNumber string    `json:"AccountNumber"`
	Surname       string    `json:"Surname"`
	Name          string    `json:"Name"`
	Patronymic    string    `json:"Patronymic"`
	PhoneNumber   string    `json:"PhoneNumber"`
	Email         string    `json:"Email"`
	INN           string    `json:"INN"`
	BirthDate     time.Time `json:"BirthDate"`
	Status        Status    `json:"Status"`
	Passport      Passport  `json:"Passport"`
	CreatedAt     time.Time `json:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt"`
}

type Passport struct {
	ID        int    `json:"ID"`
	Series    string `json:"Series"`
	Number    string `json:"Number"`
	IssuedBy  string `json:"IssuedBy"`
	IssueDate string `json:"IssueDate"`
}

type AddSubscriberRequest struct {
	AccountNumber string                       `json:"AccountNumber"`
	Surname       string                       `json:"Surname"`
	Name          string                       `json:"Name"`
	Patronymic    string                       `json:"Patronymic"`
	PhoneNumber   string                       `json:"PhoneNumber"`
	Email         string                       `json:"Email"`
	INN           string                       `json:"INN"`
	BirthDate     string                       `json:"BirthDate"`
	Passport      AddSubscriberRequestPassport `json:"Passport"`
}

type AddSubscriberRequestPassport struct {
	Series    string `json:"Series"`
	Number    string `json:"Number"`
	IssuedBy  string `json:"IssuedBy"`
	IssueDate string `json:"IssueDate"`
}

type UpsertSubscriberRequest struct {
	AccountNumber string                       `json:"AccountNumber"`
	Surname       string                       `json:"Surname"`
	Name          string                       `json:"Name"`
	Patronymic    string                       `json:"Patronymic"`
	PhoneNumber   string                       `json:"PhoneNumber"`
	Email         string                       `json:"Email"`
	INN           string                       `json:"INN"`
	BirthDate     string                       `json:"BirthDate"`
	Passport      AddSubscriberRequestPassport `json:"Passport"`
}

type UpsertSubscriberRequestPassport struct {
	Series    string `json:"Series"`
	Number    string `json:"Number"`
	IssuedBy  string `json:"IssuedBy"`
	IssueDate string `json:"IssueDate"`
}
