package dto

type ErrorDetail struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

type Error struct {
	Code        string        `json:"code"`
	Message     string        `json:"message"`
	UserMessage string        `json:"user_message"`
	RequestID   string        `json:"request_id"`
	Details     []ErrorDetail `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error Error `json:"error"`
}
