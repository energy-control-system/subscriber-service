package handler

import (
	"fmt"
	"net/http"
	"subscriber-service/service/object"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
	"github.com/sunshineOfficial/golib/pagination"
)

// AddObject godoc
// @Summary Create metering object
// @Description Creates a metering object with devices and seals.
// @Tags objects
// @Produce json
// @Param request body object.AddObjectRequest true "Object creation payload"
// @Success 200 {object} object.Object
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Security bearer
// @Router /objects [post]
func AddObject(s *object.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var request object.AddObjectRequest
		if err := c.ReadJson(&request); err != nil {
			return fmt.Errorf("failed to read add object request: %w", err)
		}

		response, err := s.AddObject(c.Ctx(), request)
		if err != nil {
			return fmt.Errorf("failed to add object: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

// GetObjectByID godoc
// @Summary Get object by ID
// @Description Returns a metering object by identifier.
// @Tags objects
// @Produce json
// @Param id path int true "Object ID"
// @Success 200 {object} object.Object
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 404 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Security bearer
// @Router /objects/{id} [get]
func GetObjectByID(s *object.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars idVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read object id: %w", err)
		}

		response, err := s.GetObjectByID(c.Ctx(), vars.ID)
		if err != nil {
			return fmt.Errorf("failed to get object: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

type deviceIDVars struct {
	DeviceID int `path:"deviceID"`
}

// GetObjectByDeviceID godoc
// @Summary Get object by device ID
// @Description Returns a metering object containing the specified device.
// @Tags objects
// @Produce json
// @Param deviceID path int true "Device ID"
// @Success 200 {object} object.Object
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 404 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /objects/devices/{deviceID} [get]
func GetObjectByDeviceID(s *object.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars deviceIDVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read device id: %w", err)
		}

		response, err := s.GetObjectByDeviceID(c.Ctx(), vars.DeviceID)
		if err != nil {
			return fmt.Errorf("failed to get object by device id: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

type sealIDVars struct {
	SealID int `path:"sealID"`
}

// GetObjectBySealID godoc
// @Summary Get object by seal ID
// @Description Returns a metering object containing the specified seal.
// @Tags objects
// @Produce json
// @Param sealID path int true "Seal ID"
// @Success 200 {object} object.Object
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 404 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /objects/seals/{sealID} [get]
func GetObjectBySealID(s *object.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars sealIDVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read seal id: %w", err)
		}

		response, err := s.GetObjectBySealID(c.Ctx(), vars.SealID)
		if err != nil {
			return fmt.Errorf("failed to get object by seal id: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

// GetAllObjects godoc
// @Summary List objects
// @Description Returns all metering objects.
// @Tags objects
// @Produce json
// @Param limit query int false "Maximum number of items to return; 0 means no limit"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} object.Object
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Security bearer
// @Router /objects [get]
func GetAllObjects(s *object.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars pagination.Pagination
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read pagination: %w", err)
		}

		response, err := s.GetAllObjects(c.Ctx(), vars)
		if err != nil {
			return fmt.Errorf("failed to get all objects: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

// UpdateObject godoc
// @Summary Update metering object
// @Description Updates a metering object's main data.
// @Tags objects
// @Produce json
// @Param id path int true "Object ID"
// @Param request body object.UpdateObjectRequest true "Object update payload"
// @Success 200 {object} object.Object
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 404 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Security bearer
// @Router /objects/{id} [patch]
func UpdateObject(s *object.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars idVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read object id: %w", err)
		}

		var request object.UpdateObjectRequest
		if err := c.ReadJson(&request); err != nil {
			return fmt.Errorf("failed to read update object request: %w", err)
		}

		response, err := s.UpdateObject(c.Ctx(), vars.ID, request)
		if err != nil {
			return fmt.Errorf("failed to update object: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

// DeleteObject godoc
// @Summary Delete metering object
// @Description Deletes a metering object by identifier.
// @Tags objects
// @Produce json
// @Param id path int true "Object ID"
// @Success 200 {object} object.Object
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 404 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Security bearer
// @Router /objects/{id} [delete]
func DeleteObject(s *object.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars idVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read object id: %w", err)
		}

		response, err := s.DeleteObject(c.Ctx(), vars.ID)
		if err != nil {
			return fmt.Errorf("failed to delete object: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}
