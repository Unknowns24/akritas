package out

// TOTPSecret is the plaintext provisioning material for a freshly generated
// secret. It must be shown to the caller at most once and never persisted
// in this form.
type TOTPSecret struct {
	Base32Key  string
	OtpauthURI string
}

type TOTPSecretGenerator interface {
	Generate(issuer, accountEmail string) (TOTPSecret, error)
}
