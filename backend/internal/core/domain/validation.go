package domain

import (
	"fmt"
	"net/url"
	"strings"
	"time"
)

func cloneStrings(values []string) []string {
	if len(values) == 0 {
		return []string{}
	}
	return append([]string(nil), values...)
}

func nonBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

func validTime(value time.Time) bool {
	return !value.IsZero()
}

func validHTTPURL(value string) bool {
	parsed, err := url.ParseRequestURI(value)
	return err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https") && parsed.Host != ""
}

func validationCause(field string) error {
	return fmt.Errorf("invalid %s", field)
}
