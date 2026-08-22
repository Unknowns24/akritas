package mapper

import (
	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func PagingToDTO(value ukerpagination.PagingBlock) commondto.PagingDTO {
	return commondto.PagingDTO{Limit: value.Limit, Total: value.Total, HasMore: value.HasMore, NextCursor: value.NextCursor, PrevCursor: value.PrevCursor}
}
