package in

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

// CreateRemediationBranchCommand carries a caller-supplied RemediationID so
// branch naming and idempotent replay never depend on IssueReference or a
// GitHub Issue existing.
type CreateRemediationBranchCommand struct {
	RemediationID uuid.UUID
	IncidentID    uuid.UUID
	WorkspacePath string
	BaseBranch    string
}

type ExecuteRemediationValidationsCommand struct {
	RemediationID uuid.UUID
	WorkspacePath string
}

// RemediationUseCase intentionally excludes AKR-49's trigger, AKR-51/52's
// change/test generation and AKR-55's failure-decision logic: those are
// deferred, and calling this usecase never requires them to exist.
type RemediationUseCase interface {
	CreateRemediationBranch(context.Context, CreateRemediationBranchCommand) (*domain.Remediation, error)
	ExecuteRemediationValidations(context.Context, ExecuteRemediationValidationsCommand) (*domain.Remediation, []domain.ValidationResult, error)
}
