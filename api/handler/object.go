package handler

import (
	"fmt"
	"net/http"
	"subscriber-service/service/object"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

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
