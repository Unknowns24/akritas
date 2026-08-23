package github

type GitHubManifestRegistrationDTO struct {
	RegistrationID string `json:"registration_id"`
	FormAction     string `json:"form_action"`
	Manifest       string `json:"manifest"`
	State          string `json:"state"`
	ExpiresAt      string `json:"expires_at"`
}
