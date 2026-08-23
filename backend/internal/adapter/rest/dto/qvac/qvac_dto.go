package qvac

import "time"

type AuthenticationRequestDTO struct {
	Type          string `json:"type"`
	BearerToken   string `json:"bearer_token,omitempty"`
	BasicUsername string `json:"basic_username,omitempty"`
	BasicPassword string `json:"basic_password,omitempty"`
}

type PutConfigurationRequestDTO struct {
	EndpointURL              string                   `json:"endpoint_url"`
	ConnectionTimeoutSeconds int                      `json:"connection_timeout_seconds"`
	ContextSize              int                      `json:"context_size"`
	Authentication           AuthenticationRequestDTO `json:"authentication"`
}

type ConfigurationDTO struct {
	EndpointURL              string    `json:"endpoint_url"`
	ConnectionTimeoutSeconds int       `json:"connection_timeout_seconds"`
	ContextSize              int       `json:"context_size"`
	AuthenticationType       string    `json:"authentication_type"`
	CredentialConfigured     bool      `json:"credential_configured"`
	UpdatedAt                time.Time `json:"updated_at"`
}

type RuntimeStatusDTO struct {
	ConnectionStatus string `json:"connection_status"`
	Runtime          string `json:"runtime"`
	ActiveModel      string `json:"active_model,omitempty"`
	Version          string `json:"version,omitempty"`
	LatencyMS        int64  `json:"latency_ms,omitempty"`
	CheckedAt        string `json:"checked_at,omitempty"`
}
