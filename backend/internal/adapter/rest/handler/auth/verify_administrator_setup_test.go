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
)

func validVerifyBody() map[string]string {
	return map[string]string{
		"enrollment_id": uuid.New().String(),
		"totp_code":     "123456",
	}
}

func newVerifyRequest(t *testing.T, body any) *http.Request {
	t.Helper()
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	return httptest.NewRequest(http.MethodPost, "/api/v1/auth/setup/verify", bytes.NewReader(payload))
}

func TestVerifyAdministratorSetupHappyPath(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	adminID := uuid.New()
	verify := &fakeVerifyAdministratorSetupUseCase{output: in.VerifyAdministratorSetupOutput{
		Administrator: domain.Administrator{
			ID: adminID, Email: "admin@example.com", DisplayName: "Akritas Administrator",
			CreatedAt: now, UpdatedAt: now,
		},
		SessionToken:      "raw-session-token",
		AuthenticatedAt:   now,
		IdleExpiresAt:     now.Add(12 * time.Hour),
		AbsoluteExpiresAt: now.Add(168 * time.Hour),
	}}
	handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, verify, &fakeLoginAdministratorUseCase{}, &fakeGetCurrentSessionUseCase{}, &fakeLogoutAdministratorUseCase{}, true)

	rec := httptest.NewRecorder()
	handler.VerifyAdministratorSetup(rec, newVerifyRequest(t, validVerifyBody()))

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusOK, rec.Body.String())
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected exactly 1 cookie, got %d", len(cookies))
	}
	cookie := cookies[0]
	if cookie.Name != "akritas_session" {
		t.Fatalf("cookie name = %q, want akritas_session", cookie.Name)
	}
	if cookie.Value != "raw-session-token" {
		t.Fatal("cookie value must be the raw session token")
	}
	if !cookie.HttpOnly {
		t.Fatal("cookie must be HttpOnly")
	}
	if !cookie.Secure {
		t.Fatal("cookie must be Secure when the handler is configured for it")
	}
	if cookie.SameSite != http.SameSiteLaxMode {
		t.Fatal("cookie must use SameSite=Lax")
	}
	if cookie.Path != "/" {
		t.Fatalf("cookie path = %q, want /", cookie.Path)
	}

	var body struct {
		Data struct {
			Administrator struct {
				ID          string `json:"id"`
				Email       string `json:"email"`
				DisplayName string `json:"display_name"`
				CreatedAt   string `json:"created_at"`
				UpdatedAt   string `json:"updated_at"`
			} `json:"administrator"`
			AuthenticatedAt   string `json:"authenticated_at"`
			IdleExpiresAt     string `json:"idle_expires_at"`
			AbsoluteExpiresAt string `json:"absolute_expires_at"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.Administrator.ID != adminID.String() || body.Data.Administrator.Email != "admin@example.com" {
		t.Fatalf("unexpected administrator in body: %+v", body.Data.Administrator)
	}
	if body.Data.AuthenticatedAt != now.Format(time.RFC3339) {
		t.Fatalf("authenticated_at = %q, want %q", body.Data.AuthenticatedAt, now.Format(time.RFC3339))
	}

	if strings.Contains(rec.Body.String(), "raw-session-token") {
		t.Fatal("the raw session token must never appear in the JSON body")
	}
}

func TestVerifyAdministratorSetupCookieNotSecureWhenConfigured(t *testing.T) {
	t.Parallel()

	verify := &fakeVerifyAdministratorSetupUseCase{}
	handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, verify, &fakeLoginAdministratorUseCase{}, &fakeGetCurrentSessionUseCase{}, &fakeLogoutAdministratorUseCase{}, false)

	rec := httptest.NewRecorder()
	handler.VerifyAdministratorSetup(rec, newVerifyRequest(t, validVerifyBody()))

	cookies := rec.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Secure {
		t.Fatal("cookie must not be Secure when the handler is configured without it")
	}
}

func TestVerifyAdministratorSetupMalformedJSON(t *testing.T) {
	t.Parallel()

	handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, &fakeVerifyAdministratorSetupUseCase{}, &fakeLoginAdministratorUseCase{}, &fakeGetCurrentSessionUseCase{}, &fakeLogoutAdministratorUseCase{}, true)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/setup/verify", strings.NewReader("{not json"))
	rec := httptest.NewRecorder()

	handler.VerifyAdministratorSetup(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestVerifyAdministratorSetupFieldValidation(t *testing.T) {
	t.Parallel()

	cases := map[string]map[string]string{
		"malformed enrollment_id": {"enrollment_id": "not-a-uuid", "totp_code": "123456"},
		"short totp_code":         {"enrollment_id": uuid.New().String(), "totp_code": "123"},
		"non numeric totp_code":   {"enrollment_id": uuid.New().String(), "totp_code": "abcdef"},
	}

	for name, body := range cases {
		body := body
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			verify := &fakeVerifyAdministratorSetupUseCase{}
			handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, verify, &fakeLoginAdministratorUseCase{}, &fakeGetCurrentSessionUseCase{}, &fakeLogoutAdministratorUseCase{}, true)

			rec := httptest.NewRecorder()
			handler.VerifyAdministratorSetup(rec, newVerifyRequest(t, body))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
			}
			if verify.receivedArgs != (in.VerifyAdministratorSetupInput{}) {
				t.Fatal("usecase must not be called when the request fails shape validation")
			}
		})
	}
}

func TestVerifyAdministratorSetupUsecaseErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		err            error
		wantStatus     int
		wantRetryAfter bool
	}{
		{"invalid verification", domain.ErrInvalidTotpEnrollmentVerification, http.StatusBadRequest, false},
		{"administrator already exists", domain.ErrAdministratorAlreadyExists, http.StatusConflict, false},
		{"rate limited", domain.ErrAuthenticationRateLimited, http.StatusTooManyRequests, true},
		{"unexpected error", errUnexpected, http.StatusInternalServerError, false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			verify := &fakeVerifyAdministratorSetupUseCase{err: tc.err}
			handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, verify, &fakeLoginAdministratorUseCase{}, &fakeGetCurrentSessionUseCase{}, &fakeLogoutAdministratorUseCase{}, true)

			body := validVerifyBody()
			rec := httptest.NewRecorder()
			handler.VerifyAdministratorSetup(rec, newVerifyRequest(t, body))

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d, body=%s", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if tc.wantRetryAfter && rec.Header().Get("Retry-After") == "" {
				t.Fatal("expected a Retry-After header")
			}
			if len(rec.Result().Cookies()) != 0 {
				t.Fatal("no cookie must be set when verification fails")
			}
			if strings.Contains(rec.Body.String(), body["totp_code"]) {
				t.Fatal("response must never echo the submitted totp_code")
			}
		})
	}
}
