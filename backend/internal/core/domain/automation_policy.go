package domain

type AutomationPolicy struct {
	AutomaticInvestigation bool
	AutomaticRemediation   bool
	AutomaticPullRequest   bool
}

func DefaultAutomationPolicy() AutomationPolicy {
	return AutomationPolicy{AutomaticInvestigation: true, AutomaticRemediation: true, AutomaticPullRequest: true}
}

func NewAutomationPolicy(investigation, remediation, pullRequest bool) (AutomationPolicy, error) {
	policy := AutomationPolicy{
		AutomaticInvestigation: investigation,
		AutomaticRemediation:   remediation,
		AutomaticPullRequest:   pullRequest,
	}
	if err := policy.Validate(); err != nil {
		return AutomationPolicy{}, err
	}
	return policy, nil
}

func (p AutomationPolicy) Validate() error {
	if (p.AutomaticRemediation && !p.AutomaticInvestigation) ||
		(p.AutomaticPullRequest && (!p.AutomaticInvestigation || !p.AutomaticRemediation)) {
		return ErrInvalidAutomationPolicy.Wrap(validationCause("automation dependencies"))
	}
	return nil
}
