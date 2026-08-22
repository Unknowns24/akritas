package middleware

import (
	"net/http"
	"strings"

	resterrors "github.com/Unknowns24/akritas/backend/internal/adapter/rest/errors"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

// RequireAllowedOrigin protects cookie-authenticated mutations from CSRF.
// Safe methods are unaffected; unsafe methods require an exact allowlisted
// Origin and deliberately reject missing or wildcard values.
func RequireAllowedOrigin(allowedOrigins []string) func(http.Handler) http.Handler {
	allowed := make(map[string]struct{}, len(allowedOrigins))
	for _, origin := range allowedOrigins {
		if normalized := strings.TrimSpace(origin); normalized != "" {
			allowed[normalized] = struct{}{}
		}
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
				next.ServeHTTP(w, r)
				return
			}
			if _, ok := allowed[strings.TrimSpace(r.Header.Get("Origin"))]; !ok {
				response.Error(w, r, resterrors.ErrOriginForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
