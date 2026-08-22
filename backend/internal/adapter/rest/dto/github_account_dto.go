package dto

type GitHubAccountDTO struct {
	ID                   string  `json:"id"`
	DisplayName          string  `json:"display_name"`
	AccountType          string  `json:"account_type"`
	AccountIdentifier    string  `json:"account_identifier"`
	AuthenticationMethod string  `json:"authentication_method"`
	AuthenticationStatus string  `json:"authentication_status"`
	CredentialConfigured bool    `json:"credential_configured"`
	RepositoryCount      int     `json:"repository_count"`
	LastCheckedAt        *string `json:"last_checked_at,omitempty"`
	ManageURL            string  `json:"manage_url,omitempty"`
	CreatedAt            string  `json:"created_at"`
	UpdatedAt            string  `json:"updated_at"`
}
