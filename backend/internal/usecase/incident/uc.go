package incident

import portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"

type UseCase struct {
	incidents      portsout.IncidentStore
	investigations portsout.LatestInvestigationFinder
	issueRefs      portsout.GitHubIssueReferenceStore
	timeline       portsout.IncidentTimelineLister
}

func New(
	incidents portsout.IncidentStore,
	investigations portsout.LatestInvestigationFinder,
	issueRefs portsout.GitHubIssueReferenceStore,
	timeline portsout.IncidentTimelineLister,
) *UseCase {
	return &UseCase{incidents: incidents, investigations: investigations, issueRefs: issueRefs, timeline: timeline}
}
