package github

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

const (
	APIVersion       = "2026-03-10"
	defaultAPIBase   = "https://api.github.com"
	maximumBodyBytes = 2 << 20
)

type ClientConfig struct {
	APIBaseURL  *url.URL
	HTTPClient  *http.Client
	Credentials portsout.CredentialStore
	Bindings    portsout.GitHubAppBindingReader
	Now         func() time.Time
}

type Client struct {
	base        *url.URL
	httpClient  *http.Client
	credentials portsout.CredentialStore
	bindings    portsout.GitHubAppBindingReader
	now         func() time.Time
	cacheMu     sync.Mutex
	tokenCache  map[string]cachedInstallationToken
}

type cachedInstallationToken struct {
	token     string
	expiresAt time.Time
}

func NewClient(config ClientConfig) (*Client, error) {
	base := config.APIBaseURL
	if base == nil {
		parsed, _ := url.Parse(defaultAPIBase)
		base = parsed
	}
	if base == nil || base.Host == "" || (base.Scheme != "https" && !(base.Scheme == "http" && isLoopbackHost(base.Hostname()))) {
		return nil, domain.ErrIntegrationUnavailable
	}
	baseCopy := *base
	baseCopy.Path = strings.TrimSuffix(baseCopy.Path, "/")
	httpClient := config.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 15 * time.Second}
	}
	clientCopy := *httpClient
	clientCopy.CheckRedirect = func(request *http.Request, _ []*http.Request) error {
		if !strings.EqualFold(request.URL.Scheme, baseCopy.Scheme) || !strings.EqualFold(request.URL.Host, baseCopy.Host) {
			return errors.New("github redirect rejected")
		}
		return nil
	}
	now := config.Now
	if now == nil {
		now = time.Now
	}
	return &Client{base: &baseCopy, httpClient: &clientCopy, credentials: config.Credentials, bindings: config.Bindings, now: now, tokenCache: make(map[string]cachedInstallationToken)}, nil
}

func isLoopbackHost(host string) bool {
	return host == "127.0.0.1" || host == "::1" || strings.EqualFold(host, "localhost")
}

var _ portsout.GitHubGateway = (*Client)(nil)
var _ portsout.GitHubAppGateway = (*Client)(nil)
