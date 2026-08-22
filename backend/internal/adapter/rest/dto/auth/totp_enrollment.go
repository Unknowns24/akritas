package auth

type TotpEnrollment struct {
	EnrollmentID   string `json:"enrollment_id"`
	OtpauthURI     string `json:"otpauth_uri"`
	ManualEntryKey string `json:"manual_entry_key"`
	ExpiresAt      string `json:"expires_at"`
}

type TotpEnrollmentResponse struct {
	Data TotpEnrollment `json:"data"`
}
