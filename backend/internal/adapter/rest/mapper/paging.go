package mapper

import (
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func PagingToDTO(value ukerpagination.PagingBlock) dto.PagingDTO {
	return dto.PagingDTO{Limit: value.Limit, Total: value.Total, HasMore: value.HasMore, NextCursor: value.NextCursor, PrevCursor: value.PrevCursor}
}
