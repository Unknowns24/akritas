package auth

type TOTPEnrollmentVerificationRequestDTO struct {
	EnrollmentID string `json:"enrollment_id"`
	TotpCode     string `json:"totp_code"`
}
