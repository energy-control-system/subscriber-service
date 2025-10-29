package handler

import (
	"fmt"
	"net/http"
	"subscriber-service/service/subscriber"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

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

func GetAllSubscribers(s *subscriber.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		response, err := s.GetAllSubscribers(c.Ctx())
		if err != nil {
			return fmt.Errorf("failed to get all subscribers: %w", err)
		}

		return c.WriteJson(http.StatusOK, response)
	}
}
