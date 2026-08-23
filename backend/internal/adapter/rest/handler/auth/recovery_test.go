package auth_test

import (
	"bytes"
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

func recoveryHandler(start *fakeStartAdministratorRecoveryUseCase, verify *fakeVerifyAdministratorRecoveryUseCase) *handlerauth.Handler {
	return handlerauth.NewHandlerWithRecovery(&fakeGetSetupStatusUseCase{}, &fakeStartAdministratorSetupUseCase{}, &fakeVerifyAdministratorSetupUseCase{}, &fakeLoginAdministratorUseCase{}, start, verify, &fakeGetCurrentSessionUseCase{}, &fakeLogoutAdministratorUseCase{}, true)
}

func TestStartRecoveryUsesRemotePeerAndReturnsNoStoreProvisioning(t *testing.T) {
	t.Parallel()
	expires := time.Date(2026, 8, 23, 12, 10, 0, 0, time.UTC)
	start := &fakeStartAdministratorRecoveryUseCase{output: in.StartAdministratorRecoveryOutput{EnrollmentID: uuid.New(), OtpauthURI: "otpauth://totp/Akritas", ManualEntryKey: "JBSWY3DPEHPK3PXP", ExpiresAt: expires}}
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/recovery", bytes.NewBufferString(`{"email":"admin@example.com","new_password":"a-new-long-password","bootstrap_token":"01234567890123456789012345678901"}`))
	req.RemoteAddr = "203.0.113.10:4321"
	req.Header.Set("X-Forwarded-For", "198.51.100.99")
	recorder := httptest.NewRecorder()
	recoveryHandler(start, &fakeVerifyAdministratorRecoveryUseCase{}).StartAdministratorRecovery(recorder, req)
	if recorder.Code != http.StatusCreated || recorder.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("status=%d headers=%v", recorder.Code, recorder.Header())
	}
	if start.receivedArgs.RateLimitKey != "203.0.113.10" {
		t.Fatalf("rate key = %q, must use direct peer", start.receivedArgs.RateLimitKey)
	}
	if strings.Contains(recorder.Body.String(), start.receivedArgs.NewPassword) || strings.Contains(recorder.Body.String(), start.receivedArgs.BootstrapToken) {
		t.Fatal("response leaked recovery credentials")
	}
}

func TestVerifyRecoveryCreatesHardenedCookie(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	administrator, _ := domain.NewAdministrator(uuid.New(), "admin@example.com", "Admin", now.Add(-time.Hour))
	verify := &fakeVerifyAdministratorRecoveryUseCase{output: in.VerifyAdministratorRecoveryOutput{Administrator: *administrator, SessionToken: "opaque-session", AuthenticatedAt: now, IdleExpiresAt: now.Add(time.Hour), AbsoluteExpiresAt: now.Add(7 * 24 * time.Hour)}}
	id := uuid.NewString()
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/recovery/verify", bytes.NewBufferString(`{"enrollment_id":"`+id+`","totp_code":"123456"}`))
	recorder := httptest.NewRecorder()
	recoveryHandler(&fakeStartAdministratorRecoveryUseCase{}, verify).VerifyAdministratorRecovery(recorder, req)
	if recorder.Code != http.StatusOK || recorder.Header().Get("Cache-Control") != "no-store" {
		t.Fatalf("status=%d headers=%v", recorder.Code, recorder.Header())
	}
	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 || !cookies[0].HttpOnly || !cookies[0].Secure || cookies[0].SameSite != http.SameSiteLaxMode || cookies[0].Path != "/" || cookies[0].MaxAge <= 0 || cookies[0].Value != "opaque-session" {
		t.Fatalf("cookie not hardened: %+v", cookies)
	}
	if strings.Contains(recorder.Body.String(), "opaque-session") || strings.Contains(recorder.Body.String(), "123456") {
		t.Fatal("response body leaked session or TOTP")
	}
}
