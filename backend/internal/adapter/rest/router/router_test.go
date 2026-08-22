package router_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	handlerauth "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/auth"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/router"
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

func newTestHandler() *handlerauth.Handler {
	return handlerauth.NewHandler(
		fakeGetSetupStatusUseCase{}, fakeStartAdministratorSetupUseCase{}, fakeVerifyAdministratorSetupUseCase{}, true,
	)
}

func TestRouterServesSetupStatus(t *testing.T) {
	t.Parallel()

	handler := router.New(router.Dependencies{Auth: newTestHandler()})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/setup-status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestRouterServesSetup(t *testing.T) {
	t.Parallel()

	handler := router.New(router.Dependencies{Auth: newTestHandler()})

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

	handler := router.New(router.Dependencies{Auth: newTestHandler()})

	body := `{"enrollment_id":"` + "00000000-0000-0000-0000-000000000000" + `","totp_code":"123456"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/setup/verify", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}
