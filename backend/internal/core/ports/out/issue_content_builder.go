package out

import "github.com/Unknowns24/akritas/backend/internal/core/domain"

type IssueContentInput struct {
	Project       domain.Project
	Incident      domain.Incident
	Investigation domain.Investigation
	Evidence      []domain.Evidence
}

type IssueContentBuilder interface {
	BuildIssueContent(IssueContentInput) (IssueContent, error)
}
