package auth

type TOTPEnrollmentDTO struct {
	EnrollmentID   string `json:"enrollment_id"`
	OtpauthURI     string `json:"otpauth_uri"`
	ManualEntryKey string `json:"manual_entry_key"`
	ExpiresAt      string `json:"expires_at"`
}
