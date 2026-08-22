package auth

type SetupRequest struct {
	Email          string `json:"email"`
	DisplayName    string `json:"display_name"`
	Password       string `json:"password"`
	BootstrapToken string `json:"bootstrap_token"`
}
