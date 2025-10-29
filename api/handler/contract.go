package handler

import (
	"fmt"
	"net/http"
	"subscriber-service/service/contract"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

func AddContract(s *contract.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var request contract.AddContractRequest
		if err := c.ReadJson(&request); err != nil {
			return fmt.Errorf("failed to read add contract request: %w", err)
		}

		response, err := s.AddContract(c.Ctx(), request)
		if err != nil {
			return fmt.Errorf("failed to add contract: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

func GetAllContracts(s *contract.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		response, err := s.GetAllContracts(c.Ctx())
		if err != nil {
			return fmt.Errorf("failed to get all contracts: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

type objectIDVars struct {
	ObjectID int `path:"objectID"`
}

func GetLastContractByObjectID(s *contract.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars objectIDVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read object id: %w", err)
		}

		response, err := s.GetLastContractByObjectID(c.Ctx(), vars.ObjectID)
		if err != nil {
			return fmt.Errorf("failed to get last contract: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}
