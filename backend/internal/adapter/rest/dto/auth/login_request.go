package auth

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TotpCode string `json:"totp_code"`
}
