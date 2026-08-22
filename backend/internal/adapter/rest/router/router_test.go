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

func TestRouterServesSetupStatus(t *testing.T) {
	t.Parallel()

	handler := router.New(router.Dependencies{
		Auth: handlerauth.NewHandler(fakeGetSetupStatusUseCase{}, fakeStartAdministratorSetupUseCase{}),
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/setup-status", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestRouterServesSetup(t *testing.T) {
	t.Parallel()

	handler := router.New(router.Dependencies{
		Auth: handlerauth.NewHandler(fakeGetSetupStatusUseCase{}, fakeStartAdministratorSetupUseCase{}),
	})

	body := `{"email":"admin@example.com","display_name":"Akritas Administrator","password":"a-long-password-from-a-password-manager","bootstrap_token":"deployment-bootstrap-secret-not-a-totp-seed"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/setup", strings.NewReader(body))
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
}
