package errors

import "github.com/Unknowns24/akritas/backend/internal/core/domain"

var ErrIntegrationPersistence = &domain.Error{
	Code: "0x202001I", Message: "integration persistence failure", UserMessage: "No se pudo guardar la integración.",
}

var ErrProjectPersistence = &domain.Error{
	Code: "0x203001I", Message: "project persistence failure", UserMessage: "No se pudo guardar el proyecto.",
}

var ErrInvestigationPersistence = &domain.Error{
	Code: "0x204001I", Message: "investigation persistence failure", UserMessage: "No se pudo guardar la investigación.",
}

var ErrOperationPersistence = &domain.Error{
	Code: "0x205001I", Message: "operation persistence failure", UserMessage: "No se pudo guardar la operación.",
}

var ErrEvidencePersistence = &domain.Error{
	Code: "0x206001I", Message: "evidence persistence failure", UserMessage: "No se pudo guardar la evidencia.",
}

func Catalog() map[string]*domain.Error {
	return map[string]*domain.Error{
		"ErrIntegrationPersistence":   ErrIntegrationPersistence,
		"ErrProjectPersistence":       ErrProjectPersistence,
		"ErrInvestigationPersistence": ErrInvestigationPersistence,
		"ErrOperationPersistence":     ErrOperationPersistence,
		"ErrEvidencePersistence":      ErrEvidencePersistence,
	}
var ErrIncidentPersistence = &domain.Error{
	Code: "0x207001I", Message: "incident persistence failure", UserMessage: "No se pudo consultar el incidente.",
}

var ErrMonitoringPersistence = &domain.Error{
	Code: "0x208001I", Message: "monitoring persistence failure", UserMessage: "No se pudo guardar el estado de monitoreo.",
}

func Catalog() map[string]*domain.Error {
	return map[string]*domain.Error{"ErrIntegrationPersistence": ErrIntegrationPersistence, "ErrProjectPersistence": ErrProjectPersistence, "ErrIncidentPersistence": ErrIncidentPersistence, "ErrMonitoringPersistence": ErrMonitoringPersistence}
}
