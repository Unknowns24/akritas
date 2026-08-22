package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
)

func TestRequireAllowedOriginRejectsUnsafeRequestsFailClosed(t *testing.T) {
	t.Parallel()
	wrapped := middleware.RequireAllowedOrigin([]string{"https://app.example.com"})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	for _, origin := range []string{"", "*", "https://evil.example.com", "https://app.example.com/"} {
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodDelete, "/api/v1/auth/session", nil)
		request.Header.Set("Origin", origin)
		wrapped.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusForbidden {
			t.Fatalf("origin %q: status=%d, want 403", origin, recorder.Code)
		}
	}
}

func TestRequireAllowedOriginAcceptsExactOriginAndSafeMethods(t *testing.T) {
	t.Parallel()
	wrapped := middleware.RequireAllowedOrigin([]string{"https://app.example.com"})(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	request := httptest.NewRequest(http.MethodDelete, "/api/v1/auth/session", nil)
	request.Header.Set("Origin", "https://app.example.com")
	recorder := httptest.NewRecorder()
	wrapped.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("allowed origin status=%d, want 204", recorder.Code)
	}

	recorder = httptest.NewRecorder()
	wrapped.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil))
	if recorder.Code != http.StatusNoContent {
		t.Fatalf("safe request status=%d, want 204", recorder.Code)
	}
}
