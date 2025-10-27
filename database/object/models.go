package object

import "time"

type Object struct {
	ID            int       `db:"id"`
	Address       string    `db:"address"`
	HaveAutomaton bool      `db:"have_automaton"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type Device struct {
	ID               int       `db:"id"`
	ObjectID         int       `db:"object_id"`
	Type             string    `db:"type"`
	Number           string    `db:"number"`
	PlaceType        int       `db:"place_type"`
	PlaceDescription string    `db:"place_description"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type Seal struct {
	ID        int       `db:"id"`
	DeviceID  int       `db:"device_id"`
	Number    string    `db:"number"`
	Place     string    `db:"place"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
