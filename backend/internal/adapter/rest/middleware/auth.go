// Package middleware holds cross-cutting REST adapters -- so far, session
// authentication.
package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

type sessionContextKey struct{}

// RequireSession resolves the session cookie via authenticate, injecting
// the active domain.AdministratorSession into the request context on
// success, or writing the appropriate error response (and never calling
// next) on failure.
func RequireSession(authenticate in.AuthenticateSessionUseCase) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := ""
			if cookie, err := r.Cookie(response.SessionCookieName); err == nil {
				token = cookie.Value
			}

			session, err := authenticate.Execute(r.Context(), token)
			if err != nil {
				var domainErr *domain.Error
				if errors.As(err, &domainErr) {
					response.WriteDomainError(w, r, domainErr)
					return
				}
				response.WriteInternalError(w, r)
				return
			}

			ctx := context.WithValue(r.Context(), sessionContextKey{}, session)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// SessionFromContext returns the session RequireSession injected, if any.
func SessionFromContext(ctx context.Context) (domain.AdministratorSession, bool) {
	session, ok := ctx.Value(sessionContextKey{}).(domain.AdministratorSession)
	return session, ok
}
