package auth

import (
	"net/http"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

func setSessionCookie(w http.ResponseWriter, secure bool, token string, authenticatedAt, absoluteExpiresAt time.Time) {
	maxAge := int(absoluteExpiresAt.Sub(authenticatedAt).Seconds())
	if maxAge < 1 {
		maxAge = 1
	}
	http.SetCookie(w, &http.Cookie{Name: response.SessionCookieName, Value: token, Path: "/", HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode, Expires: absoluteExpiresAt, MaxAge: maxAge})
}

func expireSessionCookie(w http.ResponseWriter, secure bool) {
	http.SetCookie(w, &http.Cookie{Name: response.SessionCookieName, Value: "", Path: "/", HttpOnly: true, Secure: secure, SameSite: http.SameSiteLaxMode, MaxAge: -1, Expires: time.Unix(0, 0)})
}
