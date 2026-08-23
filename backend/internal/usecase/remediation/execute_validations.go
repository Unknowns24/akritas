package remediation

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/service/evidencesafety"
	"github.com/Unknowns24/akritas/backend/internal/service/validationpolicy"
)

const maxOutputExcerptBytes = 50000

// ExecuteRemediationValidations builds a ValidationPlan, runs every step
// regardless of earlier failures (a complete auditable record beats saving
// a few seconds of CI time), and persists each result as it finishes.
//
// A failed validation is terminal for this Remediation attempt: all available
// results are persisted first, then the Remediation is marked failed so later
// commit/push/PR stages cannot proceed on known-bad changes.
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
	hasFailedValidation := false
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
		if result.Status == domain.ValidationStatusFailed {
			hasFailedValidation = true
		}
		results = append(results, *result)
	}
	if hasFailedValidation {
		if err := remediation.Fail("La remediación falló porque una o más validaciones no pasaron.", uc.now()); err != nil {
			return remediation, results, err
		}
		if err := uc.remediations.Update(ctx, remediation); err != nil {
			return remediation, results, err
		}
	}

	return remediation, results, nil
}

func (uc *UseCase) finishFromExecution(ctx context.Context, result *domain.ValidationResult, step validationpolicy.ValidationStep, workspacePath string) error {
	execResult, runErr := uc.runner.Run(ctx, step.Command, workspacePath)
	at := uc.now()

	switch {
	case runErr != nil:
		output := safeValidationOutput(runErr.Error())
		return result.FailWithOutputRedacted(at, "Akritas no pudo ejecutar la validación.", output.Value, output.Redacted)
	case execResult.Outcome == portsout.ExecutionOutcomeTimedOut:
		output := safeValidationOutput(execResult.Stdout + execResult.Stderr)
		return result.FailWithOutputRedacted(at, "La validación superó el tiempo máximo permitido.", output.Value, output.Redacted)
	case execResult.ExitCode == 0:
		output := safeValidationOutput(execResult.Stdout)
		return result.PassWithOutputRedacted(at, "La validación se ejecutó correctamente.", output.Value, output.Redacted)
	default:
		output := safeValidationOutput(execResult.Stdout + execResult.Stderr)
		return result.FailWithOutputRedacted(at, "La validación finalizó con errores.", output.Value, output.Redacted)
	}
}

func safeValidationOutput(value string) evidencesafety.RedactionResult {
	result := evidencesafety.RedactAndLimitWithReport(value, maxOutputExcerptBytes)
	result.Value = truncateWithMarker(result.Value, maxOutputExcerptBytes)
	return result
}
