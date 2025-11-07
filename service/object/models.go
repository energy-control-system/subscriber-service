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

func ParseDevicePlaceType(s string) DevicePlaceType {
	switch s {
	case "1":
		return DevicePlaceOther
	case "2":
		return DevicePlaceFlat
	case "3":
		return DevicePlaceStairLanding
	default:
		return DevicePlaceUnknown
	}
}

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

type UpsertObjectRequest struct {
	Address       string `json:"Address"`
	HaveAutomaton bool   `json:"HaveAutomaton"`
}

type UpsertDeviceRequest struct {
	ObjectAddress    string          `json:"ObjectAddress"`
	Type             string          `json:"Type"`
	Number           string          `json:"Number"`
	PlaceType        DevicePlaceType `json:"PlaceType"`
	PlaceDescription string          `json:"PlaceDescription"`
}

type UpsertSealRequest struct {
	DeviceNumber string `json:"DeviceNumber"`
	Number       string `json:"Number"`
	Place        string `json:"Place"`
}
