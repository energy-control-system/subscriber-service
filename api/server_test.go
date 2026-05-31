package api

import (
	"net/http"
	"net/http/httptest"
	"subscriber-service/config"
	"testing"

	"github.com/sunshineOfficial/golib/golog"
)

func TestSubscriberAuthorizationPolicy(t *testing.T) {
	builder := NewServerBuilder(t.Context(), golog.NewLogger("test"), config.Settings{
		Port: 80,
	})
	builder.AddSubscribers(nil)
	builder.AddContracts(nil)

	t.Run("subscriber creation requires authorization", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPost, "/subscribers", nil)

		builder.router.ServeHTTP(response, request)

		if response.Code != http.StatusUnauthorized {
			t.Fatalf("status = %d, want %d", response.Code, http.StatusUnauthorized)
		}
	})

	t.Run("last contract by object id allows internal calls without authorization", func(t *testing.T) {
		response := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/contracts/objects/1/last", nil)

		builder.router.ServeHTTP(response, request)

		if response.Code == http.StatusUnauthorized {
			t.Fatalf("status = %d, route must stay open for internal service calls", response.Code)
		}
	})
}
