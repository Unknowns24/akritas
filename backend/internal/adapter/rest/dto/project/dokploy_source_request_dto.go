package project

import "github.com/google/uuid"

type DokploySourceRequestDTO struct {
	Type               string    `json:"type"`
	DokployServerID    uuid.UUID `json:"dokploy_server_id"`
	ResourceIdentifier string    `json:"resource_identifier"`
	ServiceName        string    `json:"service_name,omitempty"`
}
