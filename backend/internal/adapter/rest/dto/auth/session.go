package auth

type Session struct {
	Administrator     Administrator `json:"administrator"`
	AuthenticatedAt   string        `json:"authenticated_at"`
	IdleExpiresAt     string        `json:"idle_expires_at"`
	AbsoluteExpiresAt string        `json:"absolute_expires_at"`
}

type SessionResponse struct {
	Data Session `json:"data"`
}
