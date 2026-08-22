package github

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"strconv"
	"time"

	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func (c *Client) appJWT(ctx context.Context, ownerType string, ownerID uuid.UUID, appID int64) (string, error) {
	privateKeyPEM, err := c.credentials.Get(ctx, ownerType, ownerID, portsout.SecretKindGitHubPrivateKey)
	if err != nil {
		return "", err
	}
	defer wipe(privateKeyPEM)
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return "", ErrInvalidGitHubAppPrivateKey
	}
	var key *rsa.PrivateKey
	parsedPKCS1, pkcs1Err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if pkcs1Err == nil {
		key = parsedPKCS1
	} else {
		parsed, parseErr := x509.ParsePKCS8PrivateKey(block.Bytes)
		if parseErr != nil {
			return "", ErrInvalidGitHubAppPrivateKey
		}
		var ok bool
		key, ok = parsed.(*rsa.PrivateKey)
		if !ok {
			return "", ErrInvalidGitHubAppPrivateKey
		}
	}
	header, _ := json.Marshal(map[string]string{"alg": "RS256", "typ": "JWT"})
	now := c.now().UTC()
	claims, _ := json.Marshal(map[string]any{"iat": now.Add(-time.Minute).Unix(), "exp": now.Add(9 * time.Minute).Unix(), "iss": strconv.FormatInt(appID, 10)})
	unsigned := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(claims)
	digest := sha256.Sum256([]byte(unsigned))
	signature, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, digest[:])
	if err != nil {
		return "", ErrInvalidGitHubAppPrivateKey.Wrap(err)
	}
	return unsigned + "." + base64.RawURLEncoding.EncodeToString(signature), nil
}
