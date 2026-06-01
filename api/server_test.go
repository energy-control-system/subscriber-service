package api

import (
	"net/http"
	"net/http/httptest"
	"subscriber-service/config"
	"testing"

	"github.com/sunshineOfficial/golib/golog"
)

func TestSubscriberRoutesAllowUnauthenticatedRequests(t *testing.T) {
	builder := NewServerBuilder(t.Context(), golog.NewLogger("test"), config.Settings{
		Port: 80,
	})
	builder.AddSubscribers(nil)
	builder.AddObjects(nil)
	builder.AddContracts(nil)
	builder.AddRegistry(nil)

	routes := []struct {
		method string
		path   string
	}{
		{method: http.MethodPost, path: "/subscribers"},
		{method: http.MethodGet, path: "/subscribers/1/extended"},
		{method: http.MethodGet, path: "/subscribers/1"},
		{method: http.MethodPatch, path: "/subscribers/1"},
		{method: http.MethodDelete, path: "/subscribers/1"},
		{method: http.MethodGet, path: "/subscribers"},
		{method: http.MethodPost, path: "/objects"},
		{method: http.MethodGet, path: "/objects/1"},
		{method: http.MethodPatch, path: "/objects/1"},
		{method: http.MethodDelete, path: "/objects/1"},
		{method: http.MethodGet, path: "/objects/devices/1"},
		{method: http.MethodGet, path: "/objects/seals/1"},
		{method: http.MethodGet, path: "/objects"},
		{method: http.MethodPost, path: "/contracts"},
		{method: http.MethodGet, path: "/contracts"},
		{method: http.MethodPatch, path: "/contracts/1"},
		{method: http.MethodDelete, path: "/contracts/1"},
		{method: http.MethodGet, path: "/contracts/objects/last?id=1"},
		{method: http.MethodGet, path: "/contracts/objects/1/last"},
		{method: http.MethodPost, path: "/registry/parse"},
	}

	for _, route := range routes {
		t.Run(route.method+" "+route.path, func(t *testing.T) {
			response := httptest.NewRecorder()
			request := httptest.NewRequest(route.method, route.path, nil)

			builder.router.ServeHTTP(response, request)

			if response.Code == http.StatusUnauthorized {
				t.Fatalf("status = %d, route must be open without authorization", response.Code)
			}
		})
	}
}
