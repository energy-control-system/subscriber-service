package handler

import (
	"fmt"
	"net/http"
	"subscriber-service/service/registry"

	"github.com/sunshineOfficial/golib/gohttp/gorouter"
)

// ParseRegistry godoc
// @Summary Import registry
// @Description Parses an Excel registry and imports subscribers, objects, devices, seals, and contracts.
// @Tags registry
// @Accept multipart/form-data
// @Param File formData file true "Registry Excel file"
// @Success 200
// @Failure 400 {object} gorouter.ErrorResponse
// @Failure 500 {object} gorouter.ErrorResponse
// @Router /registry/parse [post]
func ParseRegistry(s *registry.Service) gorouter.Handler {
	return func(c gorouter.Context) error {
		files, err := c.FormFiles("File")
		if err != nil {
			return fmt.Errorf("parse form files: %w", err)
		}
		if len(files) != 1 {
			return fmt.Errorf("parse form files: got %d files, expected 1", len(files))
		}

		err = s.Parse(c.Ctx(), c.Log(), files[0])
		if err != nil {
			return fmt.Errorf("failed to parse registry: %w", err)
		}

		c.Write(http.StatusOK)

		return nil
	}
}
