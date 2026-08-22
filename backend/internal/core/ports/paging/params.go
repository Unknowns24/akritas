package paging

import "github.com/Unknowns24/akritas/backend/internal/core/domain"

type ListQuery struct {
	Limit              int
	Offset             int
	Sort               string
	NameLike           string
	MonitoringStatusIn []domain.MonitoringStatus
}

type Page struct {
	Limit      int
	Total      int64
	HasMore    bool
	NextCursor string
	PrevCursor string
}

func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return 20
	}
	if limit > 100 {
		return 100
	}
	return limit
}
