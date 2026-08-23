package githubapp

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
	"time"

	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

var ErrInvalidConfiguration = errors.New("invalid GitHub App use case configuration")

type UseCase struct {
	store     portsout.GitHubAppRegistrationStore
	gateway   portsout.GitHubAppGateway
	publicURL *url.URL
	newID     func() uuid.UUID
	now       func() time.Time
	newState  func() (string, error)
}

func New(store portsout.GitHubAppRegistrationStore, gateway portsout.GitHubAppGateway, publicURL string, newID func() uuid.UUID, now func() time.Time) (*UseCase, error) {
	parsed, err := url.Parse(strings.TrimRight(publicURL, "/"))
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" || store == nil || gateway == nil || newID == nil || now == nil {
		return nil, ErrInvalidConfiguration
	}
	return &UseCase{store: store, gateway: gateway, publicURL: parsed, newID: newID, now: now, newState: randomState}, nil
}

func randomState() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func wipe(value []byte) { clear(value) }
