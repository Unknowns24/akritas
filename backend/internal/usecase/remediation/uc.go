// Package remediation implements PB-041 (dedicated branch creation) and
// PB-044/PB-045 (execute + persist validations) as standalone use cases.
// It deliberately does not implement AKR-49 (triggering Remediation
// creation from Investigation.resolution_status == fixable): callers
// supply a RemediationID and IncidentID explicitly, so this package never
// depends on Investigation, Issue, or IssueReference.
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
