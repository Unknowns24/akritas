package auth

type LoginRequestDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TotpCode string `json:"totp_code"`
}
