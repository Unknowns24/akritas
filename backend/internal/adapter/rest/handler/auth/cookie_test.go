package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSessionCookieUsesConfiguredSameSiteMode(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 23, 12, 0, 0, 0, time.UTC)
	recorder := httptest.NewRecorder()
	setSessionCookie(recorder, true, http.SameSiteStrictMode, "token", now, now.Add(time.Hour))

	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 || cookies[0].SameSite != http.SameSiteStrictMode {
		t.Fatalf("cookie SameSite = %+v, want Strict", cookies)
	}
}

func TestExpiredSessionCookieUsesConfiguredSameSiteMode(t *testing.T) {
	t.Parallel()

	recorder := httptest.NewRecorder()
	expireSessionCookie(recorder, true, http.SameSiteNoneMode)

	cookies := recorder.Result().Cookies()
	if len(cookies) != 1 || cookies[0].SameSite != http.SameSiteNoneMode {
		t.Fatalf("cookie SameSite = %+v, want None", cookies)
	}
}
