package mapper

import (
	incidentdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/incident"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func TimelineEventToDTO(value domain.TimelineEvent) incidentdto.TimelineEventDTO {
	return incidentdto.TimelineEventDTO{
		ID: value.ID.String(), IncidentID: value.IncidentID.String(), Type: string(value.Type),
		OccurredAt: value.OccurredAt, Summary: value.Summary, Detail: value.Detail,
	}
}
