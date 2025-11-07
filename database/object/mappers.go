package object

import (
	"fmt"
	"subscriber-service/service/object"
)

func MapAddObjectRequestToDB(request object.AddObjectRequest) Object {
	return Object{
		Address:       request.Address,
		HaveAutomaton: request.HaveAutomaton,
	}
}

func MapAddObjectRequestDevicesToDB(devices []object.AddObjectRequestDevice, objectID int) []Device {
	result := make([]Device, 0, len(devices))
	for _, device := range devices {
		result = append(result, MapAddObjectRequestDeviceToDB(device, objectID))
	}

	return result
}

func MapUpsertObjectRequestToDB(request object.UpsertObjectRequest) UpsertObjectRequest {
	return UpsertObjectRequest{
		Address:       request.Address,
		HaveAutomaton: request.HaveAutomaton,
	}
}

func MapUpsertObjectRequestsToDB(requests []object.UpsertObjectRequest) []UpsertObjectRequest {
	result := make([]UpsertObjectRequest, 0, len(requests))
	for _, request := range requests {
		result = append(result, MapUpsertObjectRequestToDB(request))
	}

	return result
}

func MapAddObjectRequestDeviceToDB(d object.AddObjectRequestDevice, objectID int) Device {
	return Device{
		ObjectID:         objectID,
		Type:             d.Type,
		Number:           d.Number,
		PlaceType:        int(d.PlaceType),
		PlaceDescription: d.PlaceDescription,
	}
}

func MapAddObjectRequestSealsToDB(seals []object.AddObjectRequestSeal, deviceID int) []Seal {
	result := make([]Seal, 0, len(seals))
	for _, seal := range seals {
		result = append(result, MapAddObjectRequestSealToDB(seal, deviceID))
	}

	return result
}

func MapUpsertDeviceRequestToDB(request object.UpsertDeviceRequest) UpsertDeviceRequest {
	return UpsertDeviceRequest{
		ObjectAddress:    request.ObjectAddress,
		Type:             request.Type,
		Number:           request.Number,
		PlaceType:        int(request.PlaceType),
		PlaceDescription: request.PlaceDescription,
	}
}

func MapUpsertDeviceRequestsToDB(requests []object.UpsertDeviceRequest) []UpsertDeviceRequest {
	result := make([]UpsertDeviceRequest, 0, len(requests))
	for _, request := range requests {
		result = append(result, MapUpsertDeviceRequestToDB(request))
	}

	return result
}

func MapAddObjectRequestSealToDB(s object.AddObjectRequestSeal, deviceID int) Seal {
	return Seal{
		DeviceID: deviceID,
		Number:   s.Number,
		Place:    s.Place,
	}
}

func MapUpsertSealRequestToDB(request object.UpsertSealRequest) UpsertSealRequest {
	return UpsertSealRequest{
		DeviceNumber: request.DeviceNumber,
		Number:       request.Number,
		Place:        request.Place,
	}
}

func MapUpsertSealRequestsToDB(requests []object.UpsertSealRequest) []UpsertSealRequest {
	result := make([]UpsertSealRequest, 0, len(requests))
	for _, request := range requests {
		result = append(result, MapUpsertSealRequestToDB(request))
	}

	return result
}

func MapObjectFromDB(o Object) object.Object {
	return object.Object{
		ID:            o.ID,
		Address:       o.Address,
		HaveAutomaton: o.HaveAutomaton,
		CreatedAt:     o.CreatedAt,
		UpdatedAt:     o.UpdatedAt,
	}
}

func MapDevicesFromDB(devices []Device) []object.Device {
	result := make([]object.Device, 0, len(devices))
	for _, device := range devices {
		result = append(result, MapDeviceFromDB(device))
	}

	return result
}

func MapDeviceFromDB(d Device) object.Device {
	return object.Device{
		ID:               d.ID,
		ObjectID:         d.ObjectID,
		Type:             d.Type,
		Number:           d.Number,
		PlaceType:        object.DevicePlaceType(d.PlaceType),
		PlaceDescription: d.PlaceDescription,
		CreatedAt:        d.CreatedAt,
		UpdatedAt:        d.UpdatedAt,
	}
}

func MapSealsFromDB(seals []Seal) []object.Seal {
	result := make([]object.Seal, 0, len(seals))
	for _, seal := range seals {
		result = append(result, MapSealFromDB(seal))
	}

	return result
}

func MapSealFromDB(s Seal) object.Seal {
	return object.Seal{
		ID:        s.ID,
		DeviceID:  s.DeviceID,
		Number:    s.Number,
		Place:     s.Place,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}

func MapObjectFullFromDB(dbObject Object, dbDevices []Device, dbSeals []Seal) (object.Object, error) {
	sealMap := make(map[int][]Seal, len(dbDevices))
	for _, seal := range dbSeals {
		sealMap[seal.DeviceID] = append(sealMap[seal.DeviceID], seal)
	}

	newObject := MapObjectFromDB(dbObject)
	newObject.Devices = MapDevicesFromDB(dbDevices)
	for i, device := range newObject.Devices {
		seals, ok := sealMap[device.ID]
		if !ok {
			return object.Object{}, fmt.Errorf("seals for device %d not found", device.ID)
		}

		newObject.Devices[i].Seals = MapSealsFromDB(seals)
	}

	return newObject, nil
}
