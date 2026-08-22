package dokploy

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net"
	"net/url"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func (c *Client) normalizeAndValidateURL(ctx context.Context, rawURL string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", domain.ErrInvalidDokployServer
	}
	if parsed.Path != "" && parsed.Path != "/" {
		return "", domain.ErrInvalidDokployServer
	}
	ips, err := c.resolve(ctx, parsed.Hostname())
	if err != nil || len(ips) == 0 {
		return "", domain.ErrIntegrationUnavailable
	}
	for _, ip := range ips {
		if isForbiddenIP(ip) || (parsed.Scheme == "http" && !ip.IsLoopback() && !ip.IsPrivate()) {
			return "", domain.ErrInvalidDokployServer
		}
	}
	parsed.Path = ""
	return strings.TrimSuffix(parsed.String(), "/"), nil
}

func (c *Client) resolve(ctx context.Context, host string) ([]net.IP, error) {
	if ip := net.ParseIP(host); ip != nil {
		return []net.IP{ip}, nil
	}
	return c.resolver.LookupIP(ctx, "ip", host)
}

func isForbiddenIP(ip net.IP) bool {
	return ip == nil || ip.IsUnspecified() || ip.IsMulticast() || ip.IsLinkLocalMulticast() || ip.IsLinkLocalUnicast()
}

func fingerprint(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}
