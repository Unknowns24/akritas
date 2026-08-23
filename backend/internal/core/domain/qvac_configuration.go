package domain

import (
	"net/url"
	"strings"
	"time"
)

type QvacAuthenticationType string

const (
	QvacAuthenticationNone   QvacAuthenticationType = "none"
	QvacAuthenticationBearer QvacAuthenticationType = "bearer"
	QvacAuthenticationBasic  QvacAuthenticationType = "basic"
	DefaultQvacContextSize   int                    = 32768
)

func (t QvacAuthenticationType) Validate() error {
	switch t {
	case QvacAuthenticationNone, QvacAuthenticationBearer, QvacAuthenticationBasic:
		return nil
	default:
		return ErrInvalidIntegrationStatus.Wrap(validationCause("QVAC authentication type"))
	}
}

type QvacConfiguration struct {
	EndpointURL              string
	ConnectionTimeoutSeconds int
	ContextSize              int
	AuthenticationType       QvacAuthenticationType
	CredentialConfigured     bool
	BasicUsername            string
	UpdatedAt                time.Time
}

func DefaultQvacConfiguration(now time.Time) QvacConfiguration {
	return QvacConfiguration{
		EndpointURL:              "http://127.0.0.1:11434/v1",
		ConnectionTimeoutSeconds: 180,
		ContextSize:              DefaultQvacContextSize,
		AuthenticationType:       QvacAuthenticationNone,
		UpdatedAt:                now,
	}
}

func NewQvacConfiguration(endpoint string, timeoutSeconds int, authType QvacAuthenticationType, credentialConfigured bool, basicUsername string, updatedAt time.Time) (QvacConfiguration, error) {
	return NewQvacConfigurationWithContext(endpoint, timeoutSeconds, DefaultQvacContextSize, authType, credentialConfigured, basicUsername, updatedAt)
}

func NewQvacConfigurationWithContext(endpoint string, timeoutSeconds int, contextSize int, authType QvacAuthenticationType, credentialConfigured bool, basicUsername string, updatedAt time.Time) (QvacConfiguration, error) {
	value := QvacConfiguration{
		EndpointURL:              strings.TrimRight(strings.TrimSpace(endpoint), "/"),
		ConnectionTimeoutSeconds: timeoutSeconds,
		ContextSize:              contextSize,
		AuthenticationType:       authType,
		CredentialConfigured:     credentialConfigured,
		BasicUsername:            strings.TrimSpace(basicUsername),
		UpdatedAt:                updatedAt,
	}
	if err := value.Validate(); err != nil {
		return QvacConfiguration{}, err
	}
	return value, nil
}

func (c QvacConfiguration) Validate() error {
	parsed, err := url.Parse(strings.TrimSpace(c.EndpointURL))
	if err != nil || parsed.Host == "" || parsed.User != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return ErrInvalidIntegrationStatus.Wrap(validationCause("QVAC endpoint"))
	}
	if c.ConnectionTimeoutSeconds < 1 || c.ConnectionTimeoutSeconds > 300 || c.AuthenticationType.Validate() != nil || !validTime(c.UpdatedAt) {
		return ErrInvalidIntegrationStatus.Wrap(validationCause("QVAC configuration"))
	}
	if c.ContextSize < 4096 || c.ContextSize > 131072 {
		return ErrInvalidIntegrationStatus.Wrap(validationCause("QVAC context size"))
	}
	if c.AuthenticationType == QvacAuthenticationBasic && c.BasicUsername == "" {
		return ErrInvalidIntegrationStatus.Wrap(validationCause("QVAC basic username"))
	}
	if c.AuthenticationType == QvacAuthenticationNone && c.CredentialConfigured {
		return ErrInvalidIntegrationStatus.Wrap(validationCause("QVAC credential"))
	}
	return nil
}
