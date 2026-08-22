package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	handlerauth "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/auth"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func TestGetSetupStatusReturnsOK(t *testing.T) {
	t.Parallel()

	getStatus := &fakeGetSetupStatusUseCase{status: in.SetupStatus{SetupRequired: true, RegistrationOpen: true}}
	handler := handlerauth.NewHandler(getStatus, &fakeStartAdministratorSetupUseCase{}, &fakeVerifyAdministratorSetupUseCase{}, true)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/setup-status", nil)
	rec := httptest.NewRecorder()

	handler.GetSetupStatus(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusOK)
	}

	var body struct {
		Data struct {
			SetupRequired    bool `json:"setup_required"`
			RegistrationOpen bool `json:"registration_open"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if !body.Data.SetupRequired || !body.Data.RegistrationOpen {
		t.Fatalf("unexpected body: %+v", body)
	}
}

func TestGetSetupStatusReturnsInternalErrorOnFailure(t *testing.T) {
	t.Parallel()

	getStatus := &fakeGetSetupStatusUseCase{err: errUnexpected}
	handler := handlerauth.NewHandler(getStatus, &fakeStartAdministratorSetupUseCase{}, &fakeVerifyAdministratorSetupUseCase{}, true)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/auth/setup-status", nil)
	rec := httptest.NewRecorder()

	handler.GetSetupStatus(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusInternalServerError)
	}
}
