package dokploy

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (c *Client) do(ctx context.Context, requestURL, credential string) ([]byte, error) {
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
	if err != nil {
		return nil, domain.ErrIntegrationUnavailable
	}
	request.Header.Set("Accept", "application/json")
	request.Header.Set("x-api-key", credential)
	clientCopy := *c.httpClient
	transport, err := c.transportForRequest(request.URL.Scheme)
	if err != nil {
		return nil, &providerError{Cause: err}
	}
	clientCopy.Transport = transport
	origin := request.URL.Scheme + "://" + request.URL.Host
	clientCopy.CheckRedirect = func(redirect *http.Request, _ []*http.Request) error {
		if redirect.URL.Scheme+"://"+redirect.URL.Host != origin {
			return errors.New("dokploy redirect rejected")
		}
		return nil
	}
	response, err := clientCopy.Do(request)
	if err != nil {
		return nil, &providerError{Cause: err}
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, maximumBodyBytes))
		return nil, &providerError{Status: response.StatusCode}
	}
	var buffer bytes.Buffer
	read, err := io.Copy(&buffer, io.LimitReader(response.Body, maximumBodyBytes+1))
	if err != nil || read > maximumBodyBytes {
		return nil, &providerError{Status: response.StatusCode}
	}
	return buffer.Bytes(), nil
}

func (c *Client) transportForRequest(scheme string) (*http.Transport, error) {
	base := c.httpClient.Transport
	if base == nil {
		base = http.DefaultTransport
	}
	transport, ok := base.(*http.Transport)
	if !ok {
		return nil, errors.New("unsupported Dokploy HTTP transport")
	}
	clone := transport.Clone()
	clone.DialContext = func(dialContext context.Context, network, address string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(address)
		if err != nil {
			return nil, errors.New("invalid Dokploy dial address")
		}
		ips, err := c.resolve(dialContext, host)
		if err != nil || len(ips) == 0 {
			return nil, errors.New("Dokploy host resolution failed")
		}
		dialer := net.Dialer{Timeout: 10 * time.Second}
		for _, ip := range ips {
			if isForbiddenIP(ip) || (scheme == "http" && !ip.IsLoopback() && !ip.IsPrivate()) {
				continue
			}
			connection, dialErr := dialer.DialContext(dialContext, network, net.JoinHostPort(ip.String(), port))
			if dialErr == nil {
				return connection, nil
			}
		}
		return nil, errors.New("Dokploy destination rejected")
	}
	return clone, nil
}

type providerError struct {
	Status int
	Cause  error
}

func (e *providerError) Error() string { return "Dokploy provider request failed" }
func (e *providerError) Unwrap() error { return e.Cause }

func normalizeProviderError(err error) error {
	var providerErr *providerError
	if errors.As(err, &providerErr) && (providerErr.Status == http.StatusUnauthorized || providerErr.Status == http.StatusForbidden) {
		return domain.ErrDokployCredentialRejected
	}
	return domain.ErrIntegrationUnavailable.Wrap(err)
}

func wipe(value []byte) { clear(value) }
