package errors

import "github.com/Unknowns24/akritas/backend/internal/core/domain"

var (
	ErrInvalidRequest = &domain.Error{Code: "0x102001V", Message: "invalid request", UserMessage: "La solicitud contiene datos inválidos."}
	ErrRequestFailed  = &domain.Error{Code: "0x102002I", Message: "internal request failure", UserMessage: "No se pudo completar la solicitud."}
)

func Catalog() map[string]*domain.Error {
	return map[string]*domain.Error{
		"ErrInvalidRequest": ErrInvalidRequest,
		"ErrRequestFailed":  ErrRequestFailed,
	}
}
