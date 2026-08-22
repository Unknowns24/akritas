package errors

import "github.com/Unknowns24/akritas/backend/internal/core/domain"

var (
	ErrInvalidRequest  = &domain.Error{Code: "0x102001V", Message: "invalid request", UserMessage: "La solicitud contiene datos inválidos."}
	ErrRequestFailed   = &domain.Error{Code: "0x102002I", Message: "internal request failure", UserMessage: "No se pudo completar la solicitud."}
	ErrRateLimited     = &domain.Error{Code: "0x102003R", Message: "rate limited", UserMessage: "Alcanzaste el límite de intentos. Probá nuevamente más tarde."}
	ErrOriginForbidden = &domain.Error{Code: "0x102004F", Message: "origin forbidden", UserMessage: "El origen de la solicitud no está permitido."}
)

func Catalog() map[string]*domain.Error {
	return map[string]*domain.Error{
		"ErrInvalidRequest":  ErrInvalidRequest,
		"ErrRequestFailed":   ErrRequestFailed,
		"ErrRateLimited":     ErrRateLimited,
		"ErrOriginForbidden": ErrOriginForbidden,
	}
}
