package qvac

import (
	"context"
	"errors"
	"log"
	"time"

	qvacexternal "github.com/Unknowns24/akritas/backend/internal/adapter/external/qvac"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

var ErrInvalidUseCase = errors.New("invalid QVAC use case configuration")

var qvacCredentialOwnerID = uuid.MustParse("00000000-0000-0000-0000-000000000001")

type UseCase struct {
	configs     portsout.QvacConfigurationStore
	credentials portsout.CredentialStore
	now         func() time.Time
}

func New(configs portsout.QvacConfigurationStore, credentials portsout.CredentialStore, now func() time.Time) (*UseCase, error) {
	if configs == nil || credentials == nil || now == nil {
		return nil, ErrInvalidUseCase
	}
	return &UseCase{configs: configs, credentials: credentials, now: now}, nil
}

func (uc *UseCase) GetConfiguration(ctx context.Context) (domain.QvacConfiguration, error) {
	return uc.configs.Get(ctx)
}

func (uc *UseCase) PutConfiguration(ctx context.Context, cmd portsin.PutQvacConfigurationCommand) (domain.QvacConfiguration, error) {
	authType := cmd.Authentication.Type
	credentialConfigured := false
	basicUsername := ""
	if authType == domain.QvacAuthenticationBearer {
		if cmd.Authentication.BearerToken == "" {
			return domain.QvacConfiguration{}, domain.ErrInvalidIntegrationStatus
		}
		credentialConfigured = true
	} else if authType == domain.QvacAuthenticationBasic {
		if cmd.Authentication.BasicUsername == "" || cmd.Authentication.BasicPassword == "" {
			return domain.QvacConfiguration{}, domain.ErrInvalidIntegrationStatus
		}
		basicUsername = cmd.Authentication.BasicUsername
		credentialConfigured = true
	} else if authType == domain.QvacAuthenticationNone {
		if cmd.Authentication.BearerToken != "" || cmd.Authentication.BasicUsername != "" || cmd.Authentication.BasicPassword != "" {
			return domain.QvacConfiguration{}, domain.ErrInvalidIntegrationStatus
		}
	}
	contextSize := cmd.ContextSize
	if contextSize == 0 {
		contextSize = domain.DefaultQvacContextSize
	}
	config, err := domain.NewQvacConfigurationWithContext(cmd.EndpointURL, cmd.ConnectionTimeoutSeconds, contextSize, authType, credentialConfigured, basicUsername, uc.now().UTC())
	if err != nil {
		return domain.QvacConfiguration{}, err
	}
	if _, err := qvacexternal.NewClient(qvacexternal.ClientConfig{EndpointURL: config.EndpointURL, Timeout: time.Duration(config.ConnectionTimeoutSeconds) * time.Second}); err != nil {
		return domain.QvacConfiguration{}, err
	}
	if err := uc.credentials.DeleteOwner(ctx, portsout.CredentialOwnerQvacConfiguration, qvacCredentialOwnerID); err != nil {
		return domain.QvacConfiguration{}, err
	}
	switch authType {
	case domain.QvacAuthenticationBearer:
		if err := uc.credentials.Put(ctx, portsout.CredentialOwnerQvacConfiguration, qvacCredentialOwnerID, portsout.SecretValue{Kind: portsout.SecretKindQvacBearerToken, Plaintext: []byte(cmd.Authentication.BearerToken)}); err != nil {
			return domain.QvacConfiguration{}, err
		}
	case domain.QvacAuthenticationBasic:
		if err := uc.credentials.Put(ctx, portsout.CredentialOwnerQvacConfiguration, qvacCredentialOwnerID, portsout.SecretValue{Kind: portsout.SecretKindQvacBasicPassword, Plaintext: []byte(cmd.Authentication.BasicPassword)}); err != nil {
			return domain.QvacConfiguration{}, err
		}
	}
	if err := uc.configs.Put(ctx, config); err != nil {
		return domain.QvacConfiguration{}, err
	}
	return config, nil
}

func (uc *UseCase) Client(ctx context.Context) (*qvacexternal.Client, error) {
	config, err := uc.configs.Get(ctx)
	if err != nil {
		return nil, err
	}
	clientConfig := qvacexternal.ClientConfig{
		EndpointURL: config.EndpointURL,
		Timeout:     time.Duration(config.ConnectionTimeoutSeconds) * time.Second,
		ContextSize: config.ContextSize,
	}
	switch config.AuthenticationType {
	case domain.QvacAuthenticationBearer:
		secret, err := uc.credentials.Get(ctx, portsout.CredentialOwnerQvacConfiguration, qvacCredentialOwnerID, portsout.SecretKindQvacBearerToken)
		if err != nil {
			return nil, err
		}
		clientConfig.APIKey = string(secret)
	case domain.QvacAuthenticationBasic:
		secret, err := uc.credentials.Get(ctx, portsout.CredentialOwnerQvacConfiguration, qvacCredentialOwnerID, portsout.SecretKindQvacBasicPassword)
		if err != nil {
			return nil, err
		}
		clientConfig.BasicUsername = config.BasicUsername
		clientConfig.BasicPassword = string(secret)
	}
	return qvacexternal.NewClient(clientConfig)
}

func (uc *UseCase) TestConnection(ctx context.Context) (portsin.ConnectionTestResult, error) {
	started := uc.now().UTC()
	client, err := uc.Client(ctx)
	if err != nil {
		log.Printf("qvac: connection test configuration unavailable error=%v", err)
		return portsin.ConnectionTestResult{Status: domain.ConnectionTestUnavailable, CheckedAt: started.Format(time.RFC3339), UserMessage: "La configuración QVAC no está disponible."}, nil
	}
	latencyDuration, err := client.Ping(ctx)
	if err != nil {
		log.Printf("qvac: connection test unavailable endpoint=%s model=%s context_size=%d error=%v", client.Endpoint(), client.Model(), client.ContextSize(), err)
		return portsin.ConnectionTestResult{Status: domain.ConnectionTestUnavailable, CheckedAt: started.Format(time.RFC3339), UserMessage: "QVAC no responde desde el backend en el endpoint configurado."}, nil
	}
	latency := latencyDuration.Milliseconds()
	log.Printf("qvac: connection test connected endpoint=%s model=%s context_size=%d latency_ms=%d", client.Endpoint(), client.Model(), client.ContextSize(), latency)
	return portsin.ConnectionTestResult{Status: domain.ConnectionTestConnected, CheckedAt: started.Format(time.RFC3339), LatencyMS: &latency, UserMessage: "QVAC respondió correctamente."}, nil
}

func (uc *UseCase) GetStatus(ctx context.Context) (portsin.QvacRuntimeStatus, error) {
	config, err := uc.configs.Get(ctx)
	if err != nil {
		return portsin.QvacRuntimeStatus{}, err
	}
	result, err := uc.TestConnection(ctx)
	if err != nil {
		return portsin.QvacRuntimeStatus{}, err
	}
	status := domain.IntegrationStatusUnavailable
	if result.Status == domain.ConnectionTestConnected {
		status = domain.IntegrationStatusConnected
	}
	latency := int64(0)
	if result.LatencyMS != nil {
		latency = *result.LatencyMS
	}
	return portsin.QvacRuntimeStatus{
		ConnectionStatus: status,
		Runtime:          "QVAC Local",
		ActiveModel:      "akritas",
		Version:          "",
		Latency:          latency,
		CheckedAt:        result.CheckedAt,
	}, config.Validate()
}
