package github

import (
	"path"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func sanitizeRepositoryPath(raw string) (string, error) {
	value := strings.TrimSpace(raw)
	value = strings.ReplaceAll(value, "\\", "/")
	if value == "" || strings.HasPrefix(value, "/") || strings.Contains(value, ":") {
		return "", domain.ErrIntegrationUnavailable.Wrap(errInvalidPath)
	}
	cleaned := path.Clean(value)
	if cleaned == "." || cleaned == ".." || strings.HasPrefix(cleaned, "../") || strings.Contains(cleaned, "/../") {
		return "", domain.ErrIntegrationUnavailable.Wrap(errInvalidPath)
	}
	for _, part := range strings.Split(cleaned, "/") {
		if part == ".." || part == "" {
			return "", domain.ErrIntegrationUnavailable.Wrap(errInvalidPath)
		}
	}
	return cleaned, nil
}

var errInvalidPath = errString("invalid repository path")

type errString string

func (e errString) Error() string { return string(e) }
