package errors

import "github.com/Unknowns24/akritas/backend/internal/core/domain"

var ErrIntegrationPersistence = &domain.Error{
	Code: "0x202001I", Message: "integration persistence failure", UserMessage: "No se pudo guardar la integración.",
}

var ErrProjectPersistence = &domain.Error{
	Code: "0x203001I", Message: "project persistence failure", UserMessage: "No se pudo guardar el proyecto.",
}

func Catalog() map[string]*domain.Error {
	return map[string]*domain.Error{"ErrIntegrationPersistence": ErrIntegrationPersistence, "ErrProjectPersistence": ErrProjectPersistence}
}
