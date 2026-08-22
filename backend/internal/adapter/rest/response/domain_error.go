package response

import (
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

// WriteDomainError maps a domain.Error to its HTTP status per the rule in
// docs/errors/aaa-map.md: V->400, U->401, F->403, N->404, C->409, I->500.
func WriteDomainError(w http.ResponseWriter, r *http.Request, err *domain.Error) {
	WriteError(w, r, domainErrorStatus(err), err.Code, err.Message, err.UserMessage)
}

func domainErrorStatus(err *domain.Error) int {
	if len(err.Code) == 0 {
		return http.StatusInternalServerError
	}
	switch err.Code[len(err.Code)-1] {
	case 'V':
		return http.StatusBadRequest
	case 'U':
		return http.StatusUnauthorized
	case 'F':
		return http.StatusForbidden
	case 'N':
		return http.StatusNotFound
	case 'C':
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
