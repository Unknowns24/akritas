package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type DokployServer struct {
	ID                   uuid.UUID         `gorm:"column:id;type:uuid;primaryKey"`
	Name                 string            `gorm:"column:name"`
	BaseURL              string            `gorm:"column:base_url"`
	ServerIdentifier     string            `gorm:"column:server_identifier"`
	ConnectionStatus     IntegrationStatus `gorm:"column:connection_status"`
	CredentialConfigured bool              `gorm:"column:credential_configured"`
	ApplicationCount     int               `gorm:"column:application_count"`
	ComposeCount         int               `gorm:"column:compose_count"`
	LastSyncedAt         *time.Time        `gorm:"column:last_synced_at"`
	CreatedAt            time.Time         `gorm:"column:created_at"`
	UpdatedAt            time.Time         `gorm:"column:updated_at"`
}

func NewDokployServer(
	id uuid.UUID,
	name, baseURL, serverIdentifier string,
	connectionStatus IntegrationStatus,
	createdAt time.Time,
) (*DokployServer, error) {
	server := &DokployServer{
		ID: id, Name: strings.TrimSpace(name), BaseURL: strings.TrimSpace(baseURL),
		ServerIdentifier: strings.TrimSpace(serverIdentifier), ConnectionStatus: connectionStatus,
		CreatedAt: createdAt, UpdatedAt: createdAt,
	}
	if err := server.Validate(); err != nil {
		return nil, err
	}
	return server, nil
}

func (s DokployServer) Validate() error {
	if s.ID == uuid.Nil || !nonBlank(s.Name) || !validHTTPURL(s.BaseURL) || !nonBlank(s.ServerIdentifier) ||
		s.ConnectionStatus.Validate() != nil || s.ApplicationCount < 0 || s.ComposeCount < 0 || !validTime(s.CreatedAt) || s.UpdatedAt.Before(s.CreatedAt) {
		return ErrInvalidDokployServer.Wrap(validationCause("Dokploy server"))
	}
	return nil
}
