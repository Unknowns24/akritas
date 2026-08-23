package common

type ErrorDTO struct {
	Code        string           `json:"code"`
	Message     string           `json:"message"`
	UserMessage string           `json:"user_message"`
	RequestID   string           `json:"request_id"`
	Details     []ErrorDetailDTO `json:"details,omitempty"`
}
