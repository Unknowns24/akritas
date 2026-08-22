package projectdto

import (
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/envelope"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
)

type ProjectListDTO struct {
	Data   []ProjectSummaryDTO `json:"data"`
	Paging envelope.PagingDTO  `json:"paging"`
}

func FromList(projects []domain.Project, page paging.Page) ProjectListDTO {
	data := make([]ProjectSummaryDTO, 0, len(projects))
	for _, project := range projects {
		data = append(data, FromProjectSummary(project))
	}
	return ProjectListDTO{
		Data: data,
		Paging: envelope.PagingDTO{
			Limit: page.Limit, Total: page.Total, HasMore: page.HasMore,
			NextCursor: page.NextCursor, PrevCursor: page.PrevCursor,
		},
	}
}
