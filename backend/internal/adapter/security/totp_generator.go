package security

import (
	"github.com/pquerna/otp/totp"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

// RFC 6238 defaults (6 digits, SHA1, 30s period, 20-byte secret) match ADR-008.
type totpSecretGenerator struct{}

func NewTOTPSecretGenerator() out.TOTPSecretGenerator {
	return &totpSecretGenerator{}
}

func (g *totpSecretGenerator) Generate(issuer, accountEmail string) (out.TOTPSecret, error) {
	key, err := totp.Generate(totp.GenerateOpts{Issuer: issuer, AccountName: accountEmail})
	if err != nil {
		return out.TOTPSecret{}, err
	}
	return out.TOTPSecret{Base32Key: key.Secret(), OtpauthURI: key.String()}, nil
}
