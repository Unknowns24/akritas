package auth_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	handlerauth "github.com/Unknowns24/akritas/backend/internal/adapter/rest/handler/auth"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	authusecase "github.com/Unknowns24/akritas/backend/internal/usecase/auth"
)

func validLoginBody() map[string]string {
	return map[string]string{
		"email":     "admin@example.com",
		"password":  "a-long-password-from-a-password-manager",
		"totp_code": "123456",
	}
}

func newLoginRequest(t *testing.T, body any) *http.Request {
	t.Helper()
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(payload))
	req.RemoteAddr = "203.0.113.10:54321"
	return req
}

func TestLoginHappyPath(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	adminID := uuid.New()
	login := &fakeLoginAdministratorUseCase{output: in.LoginAdministratorOutput{
		Administrator: domain.Administrator{
			ID: adminID, Email: "admin@example.com", DisplayName: "Akritas Administrator",
			CreatedAt: now, UpdatedAt: now,
		},
		SessionToken:      "raw-session-token",
		AuthenticatedAt:   now,
		IdleExpiresAt:     now.Add(12 * time.Hour),
		AbsoluteExpiresAt: now.Add(168 * time.Hour),
	}}
	handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, &fakeVerifyAdministratorSetupUseCase{}, login, &fakeGetCurrentSessionUseCase{}, &fakeLogoutAdministratorUseCase{}, true)

	rec := httptest.NewRecorder()
	handler.Login(rec, newLoginRequest(t, validLoginBody()))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Value != "raw-session-token" {
		t.Fatal("expected exactly 1 cookie carrying the raw session token")
	}

	if login.receivedArgs.RateLimitKey != "203.0.113.10" {
		t.Fatalf("rate limit key = %q, want the request IP without port", login.receivedArgs.RateLimitKey)
	}
	if strings.Contains(rec.Body.String(), "raw-session-token") {
		t.Fatal("the raw session token must never appear in the JSON body")
	}
}

func TestLoginMalformedJSON(t *testing.T) {
	t.Parallel()

	handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, &fakeVerifyAdministratorSetupUseCase{}, &fakeLoginAdministratorUseCase{}, &fakeGetCurrentSessionUseCase{}, &fakeLogoutAdministratorUseCase{}, true)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", strings.NewReader("{not json"))
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestLoginFieldValidation(t *testing.T) {
	t.Parallel()

	cases := map[string]map[string]string{
		"invalid email":  {"email": "not-an-email", "password": "a-long-password-from-a-password-manager", "totp_code": "123456"},
		"short password": {"email": "admin@example.com", "password": "too-short", "totp_code": "123456"},
		"malformed totp": {"email": "admin@example.com", "password": "a-long-password-from-a-password-manager", "totp_code": "abc"},
	}

	for name, body := range cases {
		body := body
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			login := &fakeLoginAdministratorUseCase{}
			handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, &fakeVerifyAdministratorSetupUseCase{}, login, &fakeGetCurrentSessionUseCase{}, &fakeLogoutAdministratorUseCase{}, true)

			rec := httptest.NewRecorder()
			handler.Login(rec, newLoginRequest(t, body))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
			}
			if login.receivedArgs != (in.LoginAdministratorInput{}) {
				t.Fatal("usecase must not be called when the request fails shape validation")
			}
		})
	}
}

func TestLoginUsecaseErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		err            error
		wantStatus     int
		wantRetryAfter bool
	}{
		{"invalid credentials", domain.ErrInvalidCredentials, http.StatusUnauthorized, false},
		{"rate limited", authusecase.ErrLoginRateLimited, http.StatusTooManyRequests, true},
		{"unexpected error", errUnexpected, http.StatusInternalServerError, false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			login := &fakeLoginAdministratorUseCase{err: tc.err}
			handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, &fakeVerifyAdministratorSetupUseCase{}, login, &fakeGetCurrentSessionUseCase{}, &fakeLogoutAdministratorUseCase{}, true)

			body := validLoginBody()
			rec := httptest.NewRecorder()
			handler.Login(rec, newLoginRequest(t, body))

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d, body=%s", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if tc.wantRetryAfter && rec.Header().Get("Retry-After") == "" {
				t.Fatal("expected a Retry-After header")
			}
			if len(rec.Result().Cookies()) != 0 {
				t.Fatal("no cookie must be set when login fails")
			}
			if strings.Contains(rec.Body.String(), body["password"]) || strings.Contains(rec.Body.String(), body["totp_code"]) {
				t.Fatal("response must never echo the password or totp_code")
			}
		})
	}
}
