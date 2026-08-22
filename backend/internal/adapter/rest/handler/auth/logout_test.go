package auth_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"

	handlerauth "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/auth"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func TestLogoutHappyPath(t *testing.T) {
	t.Parallel()

	logout := &fakeLogoutAdministratorUseCase{}
	handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, &fakeVerifyAdministratorSetupUseCase{}, &fakeLoginAdministratorUseCase{}, &fakeGetCurrentSessionUseCase{}, logout, true)

	session := domain.AdministratorSession{ID: uuid.New(), AdministratorID: uuid.New()}
	req := withInjectedSession(httptest.NewRequest(http.MethodDelete, "/api/v1/auth/session", nil), session)
	rec := httptest.NewRecorder()

	handler.Logout(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusNoContent, rec.Body.String())
	}
	if logout.receivedInput.ID != session.ID {
		t.Fatal("usecase must receive the session injected by the middleware")
	}

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected exactly 1 cookie (the expiring one), got %d", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Value != "" {
		t.Fatalf("expiring cookie must carry an empty value, got %q", cookie.Value)
	}
	if cookie.MaxAge >= 0 {
		t.Fatalf("expiring cookie must have a negative Max-Age, got %d", cookie.MaxAge)
	}
}

func TestLogoutUsecaseErrorMapsToDomainStatus(t *testing.T) {
	t.Parallel()

	logout := &fakeLogoutAdministratorUseCase{err: domain.ErrAdministratorSessionTransition}
	handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, &fakeVerifyAdministratorSetupUseCase{}, &fakeLoginAdministratorUseCase{}, &fakeGetCurrentSessionUseCase{}, logout, true)

	req := withInjectedSession(httptest.NewRequest(http.MethodDelete, "/api/v1/auth/session", nil), domain.AdministratorSession{ID: uuid.New()})
	rec := httptest.NewRecorder()

	handler.Logout(rec, req)

	if rec.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusConflict, rec.Body.String())
	}
	if len(rec.Result().Cookies()) != 0 {
		t.Fatal("no cookie must be set when logout fails")
	}
}

func TestLogoutUnexpectedErrorMapsTo500(t *testing.T) {
	t.Parallel()

	logout := &fakeLogoutAdministratorUseCase{err: errUnexpected}
	handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, &fakeVerifyAdministratorSetupUseCase{}, &fakeLoginAdministratorUseCase{}, &fakeGetCurrentSessionUseCase{}, logout, true)

	req := withInjectedSession(httptest.NewRequest(http.MethodDelete, "/api/v1/auth/session", nil), domain.AdministratorSession{ID: uuid.New()})
	rec := httptest.NewRecorder()

	handler.Logout(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusInternalServerError, rec.Body.String())
	}
}
