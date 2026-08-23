package mapper

import (
	automationdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/automation"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func AutomationPolicyToDTO(value domain.AutomationPolicy) automationdto.PolicyDTO {
	return automationdto.PolicyDTO{
		AutomaticInvestigation: value.AutomaticInvestigation,
		AutomaticRemediation:   value.AutomaticRemediation,
		AutomaticPullRequest:   value.AutomaticPullRequest,
	}
}

func AutomationPolicyToDomain(value automationdto.PolicyDTO) (domain.AutomationPolicy, error) {
	return domain.NewAutomationPolicy(value.AutomaticInvestigation, value.AutomaticRemediation, value.AutomaticPullRequest)
}
