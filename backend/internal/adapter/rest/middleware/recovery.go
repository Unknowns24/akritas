package middleware

import (
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

// RecoverPanics keeps unexpected transport failures inside the stable REST
// error contract without exposing panic values or stack traces to clients.
func RecoverPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				if recovered == http.ErrAbortHandler {
					panic(recovered)
				}
				response.WriteInternalError(w, r)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
