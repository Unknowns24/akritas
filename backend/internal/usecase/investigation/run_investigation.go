package investigation

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

const (
	maximumPersistedEvidence      = 25
	maximumPersistedEvidenceBytes = 128 << 10
	publicInvestigationFailure    = "No se pudo completar la investigación."
)

// Execute performs database transitions atomically around a slow, bounded
// local QVAC call. No transaction remains open during inference or GitHub I/O.
func (uc *RunUseCase) Execute(ctx context.Context, investigationID, operationID uuid.UUID) error {
	var investigation *domain.Investigation
	var operation *domain.Operation
	startedAt := uc.now().UTC()
	if err := uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		var err error
		investigation, err = uc.investigations.FindByID(txCtx, investigationID)
		if err != nil {
			return err
		}
		operation, err = uc.operations.FindByID(txCtx, operationID)
		if err != nil {
			return err
		}
		if err = investigation.Start(startedAt); err != nil {
			return err
		}
		if err = operation.Start(startedAt); err != nil {
			return err
		}
		if err = uc.investigations.Update(txCtx, investigation); err != nil {
			return err
		}
		return uc.operations.Update(txCtx, operation)
	}); err != nil {
		return err
	}

	runContext, err := uc.assembler.Assemble(ctx, *investigation)
	if err != nil {
		return uc.failInvestigation(ctx, investigation, operation, uc.now().UTC(), err)
	}
	if err = uc.persistEvidence(ctx, investigation, runContext.Evidence); err != nil {
		return uc.failInvestigation(ctx, investigation, operation, uc.now().UTC(), err)
	}
	runContext.Investigation = *investigation

	result, runErr := uc.runner.Run(ctx, runContext)
	if persistErr := uc.persistEvidence(ctx, investigation, result.DiscoveredEvidence); persistErr != nil {
		return uc.failInvestigation(ctx, investigation, operation, uc.now().UTC(), persistErr)
	}
	finishedAt := uc.now().UTC()
	if runErr != nil {
		return uc.failInvestigation(ctx, investigation, operation, finishedAt, runErr)
	}

	completedInvestigation := *investigation
	succeededOperation := *operation
	if err = completedInvestigation.Complete(
		finishedAt, result.Summary, result.RootCause, result.RootCauseStatus, result.ResolutionStatus,
		result.Confidence, result.Hypotheses, result.RelevantFiles, result.RelevantCommits,
		result.RecommendedActions, result.EvidenceIDs,
	); err != nil {
		return uc.failInvestigation(ctx, investigation, operation, finishedAt, err)
	}
	if err = succeededOperation.Succeed(finishedAt, "La investigación finalizó correctamente."); err != nil {
		return uc.failInvestigation(ctx, investigation, operation, finishedAt, err)
	}
	err = uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if updateErr := uc.investigations.Update(txCtx, &completedInvestigation); updateErr != nil {
			return updateErr
		}
		return uc.operations.Update(txCtx, &succeededOperation)
	})
	if err != nil {
		return uc.failInvestigation(ctx, investigation, operation, finishedAt, err)
	}
	return nil
}

func (uc *RunUseCase) persistEvidence(ctx context.Context, investigation *domain.Investigation, candidates []domain.Evidence) error {
	if len(candidates) == 0 {
		return nil
	}
	return uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		existing, err := uc.evidence.ListByInvestigation(txCtx, investigation.ID, paging.Params{Limit: maximumPersistedEvidence})
		if err != nil {
			return err
		}
		count := len(existing.Items)
		used := 0
		seen := make(map[string]struct{}, len(existing.Items)+len(candidates))
		for _, evidence := range existing.Items {
			used += len(evidence.Summary) + len(evidence.Content) + len(evidence.FilePath) + len(evidence.CommitSHA) + len(evidence.Patch)
			seen[evidenceDedupKey(evidence)] = struct{}{}
		}
		created := 0
		for index := range candidates {
			candidate := candidates[index]
			key := evidenceDedupKey(candidate)
			if _, duplicate := seen[key]; duplicate {
				continue
			}
			size := len(candidate.Summary) + len(candidate.Content) + len(candidate.FilePath) + len(candidate.CommitSHA) + len(candidate.Patch)
			if count+created+1 > maximumPersistedEvidence || used+size > maximumPersistedEvidenceBytes {
				return domain.ErrInvalidEvidence
			}
			if err := uc.evidence.Create(txCtx, &candidate); err != nil {
				return err
			}
			seen[key] = struct{}{}
			used += size
			created++
		}
		investigation.EvidenceCount = count + created
		return uc.investigations.Update(txCtx, investigation)
	})
}

func evidenceDedupKey(evidence domain.Evidence) string {
	return fmt.Sprintf("%s\x00%s\x00%s\x00%s", evidence.Type, evidence.FilePath, evidence.CommitSHA, evidence.Summary)
}

func (uc *RunUseCase) failInvestigation(ctx context.Context, investigation *domain.Investigation, operation *domain.Operation, at time.Time, cause error) error {
	message, failureCode := publicFailure(cause)
	persistErr := uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		incident, err := uc.incidents.Lock(txCtx, investigation.IncidentID)
		if err != nil {
			return err
		}
		if err = investigation.Fail(at, message); err != nil {
			return err
		}
		if err = operation.Fail(at, message, failureCode); err != nil {
			return err
		}
		if err = incident.FailInvestigation(); err != nil {
			return err
		}
		if err = uc.investigations.Update(txCtx, investigation); err != nil {
			return err
		}
		if err = uc.operations.Update(txCtx, operation); err != nil {
			return err
		}
		return uc.incidents.Update(txCtx, incident)
	})
	if persistErr != nil {
		return persistErr
	}
	return cause
}

func publicFailure(cause error) (string, *string) {
	var stable *domain.Error
	if errors.As(cause, &stable) {
		code := stable.Code
		return stable.UserMessage, &code
	}
	return publicInvestigationFailure, nil
}
