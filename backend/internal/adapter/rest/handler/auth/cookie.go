package auth

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

var ErrInvalidSessionCookieSameSite = errors.New("invalid session cookie SameSite mode")

func ParseSessionCookieSameSite(value string) (http.SameSite, error) {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "lax":
		return http.SameSiteLaxMode, nil
	case "strict":
		return http.SameSiteStrictMode, nil
	case "none":
		return http.SameSiteNoneMode, nil
	default:
		return http.SameSiteDefaultMode, ErrInvalidSessionCookieSameSite
	}
}

func setSessionCookie(w http.ResponseWriter, secure bool, sameSite http.SameSite, token string, authenticatedAt, absoluteExpiresAt time.Time) {
	maxAge := int(absoluteExpiresAt.Sub(authenticatedAt).Seconds())
	if maxAge < 1 {
		maxAge = 1
	}
	http.SetCookie(w, &http.Cookie{Name: response.SessionCookieName, Value: token, Path: "/", HttpOnly: true, Secure: secure, SameSite: sameSite, Expires: absoluteExpiresAt, MaxAge: maxAge})
}

func expireSessionCookie(w http.ResponseWriter, secure bool, sameSite http.SameSite) {
	http.SetCookie(w, &http.Cookie{Name: response.SessionCookieName, Value: "", Path: "/", HttpOnly: true, Secure: secure, SameSite: sameSite, MaxAge: -1, Expires: time.Unix(0, 0)})
}
