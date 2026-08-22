package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type DokployServer struct {
	ID                   uuid.UUID
	Name                 string
	BaseURL              string
	ServerIdentifier     string
	ConnectionStatus     IntegrationStatus
	CredentialConfigured bool
	ApplicationCount     int
	LastSyncedAt         *time.Time
	CreatedAt            time.Time
	UpdatedAt            time.Time
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
		s.ConnectionStatus.Validate() != nil || s.ApplicationCount < 0 || !validTime(s.CreatedAt) || s.UpdatedAt.Before(s.CreatedAt) {
		return ErrInvalidDokployServer.Wrap(validationCause("Dokploy server"))
	}
	return nil
}
