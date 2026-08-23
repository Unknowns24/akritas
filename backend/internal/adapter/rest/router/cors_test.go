package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

const testFrontendOrigin = "http://localhost:3000"

func TestRouterAddsCredentialedCORSHeadersToAllowedOriginErrors(t *testing.T) {
	fixture := newRouterFixture()
	fixture.config.AllowedOrigins = []string{testFrontendOrigin}
	fixture.authenticate.err = errInactiveSession()

	request := httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil)
	request.Header.Set("Origin", testFrontendOrigin)
	recorder := httptest.NewRecorder()
	mustRouter(t, fixture.config).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", recorder.Code)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != testFrontendOrigin {
		t.Fatalf("Access-Control-Allow-Origin = %q, want %q", got, testFrontendOrigin)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Credentials"); got != "true" {
		t.Fatalf("Access-Control-Allow-Credentials = %q, want true", got)
	}
	if vary := strings.Join(recorder.Header().Values("Vary"), ","); !strings.Contains(vary, "Origin") {
		t.Fatalf("Vary = %q, want Origin", vary)
	}
}

func TestRouterHandlesAllowedCORSPreflightAtTopLevel(t *testing.T) {
	config := validConfig()
	config.AllowedOrigins = []string{testFrontendOrigin}
	request := httptest.NewRequest(http.MethodOptions, "/api/v1/projects/"+uuid.NewString()+"/monitoring-configuration", nil)
	request.Header.Set("Origin", testFrontendOrigin)
	request.Header.Set("Access-Control-Request-Method", http.MethodPut)
	request.Header.Set("Access-Control-Request-Headers", "Content-Type, Idempotency-Key")
	recorder := httptest.NewRecorder()

	mustRouter(t, config).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", recorder.Code)
	}
	for name, want := range map[string]string{
		"Access-Control-Allow-Origin":      testFrontendOrigin,
		"Access-Control-Allow-Credentials": "true",
		"Access-Control-Allow-Methods":     http.MethodPut,
		"Access-Control-Max-Age":           "300",
	} {
		if got := recorder.Header().Get(name); got != want {
			t.Fatalf("%s = %q, want %q", name, got, want)
		}
	}
	allowedHeaders := recorder.Header().Get("Access-Control-Allow-Headers")
	for _, name := range []string{"Content-Type", "Idempotency-Key"} {
		if !strings.Contains(allowedHeaders, name) {
			t.Fatalf("Access-Control-Allow-Headers = %q, want %s", allowedHeaders, name)
		}
	}
}

func TestRouterDoesNotAuthorizeUnconfiguredCORSOrigin(t *testing.T) {
	fixture := newRouterFixture()
	fixture.config.AllowedOrigins = []string{testFrontendOrigin}
	fixture.authenticate.err = errInactiveSession()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil)
	request.Header.Set("Origin", "https://untrusted.example.com")
	recorder := httptest.NewRecorder()

	mustRouter(t, fixture.config).ServeHTTP(recorder, request)

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", recorder.Code)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != "" {
		t.Fatalf("untrusted Access-Control-Allow-Origin = %q, want empty", got)
	}
	if got := recorder.Header().Get("Access-Control-Allow-Credentials"); got != "" {
		t.Fatalf("untrusted Access-Control-Allow-Credentials = %q, want empty", got)
	}
}

func TestRouterPreservesRequestsWithoutOrigin(t *testing.T) {
	fixture := newRouterFixture()
	fixture.authenticate.err = errInactiveSession()
	recorder := serve(mustRouter(t, fixture.config), http.MethodGet, "/api/v1/auth/session")

	if recorder.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", recorder.Code)
	}
	if recorder.Header().Get("Access-Control-Allow-Origin") != "" || recorder.Header().Get("Access-Control-Allow-Credentials") != "" {
		t.Fatal("request without Origin received CORS authorization headers")
	}
}

func TestRouterReflectsOnlyTheExactConfiguredOrigin(t *testing.T) {
	fixture := newRouterFixture()
	fixture.config.AllowedOrigins = []string{"https://app.example.com", testFrontendOrigin}
	fixture.authenticate.err = errInactiveSession()
	request := httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil)
	request.Header.Set("Origin", testFrontendOrigin)
	recorder := httptest.NewRecorder()

	mustRouter(t, fixture.config).ServeHTTP(recorder, request)

	if got := recorder.Header().Get("Access-Control-Allow-Origin"); got != testFrontendOrigin || got == "*" {
		t.Fatalf("Access-Control-Allow-Origin = %q, want exact %q", got, testFrontendOrigin)
	}
}

func TestRouterKeepsCSRFOriginProtectionWithCORS(t *testing.T) {
	fixture := newRouterFixture()
	fixture.config.AllowedOrigins = []string{testFrontendOrigin}
	accountID := uuid.NewString()

	allowed := httptest.NewRequest(http.MethodDelete, "/api/v1/integrations/github/accounts/"+accountID, nil)
	allowed.Header.Set("Origin", testFrontendOrigin)
	allowedRecorder := httptest.NewRecorder()
	mustRouter(t, fixture.config).ServeHTTP(allowedRecorder, allowed)
	if allowedRecorder.Code != http.StatusNoContent {
		t.Fatalf("allowed mutation status = %d, want 204", allowedRecorder.Code)
	}

	denied := httptest.NewRequest(http.MethodDelete, "/api/v1/integrations/github/accounts/"+accountID, nil)
	denied.Header.Set("Origin", "https://untrusted.example.com")
	deniedRecorder := httptest.NewRecorder()
	mustRouter(t, fixture.config).ServeHTTP(deniedRecorder, denied)
	if deniedRecorder.Code != http.StatusForbidden {
		t.Fatalf("denied mutation status = %d, want 403", deniedRecorder.Code)
	}
}

func errInactiveSession() error {
	return domain.ErrInactiveAdministratorSession
}
