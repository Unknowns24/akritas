package mapper

import (
	operationdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/operation"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func OperationToDTO(value domain.Operation) operationdto.OperationDTO {
	dto := operationdto.OperationDTO{
		ID: value.ID.String(), Type: string(value.Type), Status: string(value.Status),
		UserMessage: value.UserMessage, CreatedAt: value.CreatedAt, UpdatedAt: value.UpdatedAt, FinishedAt: value.FinishedAt,
	}
	if value.ResourceType != nil {
		dto.ResourceType = string(*value.ResourceType)
	}
	if value.ResourceID != nil {
		dto.ResourceID = value.ResourceID.String()
	}
	if value.FailureCode != nil {
		dto.FailureCode = *value.FailureCode
	}
	return dto
}
