package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
)

func TestRecoverPanicsPreservesAbortHandler(t *testing.T) {
	handler := middleware.RecoverPanics(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic(http.ErrAbortHandler)
	}))

	defer func() {
		if recovered := recover(); recovered != http.ErrAbortHandler {
			t.Fatalf("recovered = %v, want http.ErrAbortHandler", recovered)
		}
	}()
	handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, "/", nil))
}
