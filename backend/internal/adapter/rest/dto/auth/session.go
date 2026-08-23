package auth

type SessionDTO struct {
	Administrator     AdministratorDTO `json:"administrator"`
	AuthenticatedAt   string           `json:"authenticated_at"`
	IdleExpiresAt     string           `json:"idle_expires_at"`
	AbsoluteExpiresAt string           `json:"absolute_expires_at"`
}
