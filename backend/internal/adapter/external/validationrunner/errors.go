package validationrunner

import "github.com/Unknowns24/akritas/backend/internal/core/domain"

var ErrValidationToolUnavailable = &domain.Error{
	Code: "0x304001I", Message: "validation tool unavailable", UserMessage: "La herramienta de validación no está disponible.",
}

var ErrValidationExecutionFailed = &domain.Error{
	Code: "0x304002I", Message: "validation execution failed", UserMessage: "No se pudo ejecutar la validación.",
}

func Catalog() map[string]*domain.Error {
	return map[string]*domain.Error{
		"ErrValidationToolUnavailable": ErrValidationToolUnavailable,
		"ErrValidationExecutionFailed": ErrValidationExecutionFailed,
	}
}
