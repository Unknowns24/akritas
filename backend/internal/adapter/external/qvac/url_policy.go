package qvac

import (
	"net"
	"net/url"
	"strings"
)

func validateEndpoint(raw string) (*url.URL, error) {
	parsed, err := url.Parse(strings.TrimSpace(raw))
	if err != nil || parsed.Host == "" || parsed.User != nil {
		return nil, ErrInvalidEndpoint
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return nil, ErrInvalidEndpoint
	}
	host := parsed.Hostname()
	ip := net.ParseIP(host)
	switch {
	case ip != nil:
		if !ip.IsLoopback() && !ip.IsPrivate() {
			return nil, ErrInvalidEndpoint
		}
	case strings.EqualFold(host, "localhost"):
	default:
		return nil, ErrInvalidEndpoint
	}
	copyURL := *parsed
	copyURL.Path = strings.TrimRight(copyURL.Path, "/")
	if copyURL.Path == "" {
		copyURL.Path = "/v1"
	}
	return &copyURL, nil
}
