package handler

import (
	"fmt"
	"net/http"
	"subscriber-service/service/subscriber"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
	"github.com/sunshineOfficial/golib/pagination"
)

// AddSubscriber godoc
// @Summary Create subscriber
// @Description Creates a subscriber with passport data.
// @Tags subscribers
// @Produce json
// @Param request body subscriber.AddSubscriberRequest true "Subscriber creation payload"
// @Success 200 {object} subscriber.Subscriber
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /subscribers [post]
func AddSubscriber(s *subscriber.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var request subscriber.AddSubscriberRequest
		if err := c.ReadJson(&request); err != nil {
			return fmt.Errorf("failed to read add subscriber request: %w", err)
		}

		response, err := s.AddSubscriber(c.Ctx(), request)
		if err != nil {
			return fmt.Errorf("failed to add subscriber: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

type idVars struct {
	ID int `path:"id"`
}

// GetSubscriberByID godoc
// @Summary Get subscriber by ID
// @Description Returns subscriber data by identifier.
// @Tags subscribers
// @Produce json
// @Param id path int true "Subscriber ID"
// @Success 200 {object} subscriber.Subscriber
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 404 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /subscribers/{id} [get]
func GetSubscriberByID(s *subscriber.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars idVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read subscriber id: %w", err)
		}

		response, err := s.GetSubscriberByID(c.Ctx(), vars.ID)
		if err != nil {
			return fmt.Errorf("failed to get subscriber: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

// GetSubscriberExtendedByID godoc
// @Summary Get subscriber with contracts and objects
// @Description Returns subscriber data with all related contracts and objects.
// @Tags subscribers
// @Produce json
// @Param id path int true "Subscriber ID"
// @Success 200 {object} subscriber.ExtendedSubscriber
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 404 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /subscribers/{id}/extended [get]
func GetSubscriberExtendedByID(s *subscriber.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars idVars
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read subscriber id: %w", err)
		}

		response, err := s.GetSubscriberExtendedByID(c.Ctx(), vars.ID)
		if err != nil {
			return fmt.Errorf("failed to get extended subscriber: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}

// GetAllSubscribers godoc
// @Summary List subscribers
// @Description Returns all subscribers.
// @Tags subscribers
// @Produce json
// @Param limit query int false "Maximum number of items to return; 0 means no limit"
// @Param offset query int false "Number of items to skip"
// @Param surname query string false "Search by subscriber surname"
// @Param name query string false "Search by subscriber name"
// @Param patronymic query string false "Search by subscriber patronymic"
// @Param accountNumber query string false "Search by subscriber account number"
// @Param phoneNumber query string false "Search by subscriber phone number"
// @Param address query string false "Search by metering object address"
// @Success 200 {array} subscriber.Subscriber
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /subscribers [get]
func GetAllSubscribers(s *subscriber.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		var vars struct {
			Limit         int    `query:"limit"`
			Offset        int    `query:"offset"`
			Surname       string `query:"surname"`
			Name          string `query:"name"`
			Patronymic    string `query:"patronymic"`
			AccountNumber string `query:"accountNumber"`
			PhoneNumber   string `query:"phoneNumber"`
			Address       string `query:"address"`
		}
		if err := c.Vars(&vars); err != nil {
			return fmt.Errorf("failed to read list params: %w", err)
		}

		response, err := s.GetAllSubscribers(
			c.Ctx(),
			pagination.Pagination{Limit: vars.Limit, Offset: vars.Offset},
			subscriber.ListFilter{
				Surname:       vars.Surname,
				Name:          vars.Name,
				Patronymic:    vars.Patronymic,
				AccountNumber: vars.AccountNumber,
				PhoneNumber:   vars.PhoneNumber,
				Address:       vars.Address,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to get all subscribers: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}
