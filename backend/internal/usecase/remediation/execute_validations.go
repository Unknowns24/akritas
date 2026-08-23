package remediation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/service/validationpolicy"
)

const maxOutputExcerptBytes = 50000

// ExecuteRemediationValidations builds a ValidationPlan, runs every step
// regardless of earlier failures (a complete auditable record beats saving
// a few seconds of CI time), and persists each result as it finishes.
//
// It deliberately never calls Remediation.MarkValidated or Remediation.Fail:
// MarkValidated requires at least one CodeChange, which AKR-51 (out of
// scope) would produce, and deciding what a validation outcome means for
// the Remediation's fate is AKR-55's job. Remediation.Status is returned
// unchanged (in_progress).
func (uc *UseCase) ExecuteRemediationValidations(ctx context.Context, cmd portsin.ExecuteRemediationValidationsCommand) (*domain.Remediation, []domain.ValidationResult, error) {
	remediation, err := uc.remediations.Get(ctx, cmd.RemediationID)
	if err != nil {
		return nil, nil, err
	}
	if remediation.Status != domain.RemediationStatusInProgress {
		return nil, nil, domain.ErrRemediationTransition
	}

	plan, err := uc.policy.Plan(ctx, cmd.WorkspacePath)
	if err != nil {
		return nil, nil, err
	}
	if !plan.Supported {
		return nil, nil, domain.ErrValidationStackUnsupported
	}

	results := make([]domain.ValidationResult, 0, len(plan.Steps))
	for _, step := range plan.Steps {
		result, err := domain.NewValidationResult(uc.newID(), remediation.ID, step.Type, step.Name, uc.now())
		if err != nil {
			return remediation, results, err
		}
		if err := result.Start(uc.now()); err != nil {
			return remediation, results, err
		}

		if err := uc.finishFromExecution(ctx, result, step, cmd.WorkspacePath); err != nil {
			return remediation, results, err
		}

		if err := remediation.AddValidationResult(*result, uc.now()); err != nil {
			return remediation, results, err
		}
		if err := uc.validationResults.Create(ctx, result); err != nil {
			return remediation, results, err
		}
		results = append(results, *result)
	}

	return remediation, results, nil
}

func (uc *UseCase) finishFromExecution(ctx context.Context, result *domain.ValidationResult, step validationpolicy.ValidationStep, workspacePath string) error {
	execResult, runErr := uc.runner.Run(ctx, step.Command, workspacePath)
	at := uc.now()

	switch {
	case runErr != nil:
		return result.Fail(at, "Akritas no pudo ejecutar la validación.", truncateWithMarker(runErr.Error(), maxOutputExcerptBytes))
	case execResult.Outcome == portsout.ExecutionOutcomeTimedOut:
		return result.Fail(at, "La validación superó el tiempo máximo permitido.", truncateWithMarker(execResult.Stdout+execResult.Stderr, maxOutputExcerptBytes))
	case execResult.ExitCode == 0:
		return result.Pass(at, "La validación se ejecutó correctamente.", truncateWithMarker(execResult.Stdout, maxOutputExcerptBytes))
	default:
		return result.Fail(at, "La validación finalizó con errores.", truncateWithMarker(execResult.Stdout+execResult.Stderr, maxOutputExcerptBytes))
	}
}
