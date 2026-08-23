package githubapp

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"net"
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
	if err != nil || !validPublicURL(parsed) || store == nil || gateway == nil || newID == nil || now == nil {
		return nil, ErrInvalidConfiguration
	}
	return &UseCase{store: store, gateway: gateway, publicURL: parsed, newID: newID, now: now, newState: randomState}, nil
}

func validPublicURL(parsed *url.URL) bool {
	if parsed == nil || parsed.Host == "" || parsed.User != nil || parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" {
		return false
	}
	if parsed.Scheme == "https" {
		return true
	}
	return parsed.Scheme == "http" && isLoopbackHost(parsed.Hostname())
}

func isLoopbackHost(host string) bool {
	normalized := strings.ToLower(strings.TrimSpace(host))
	if normalized == "localhost" {
		return true
	}
	ip := net.ParseIP(normalized)
	return ip != nil && ip.IsLoopback()
}

func randomState() (string, error) {
	value := make([]byte, 32)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func wipe(value []byte) { clear(value) }
