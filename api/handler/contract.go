package handler

import (
	"fmt"
	"net/http"
	"subscriber-service/service/contract"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
	"github.com/sunshineOfficial/golib/pagination"
)

// AddContract godoc
// @Summary Create contract
// @Description Creates a contract between subscriber and metering object.
// @Tags contracts
// @Produce json
// @Param request body contract.AddContractRequest true "Contract creation payload"
// @Success 200 {object} contract.Contract
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /contracts [post]
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

// GetAllContracts godoc
// @Summary List contracts
// @Description Returns all contracts.
// @Tags contracts
// @Produce json
// @Param limit query int false "Maximum number of items to return; 0 means no limit"
// @Param offset query int false "Number of items to skip"
// @Success 200 {array} contract.Contract
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /contracts [get]
func GetAllContracts(s *contract.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars pagination.Pagination
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read pagination: %w", err)
		}

		response, err := s.GetAllContracts(c.Ctx(), vars)
		if err != nil {
			return fmt.Errorf("failed to get all contracts: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

type objectIDVars struct {
	ObjectID int `path:"objectID"`
}

type objectIDsVars struct {
	IDs []int `query:"id"`
}

// GetLastContractsByObjectIDs godoc
// @Summary Get latest object contracts
// @Description Returns latest contracts for several metering objects.
// @Tags contracts
// @Produce json
// @Param id query []int true "Object IDs" collectionFormat(multi)
// @Success 200 {array} contract.Contract
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /contracts/objects/last [get]
func GetLastContractsByObjectIDs(s *contract.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars objectIDsVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read object ids: %w", err)
		}

		response, err := s.GetLastContractsByObjectIDs(c.Ctx(), vars.IDs)
		if err != nil {
			return fmt.Errorf("failed to get last contracts: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

// GetLastContractByObjectID godoc
// @Summary Get latest object contract
// @Description Returns the latest contract for a metering object.
// @Tags contracts
// @Produce json
// @Param objectID path int true "Object ID"
// @Success 200 {object} contract.Contract
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 404 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /contracts/objects/{objectID}/last [get]
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
