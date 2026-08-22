package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

type fakeAuthenticateSessionUseCase struct {
	session       domain.AdministratorSession
	err           error
	receivedToken string
}

func (f *fakeAuthenticateSessionUseCase) Execute(ctx context.Context, sessionToken string) (domain.AdministratorSession, error) {
	f.receivedToken = sessionToken
	return f.session, f.err
}

func TestRequireSessionWithoutCookie(t *testing.T) {
	t.Parallel()

	authenticate := &fakeAuthenticateSessionUseCase{err: domain.ErrInactiveAdministratorSession}
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })

	handler := middleware.RequireSession(authenticate)(next)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if called {
		t.Fatal("next handler must not be called without a valid session")
	}
	if authenticate.receivedToken != "" {
		t.Fatalf("expected an empty token to be passed through, got %q", authenticate.receivedToken)
	}
}

func TestRequireSessionWithInvalidCookie(t *testing.T) {
	t.Parallel()

	authenticate := &fakeAuthenticateSessionUseCase{err: domain.ErrInactiveAdministratorSession}
	called := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { called = true })

	handler := middleware.RequireSession(authenticate)(next)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil)
	req.AddCookie(&http.Cookie{Name: "akritas_session", Value: "expired-or-unknown-token"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusUnauthorized)
	}
	if called {
		t.Fatal("next handler must not be called for an invalid session")
	}
	if authenticate.receivedToken != "expired-or-unknown-token" {
		t.Fatalf("expected the cookie value to be passed through, got %q", authenticate.receivedToken)
	}
	if strings.Contains(rec.Body.String(), "expired-or-unknown-token") {
		t.Fatal("response must never echo the cookie value")
	}
}

func TestRequireSessionWithValidCookieInjectsSession(t *testing.T) {
	t.Parallel()

	session := domain.AdministratorSession{ID: uuid.New(), AdministratorID: uuid.New()}
	authenticate := &fakeAuthenticateSessionUseCase{session: session}

	var gotSession domain.AdministratorSession
	var gotOK bool
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSession, gotOK = middleware.SessionFromContext(r.Context())
		w.WriteHeader(http.StatusOK)
	})

	handler := middleware.RequireSession(authenticate)(next)
	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil)
	req.AddCookie(&http.Cookie{Name: "akritas_session", Value: "a-valid-token"})
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}
	if !gotOK {
		t.Fatal("SessionFromContext must find the session RequireSession injected")
	}
	if gotSession.ID != session.ID {
		t.Fatal("injected session must match the one AuthenticateSession returned")
	}
	if authenticate.receivedToken != "a-valid-token" {
		t.Fatalf("expected the cookie value to be passed through, got %q", authenticate.receivedToken)
	}
}
