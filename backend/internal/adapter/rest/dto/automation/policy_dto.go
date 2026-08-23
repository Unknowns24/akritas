package automation

type PolicyDTO struct {
	AutomaticInvestigation bool `json:"automatic_investigation"`
	AutomaticRemediation   bool `json:"automatic_remediation"`
	AutomaticPullRequest   bool `json:"automatic_pull_request"`
}
