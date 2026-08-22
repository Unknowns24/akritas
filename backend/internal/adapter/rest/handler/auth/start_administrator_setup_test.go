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

func validSetupBody() map[string]string {
	return map[string]string{
		"email":           "admin@example.com",
		"display_name":    "Akritas Administrator",
		"password":        "a-long-password-from-a-password-manager",
		"bootstrap_token": "deployment-bootstrap-secret-not-a-totp-seed",
	}
}

func newSetupRequest(t *testing.T, body any) *http.Request {
	t.Helper()
	payload, err := json.Marshal(body)
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/setup", bytes.NewReader(payload))
	req.RemoteAddr = "203.0.113.10:54321"
	return req
}

func TestStartAdministratorSetupHappyPath(t *testing.T) {
	t.Parallel()

	expiresAt := time.Date(2026, 8, 22, 12, 10, 0, 0, time.UTC)
	startSetup := &fakeStartAdministratorSetupUseCase{output: in.StartAdministratorSetupOutput{
		EnrollmentID:   uuid.New(),
		OtpauthURI:     "otpauth://totp/Akritas:admin@example.com?secret=ABC&issuer=Akritas",
		ManualEntryKey: "JBSWY3DPEHPK3PXP",
		ExpiresAt:      expiresAt,
	}}
	handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, startSetup)

	rec := httptest.NewRecorder()
	handler.StartAdministratorSetup(rec, newSetupRequest(t, validSetupBody()))

	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusCreated, rec.Body.String())
	}
	if got := rec.Header().Get("Cache-Control"); got != "no-store" {
		t.Fatalf("Cache-Control = %q, want no-store", got)
	}

	var body struct {
		Data struct {
			EnrollmentID   string `json:"enrollment_id"`
			OtpauthURI     string `json:"otpauth_uri"`
			ManualEntryKey string `json:"manual_entry_key"`
			ExpiresAt      string `json:"expires_at"`
		} `json:"data"`
	}
	if err := json.NewDecoder(rec.Body).Decode(&body); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if body.Data.EnrollmentID != startSetup.output.EnrollmentID.String() ||
		body.Data.OtpauthURI != startSetup.output.OtpauthURI ||
		body.Data.ManualEntryKey != startSetup.output.ManualEntryKey ||
		body.Data.ExpiresAt != expiresAt.Format(time.RFC3339) {
		t.Fatalf("unexpected body: %+v", body)
	}

	if startSetup.receivedArgs.RateLimitKey != "203.0.113.10" {
		t.Fatalf("rate limit key = %q, want the request IP without port", startSetup.receivedArgs.RateLimitKey)
	}
}

func TestStartAdministratorSetupMalformedJSON(t *testing.T) {
	t.Parallel()

	startSetup := &fakeStartAdministratorSetupUseCase{}
	handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, startSetup)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/setup", strings.NewReader("{not json"))
	rec := httptest.NewRecorder()

	handler.StartAdministratorSetup(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestStartAdministratorSetupFieldValidation(t *testing.T) {
	t.Parallel()

	cases := map[string]map[string]string{
		"short password": {
			"email": "admin@example.com", "display_name": "Akritas Administrator",
			"password": "too-short", "bootstrap_token": strings.Repeat("a", 32),
		},
		"invalid email": {
			"email": "not-an-email", "display_name": "Akritas Administrator",
			"password": "a-long-password-from-a-password-manager", "bootstrap_token": strings.Repeat("a", 32),
		},
		"blank display name": {
			"email": "admin@example.com", "display_name": "",
			"password": "a-long-password-from-a-password-manager", "bootstrap_token": strings.Repeat("a", 32),
		},
		"short bootstrap token": {
			"email": "admin@example.com", "display_name": "Akritas Administrator",
			"password": "a-long-password-from-a-password-manager", "bootstrap_token": "short",
		},
	}

	for name, body := range cases {
		body := body
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			startSetup := &fakeStartAdministratorSetupUseCase{}
			handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, startSetup)

			rec := httptest.NewRecorder()
			handler.StartAdministratorSetup(rec, newSetupRequest(t, body))

			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want %d, body=%s", rec.Code, http.StatusBadRequest, rec.Body.String())
			}
			if startSetup.receivedArgs != (in.StartAdministratorSetupInput{}) {
				t.Fatal("usecase must not be called when the request fails shape validation")
			}
		})
	}
}

func TestStartAdministratorSetupUsecaseErrors(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		err            error
		wantStatus     int
		wantRetryAfter bool
	}{
		{"invalid bootstrap token", domain.ErrInvalidBootstrapToken, http.StatusBadRequest, false},
		{"administrator already exists", domain.ErrAdministratorAlreadyExists, http.StatusConflict, false},
		{"rate limited", authusecase.ErrSetupRateLimited, http.StatusTooManyRequests, true},
		{"unexpected error", errUnexpected, http.StatusInternalServerError, false},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			startSetup := &fakeStartAdministratorSetupUseCase{err: tc.err}
			handler := handlerauth.NewHandler(&fakeGetSetupStatusUseCase{}, startSetup)

			rec := httptest.NewRecorder()
			handler.StartAdministratorSetup(rec, newSetupRequest(t, validSetupBody()))

			if rec.Code != tc.wantStatus {
				t.Fatalf("status = %d, want %d, body=%s", rec.Code, tc.wantStatus, rec.Body.String())
			}
			if tc.wantRetryAfter && rec.Header().Get("Retry-After") == "" {
				t.Fatal("expected a Retry-After header")
			}
			if rec.Body.String() != "" && strings.Contains(rec.Body.String(), "deployment-bootstrap-secret-not-a-totp-seed") {
				t.Fatal("response must never echo the bootstrap token")
			}
			if rec.Body.String() != "" && strings.Contains(rec.Body.String(), "a-long-password-from-a-password-manager") {
				t.Fatal("response must never echo the password")
			}
		})
	}
}
