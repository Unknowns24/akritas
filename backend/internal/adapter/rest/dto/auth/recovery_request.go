package auth

type RecoveryRequestDTO struct {
	Email          string `json:"email"`
	NewPassword    string `json:"new_password"`
	BootstrapToken string `json:"bootstrap_token"`
}
