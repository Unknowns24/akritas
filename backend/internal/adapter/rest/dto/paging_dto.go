package dto

type PagingDTO struct {
	Limit      int    `json:"limit"`
	Total      int64  `json:"total"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
	PrevCursor string `json:"prev_cursor"`
}
