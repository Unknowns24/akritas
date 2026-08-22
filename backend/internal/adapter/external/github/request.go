package github

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (c *Client) doJSON(ctx context.Context, method, path, token string, body io.Reader, target any) (*http.Response, error) {
	reference, err := url.Parse(path)
	if err != nil {
		return nil, domain.ErrIntegrationUnavailable
	}
	requestURL := c.base.ResolveReference(reference)
	request, err := http.NewRequestWithContext(ctx, method, requestURL.String(), body)
	if err != nil {
		return nil, domain.ErrIntegrationUnavailable
	}
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", APIVersion)
	if method == http.MethodPost || method == http.MethodPatch {
		request.Header.Set("Content-Type", "application/json")
	}
	if token != "" {
		request.Header.Set("Authorization", "Bearer "+token)
	}
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, &providerError{Cause: err}
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, maximumBodyBytes))
		return response, &providerError{Status: response.StatusCode}
	}
	if target != nil {
		var buffer bytes.Buffer
		read, readErr := io.Copy(&buffer, io.LimitReader(response.Body, maximumBodyBytes+1))
		if readErr != nil || read > maximumBodyBytes {
			return response, &providerError{Status: response.StatusCode}
		}
		if err := json.Unmarshal(buffer.Bytes(), target); err != nil && !errors.Is(err, io.EOF) {
			return response, &providerError{Status: response.StatusCode, Cause: err}
		}
	}
	return response, nil
}

type providerError struct {
	Status int
	Cause  error
}

func (e *providerError) Error() string { return "GitHub provider request failed" }
func (e *providerError) Unwrap() error { return e.Cause }

func normalizeCredentialError(err error) error {
	var providerErr *providerError
	if errors.As(err, &providerErr) && (providerErr.Status == http.StatusUnauthorized || providerErr.Status == http.StatusForbidden || providerErr.Status == http.StatusNotFound) {
		return domain.ErrGitHubCredentialRejected
	}
	return domain.ErrIntegrationUnavailable.Wrap(err)
}

func normalizeProviderError(err error) error {
	var providerErr *providerError
	if errors.As(err, &providerErr) && (providerErr.Status == http.StatusUnauthorized || providerErr.Status == http.StatusForbidden) {
		return domain.ErrGitHubCredentialRejected
	}
	return domain.ErrIntegrationUnavailable.Wrap(err)
}

func wipe(value []byte) { clear(value) }
