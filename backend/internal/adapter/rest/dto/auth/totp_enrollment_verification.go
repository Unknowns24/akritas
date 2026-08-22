package auth

type TotpEnrollmentVerificationRequest struct {
	EnrollmentID string `json:"enrollment_id"`
	TotpCode     string `json:"totp_code"`
}
