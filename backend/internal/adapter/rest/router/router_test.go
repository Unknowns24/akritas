package router_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"

	handlerauth "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/auth"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/router"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

type fakeGetSetupStatusUseCase struct{}

func (fakeGetSetupStatusUseCase) Execute(ctx context.Context) (in.SetupStatus, error) {
	return in.SetupStatus{SetupRequired: true, RegistrationOpen: true}, nil
}

type fakeStartAdministratorSetupUseCase struct{}

func (fakeStartAdministratorSetupUseCase) Execute(ctx context.Context, input in.StartAdministratorSetupInput) (in.StartAdministratorSetupOutput, error) {
	return in.StartAdministratorSetupOutput{}, nil
}

type fakeVerifyAdministratorSetupUseCase struct{}

func (fakeVerifyAdministratorSetupUseCase) Execute(ctx context.Context, input in.VerifyAdministratorSetupInput) (in.VerifyAdministratorSetupOutput, error) {
	return in.VerifyAdministratorSetupOutput{}, nil
}

type fakeLoginAdministratorUseCase struct{}

func (fakeLoginAdministratorUseCase) Execute(ctx context.Context, input in.LoginAdministratorInput) (in.LoginAdministratorOutput, error) {
	return in.LoginAdministratorOutput{}, nil
}

type fakeGetCurrentSessionUseCase struct{}

func (fakeGetCurrentSessionUseCase) Execute(ctx context.Context, session domain.AdministratorSession) (in.GetCurrentSessionOutput, error) {
	return in.GetCurrentSessionOutput{}, nil
}

type fakeLogoutAdministratorUseCase struct{}

func (fakeLogoutAdministratorUseCase) Execute(ctx context.Context, session domain.AdministratorSession) error {
	return nil
}

// fakeAuthenticateSessionUseCase drives the router's real RequireSession
// middleware end to end: succeed to let requests through, or return err to
// exercise the 401 path.
type fakeAuthenticateSessionUseCase struct {
	session domain.AdministratorSession
	err     error
}

func (f fakeAuthenticateSessionUseCase) Execute(ctx context.Context, sessionToken string) (domain.AdministratorSession, error) {
	return f.session, f.err
}

func newTestHandler() *handlerauth.Handler {
	return handlerauth.NewHandler(
		fakeGetSetupStatusUseCase{}, fakeStartAdministratorSetupUseCase{}, fakeVerifyAdministratorSetupUseCase{},
		fakeLoginAdministratorUseCase{}, fakeGetCurrentSessionUseCase{}, fakeLogoutAdministratorUseCase{}, true,
	)
}

func TestRouterServesSetupStatus(t *testing.T) {
	t.Parallel()

	handler := router.New(router.Dependencies{Auth: newTestHandler(), Authenticate: fakeAuthenticateSessionUseCase{}})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/setup-status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestRouterServesSetup(t *testing.T) {
	t.Parallel()

	handler := router.New(router.Dependencies{Auth: newTestHandler(), Authenticate: fakeAuthenticateSessionUseCase{}})

	body := `{"email":"admin@example.com","display_name":"Akritas Administrator","password":"a-long-password-from-a-password-manager","bootstrap_token":"deployment-bootstrap-secret-not-a-totp-seed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/setup", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
}

func TestRouterServesSetupVerify(t *testing.T) {
	t.Parallel()

	handler := router.New(router.Dependencies{Auth: newTestHandler(), Authenticate: fakeAuthenticateSessionUseCase{}})

	body := `{"enrollment_id":"` + "00000000-0000-0000-0000-000000000000" + `","totp_code":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/setup/verify", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestRouterServesLogin(t *testing.T) {
	t.Parallel()

	handler := router.New(router.Dependencies{Auth: newTestHandler(), Authenticate: fakeAuthenticateSessionUseCase{}})

	body := `{"email":"admin@example.com","password":"a-long-password-from-a-password-manager","totp_code":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestRouterRequiresSessionForGetSession(t *testing.T) {
	t.Parallel()

	handler := router.New(router.Dependencies{
		Auth:         newTestHandler(),
		Authenticate: fakeAuthenticateSessionUseCase{err: domain.ErrInactiveAdministratorSession},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusUnauthorized, rec.Body.String())
	}
}

func TestRouterServesGetSessionWhenAuthenticated(t *testing.T) {
	t.Parallel()

	session := domain.AdministratorSession{ID: uuid.New(), AdministratorID: uuid.New()}
	handler := router.New(router.Dependencies{
		Auth:         newTestHandler(),
		Authenticate: fakeAuthenticateSessionUseCase{session: session},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil)
	req.AddCookie(&http.Cookie{Name: "akritas_session", Value: "raw-token"})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestRouterServesLogoutWhenAuthenticated(t *testing.T) {
	t.Parallel()

	session := domain.AdministratorSession{ID: uuid.New(), AdministratorID: uuid.New()}
	handler := router.New(router.Dependencies{
		Auth:         newTestHandler(),
		Authenticate: fakeAuthenticateSessionUseCase{session: session},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/auth/session", nil)
	req.AddCookie(&http.Cookie{Name: "akritas_session", Value: "raw-token"})
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusNoContent, rec.Body.String())
	}
}

func TestRouterRequiresSessionForLogout(t *testing.T) {
	t.Parallel()

	handler := router.New(router.Dependencies{
		Auth:         newTestHandler(),
		Authenticate: fakeAuthenticateSessionUseCase{err: domain.ErrInactiveAdministratorSession},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/auth/session", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusUnauthorized, rec.Body.String())
	}
}
