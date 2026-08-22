package middleware

import (
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/httperr"
	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
)

const SessionCookie = "akritas_session"

type SessionValidator interface {
	Valid(token string) bool
}

type RejectAllSessions struct{}

func (RejectAllSessions) Valid(string) bool { return false }

type AllowNonEmptySessions struct{}

func (AllowNonEmptySessions) Valid(token string) bool { return token != "" }

func RequireSession(validator SessionValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(SessionCookie)
			if err != nil || !validator.Valid(cookie.Value) {
				httperr.Write(w, apperr.ErrUnauthenticated, RequestID(r))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
