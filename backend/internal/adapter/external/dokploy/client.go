package dokploy

import (
	"context"
	"net"
	"net/http"
	"time"

	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

const maximumBodyBytes = 2 << 20

type Resolver interface {
	LookupIP(context.Context, string, string) ([]net.IP, error)
}

type ClientConfig struct {
	HTTPClient  *http.Client
	Resolver    Resolver
	Credentials portsout.CredentialStore
}

type Client struct {
	httpClient  *http.Client
	resolver    Resolver
	credentials portsout.CredentialStore
}

func NewClient(config ClientConfig) (*Client, error) {
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	resolver := config.Resolver
	if resolver == nil {
		resolver = net.DefaultResolver
	}
	return &Client{httpClient: httpClient, resolver: resolver, credentials: config.Credentials}, nil
}

var _ portsout.DokployGateway = (*Client)(nil)
