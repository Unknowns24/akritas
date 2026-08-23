// Package remediation implements dedicated branch creation, validation
// execution and explicit pull-request creation as standalone use cases. It
// deliberately does not implement automatic change generation, merge, deploy
// or rollback.
package remediation

import (
	"time"

	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/service/validationpolicy"
	"github.com/google/uuid"
)

type UseCase struct {
	workspace         portsout.RepositoryWorkspace
	runner            portsout.ValidationRunner
	remediations      portsout.RemediationStore
	validationResults portsout.ValidationResultStore
	incidents         portsout.IncidentGetter
	projects          portsout.ProjectStore
	githubAccounts    portsout.GitHubAccountReader
	pullRequests      portsout.PullRequestPublisher
	policy            *validationpolicy.Policy
	newID             func() uuid.UUID
	now               func() time.Time
}

func New(
	workspace portsout.RepositoryWorkspace,
	runner portsout.ValidationRunner,
	remediations portsout.RemediationStore,
	validationResults portsout.ValidationResultStore,
	policy *validationpolicy.Policy,
	newID func() uuid.UUID,
	now func() time.Time,
) portsin.RemediationUseCase {
	return &UseCase{
		workspace: workspace, runner: runner, remediations: remediations, validationResults: validationResults,
		policy: policy, newID: newID, now: now,
	}
}

func NewWithPullRequests(
	workspace portsout.RepositoryWorkspace,
	runner portsout.ValidationRunner,
	remediations portsout.RemediationStore,
	validationResults portsout.ValidationResultStore,
	incidents portsout.IncidentGetter,
	projects portsout.ProjectStore,
	githubAccounts portsout.GitHubAccountReader,
	pullRequests portsout.PullRequestPublisher,
	policy *validationpolicy.Policy,
	newID func() uuid.UUID,
	now func() time.Time,
) portsin.RemediationUseCase {
	uc := New(workspace, runner, remediations, validationResults, policy, newID, now).(*UseCase)
	uc.incidents = incidents
	uc.projects = projects
	uc.githubAccounts = githubAccounts
	uc.pullRequests = pullRequests
	return uc
}
