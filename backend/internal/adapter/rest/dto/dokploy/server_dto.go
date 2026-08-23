package dokploy

type DokployServerDTO struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	BaseURL              string  `json:"base_url"`
	ServerIdentifier     string  `json:"server_identifier"`
	ConnectionStatus     string  `json:"connection_status"`
	CredentialConfigured bool    `json:"credential_configured"`
	ApplicationCount     int     `json:"application_count"`
	ComposeCount         int     `json:"compose_count"`
	LastSyncedAt         *string `json:"last_synced_at,omitempty"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}
