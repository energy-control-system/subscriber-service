package object

import "time"

type Object struct {
	ID            int       `json:"ID"`
	Address       string    `json:"Address"`
	HaveAutomaton bool      `json:"HaveAutomaton"`
	CreatedAt     time.Time `json:"CreatedAt"`
	UpdatedAt     time.Time `json:"UpdatedAt"`
	Devices       []Device  `json:"Devices"`
}

type DevicePlaceType int

const (
	DevicePlaceUnknown DevicePlaceType = iota
	DevicePlaceOther
	DevicePlaceFlat
	DevicePlaceStairLanding
)

type Device struct {
	ID               int             `json:"ID"`
	ObjectID         int             `json:"ObjectID"`
	Type             string          `json:"Type"`
	Number           string          `json:"Number"`
	PlaceType        DevicePlaceType `json:"PlaceType"`
	PlaceDescription string          `json:"PlaceDescription"`
	CreatedAt        time.Time       `json:"CreatedAt"`
	UpdatedAt        time.Time       `json:"UpdatedAt"`
	Seals            []Seal          `json:"Seals"`
}

type Seal struct {
	ID        int       `json:"ID"`
	DeviceID  int       `json:"DeviceID"`
	Number    string    `json:"Number"`
	Place     string    `json:"Place"`
	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`
}

type AddObjectRequest struct {
	Address       string                   `json:"Address"`
	HaveAutomaton bool                     `json:"HaveAutomaton"`
	Devices       []AddObjectRequestDevice `json:"Devices"`
}

type AddObjectRequestDevice struct {
	Type             string                 `json:"Type"`
	Number           string                 `json:"Number"`
	PlaceType        DevicePlaceType        `json:"PlaceType"`
	PlaceDescription string                 `json:"PlaceDescription"`
	Seals            []AddObjectRequestSeal `json:"Seals"`
}

type AddObjectRequestSeal struct {
	Number string `json:"Number"`
	Place  string `json:"Place"`
}
