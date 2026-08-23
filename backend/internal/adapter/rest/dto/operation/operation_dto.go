package operation

import "time"

type OperationDTO struct {
	ID           string     `json:"id"`
	Type         string     `json:"type"`
	Status       string     `json:"status"`
	ResourceType string     `json:"resource_type,omitempty"`
	ResourceID   string     `json:"resource_id,omitempty"`
	UserMessage  string     `json:"user_message,omitempty"`
	FailureCode  string     `json:"failure_code,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
	FinishedAt   *time.Time `json:"finished_at,omitempty"`
}
