package mapper

import (
	qvacdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/qvac"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
)

func QvacConfigurationToDTO(value domain.QvacConfiguration) qvacdto.ConfigurationDTO {
	return qvacdto.ConfigurationDTO{
		EndpointURL: value.EndpointURL, ConnectionTimeoutSeconds: value.ConnectionTimeoutSeconds,
		ContextSize: value.ContextSize, AuthenticationType: string(value.AuthenticationType), CredentialConfigured: value.CredentialConfigured,
		UpdatedAt: value.UpdatedAt,
	}
}

func QvacStatusToDTO(value portsin.QvacRuntimeStatus) qvacdto.RuntimeStatusDTO {
	return qvacdto.RuntimeStatusDTO{
		ConnectionStatus: string(value.ConnectionStatus), Runtime: value.Runtime, ActiveModel: value.ActiveModel,
		Version: value.Version, LatencyMS: value.Latency, CheckedAt: value.CheckedAt,
	}
}

func PutQvacConfigurationToCommand(value qvacdto.PutConfigurationRequestDTO) portsin.PutQvacConfigurationCommand {
	return portsin.PutQvacConfigurationCommand{
		EndpointURL: value.EndpointURL, ConnectionTimeoutSeconds: value.ConnectionTimeoutSeconds,
		ContextSize: value.ContextSize,
		Authentication: portsin.QvacAuthenticationCommand{
			Type: domain.QvacAuthenticationType(value.Authentication.Type), BearerToken: value.Authentication.BearerToken,
			BasicUsername: value.Authentication.BasicUsername, BasicPassword: value.Authentication.BasicPassword,
		},
	}
}
