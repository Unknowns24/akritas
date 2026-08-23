package investigation

import (
	"context"
	"errors"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

// Execute drives an already-queued Investigation through pending->running and
// then to completed or failed, persisting each transition discretely so no
// state is lost if the process is interrupted mid-flight. It is invoked by
// InvestigationDispatcher on its own detached lifetime, never by REST.
func (uc *RunUseCase) Execute(ctx context.Context, investigationID, operationID uuid.UUID) error {
	investigation, err := uc.investigations.FindByID(ctx, investigationID)
	if err != nil {
		return err
	}
	operation, err := uc.operations.FindByID(ctx, operationID)
	if err != nil {
		return err
	}

	startedAt := uc.now().UTC()
	if err := investigation.Start(startedAt); err != nil {
		return err
	}
	if err := uc.investigations.Update(ctx, investigation); err != nil {
		return err
	}
	if err := operation.Start(startedAt); err != nil {
		return err
	}
	if err := uc.operations.Update(ctx, operation); err != nil {
		return err
	}

	assembled, err := uc.assembler.Assemble(ctx, *investigation)
	if err != nil {
		return err
	}
	for i := range assembled {
		if err := uc.evidence.Create(ctx, &assembled[i]); err != nil {
			return err
		}
	}
	investigation.EvidenceCount = len(assembled)
	if err := uc.investigations.Update(ctx, investigation); err != nil {
		return err
	}

	result, runErr := uc.runner.Run(ctx, *investigation)
	finishedAt := uc.now().UTC()
	if runErr != nil {
		return uc.failInvestigation(ctx, investigation, operation, finishedAt, runErr)
	}

	if err := investigation.Complete(
		finishedAt, result.Summary, result.RootCause, result.RootCauseStatus, result.ResolutionStatus,
		result.Confidence, result.Hypotheses, result.RelevantFiles, result.RelevantCommits, result.RecommendedActions,
	); err != nil {
		return uc.failInvestigation(ctx, investigation, operation, finishedAt, err)
	}
	if err := uc.investigations.Update(ctx, investigation); err != nil {
		return err
	}
	if err := operation.Succeed(finishedAt, "La investigación finalizó correctamente."); err != nil {
		return err
	}
	return uc.operations.Update(ctx, operation)
}

func (uc *RunUseCase) failInvestigation(ctx context.Context, investigation *domain.Investigation, operation *domain.Operation, at time.Time, cause error) error {
	message := cause.Error()
	if err := investigation.Fail(at, message); err != nil {
		return err
	}
	if err := uc.investigations.Update(ctx, investigation); err != nil {
		return err
	}

	var failureCode *string
	var stable *domain.Error
	if errors.As(cause, &stable) {
		code := stable.Code
		failureCode = &code
	}
	if err := operation.Fail(at, message, failureCode); err != nil {
		return err
	}
	return uc.operations.Update(ctx, operation)
}
