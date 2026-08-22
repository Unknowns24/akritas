package auth_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"

	handlerauth "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/auth"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

// withInjectedSession routes req through the real RequireSession middleware
// (with a fake that always succeeds) so the request carries a session in
// its context exactly the way the router would produce it -- the injection
// key itself is unexported inside the middleware package.
func withInjectedSession(req *http.Request, session domain.AdministratorSession) *http.Request {
	var injected *http.Request
	shim := middleware.RequireSession(fakeAlwaysAuthenticate{session: session})(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { injected = r }),
	)
	rec := httptest.NewRecorder()
	shim.ServeHTTP(rec, req)
	return injected
}

type fakeAlwaysAuthenticate struct {
	session domain.AdministratorSession
}

func (f fakeAlwaysAuthenticate) Execute(ctx context.Context, sessionToken string) (domain.AdministratorSession, error) {
	return f.session, nil
}

func TestGetCurrentSessionHappyPath(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	adminID := uuid.New()
	getCurrentSession := &fakeGetCurrentSessionUseCase{output: in.GetCurrentSessionOutput{
		Administrator: domain.Administrator{
			ID: adminID, Email: "admin@example.com", DisplayName: "Akritas Administrator",
			CreatedAt: now, UpdatedAt: now,
		},
		AuthenticatedAt:   now,
		IdleExpiresAt:     now.Add(12 * time.Hour),
		AbsoluteExpiresAt: now.Add(168 * time.Hour),
	}}
	handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, &fakeVerifyAdministratorSetupUseCase{}, &fakeLoginAdministratorUseCase{}, getCurrentSession, &fakeLogoutAdministratorUseCase{}, true)

	session := domain.AdministratorSession{ID: uuid.New(), AdministratorID: adminID}
	req := withInjectedSession(httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil), session)
	rec := httptest.NewRecorder()

	handler.GetCurrentSession(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if getCurrentSession.receivedInput.ID != session.ID {
		t.Fatal("usecase must receive the session injected by the middleware")
	}

	var body struct {
		Data struct {
			Administrator struct {
				ID string `json:"id"`
			} `json:"administrator"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Administrator.ID != adminID.String() {
		t.Fatalf("unexpected administrator id in body: %q", body.Data.Administrator.ID)
	}
}

func TestGetCurrentSessionUsecaseErrorMapsToDomainStatus(t *testing.T) {
	t.Parallel()

	getCurrentSession := &fakeGetCurrentSessionUseCase{err: domain.ErrInactiveAdministratorSession}
	handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, &fakeVerifyAdministratorSetupUseCase{}, &fakeLoginAdministratorUseCase{}, getCurrentSession, &fakeLogoutAdministratorUseCase{}, true)

	req := withInjectedSession(httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil), domain.AdministratorSession{ID: uuid.New()})
	rec := httptest.NewRecorder()

	handler.GetCurrentSession(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusUnauthorized, rec.Body.String())
	}
}

func TestGetCurrentSessionUnexpectedErrorMapsTo500(t *testing.T) {
	t.Parallel()

	getCurrentSession := &fakeGetCurrentSessionUseCase{err: errUnexpected}
	handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, &fakeVerifyAdministratorSetupUseCase{}, &fakeLoginAdministratorUseCase{}, getCurrentSession, &fakeLogoutAdministratorUseCase{}, true)

	req := withInjectedSession(httptest.NewRequest(http.MethodGet, "/api/v1/auth/session", nil), domain.AdministratorSession{ID: uuid.New()})
	rec := httptest.NewRecorder()

	handler.GetCurrentSession(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusInternalServerError, rec.Body.String())
	}
}
