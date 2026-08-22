package dto

type ConnectionTestDTO struct {
	Status      string `json:"status"`
	CheckedAt   string `json:"checked_at"`
	LatencyMS   *int64 `json:"latency_ms,omitempty"`
	UserMessage string `json:"user_message"`
}
