package investigation

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

const (
	maximumPersistedEvidence      = 25
	maximumPersistedEvidenceBytes = 128 << 10
	publicInvestigationFailure    = "No se pudo completar la investigación."
	publicIssuePublicationFailure = "No se pudo publicar la Issue en GitHub."
)

// Execute performs database transitions atomically around a slow, bounded
// local QVAC call. No transaction remains open during inference or GitHub I/O.
func (uc *RunUseCase) Execute(ctx context.Context, investigationID, operationID uuid.UUID) error {
	var investigation *domain.Investigation
	var operation *domain.Operation
	startedAt := uc.now().UTC()
	log.Printf("investigation: starting investigation_id=%s operation_id=%s", investigationID, operationID)
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
		log.Printf("investigation: failed to mark running investigation_id=%s operation_id=%s error=%v", investigationID, operationID, err)
		return err
	}

	runContext, err := uc.assembler.Assemble(ctx, *investigation)
	if err != nil {
		log.Printf("investigation: context assembly failed investigation_id=%s operation_id=%s incident_id=%s error=%v", investigation.ID, operation.ID, investigation.IncidentID, err)
		return uc.failInvestigation(ctx, investigation, operation, uc.now().UTC(), err)
	}
	log.Printf("investigation: context assembled investigation_id=%s operation_id=%s incident_id=%s evidence=%d repository=%s/%s", investigation.ID, operation.ID, investigation.IncidentID, len(runContext.Evidence), runContext.Repository.Owner, runContext.Repository.Name)
	if err = uc.persistEvidence(ctx, investigation, runContext.Evidence); err != nil {
		log.Printf("investigation: initial evidence persistence failed investigation_id=%s operation_id=%s incident_id=%s evidence=%d error=%v", investigation.ID, operation.ID, investigation.IncidentID, len(runContext.Evidence), err)
		return uc.failInvestigation(ctx, investigation, operation, uc.now().UTC(), err)
	}
	runContext.Investigation = *investigation

	result, runErr := uc.runner.Run(ctx, runContext)
	if runErr != nil {
		log.Printf("investigation: runner failed investigation_id=%s operation_id=%s incident_id=%s error=%v cause=%v", investigation.ID, operation.ID, investigation.IncidentID, runErr, rootCause(runErr))
	} else {
		log.Printf("investigation: runner completed investigation_id=%s operation_id=%s incident_id=%s discovered_evidence=%d confidence=%.4f", investigation.ID, operation.ID, investigation.IncidentID, len(result.DiscoveredEvidence), result.Confidence)
	}
	if persistErr := uc.persistEvidence(ctx, investigation, result.DiscoveredEvidence); persistErr != nil {
		log.Printf("investigation: discovered evidence persistence failed investigation_id=%s operation_id=%s incident_id=%s evidence=%d error=%v", investigation.ID, operation.ID, investigation.IncidentID, len(result.DiscoveredEvidence), persistErr)
		return uc.failInvestigation(ctx, investigation, operation, uc.now().UTC(), persistErr)
	}
	finishedAt := uc.now().UTC()
	if runErr != nil {
		return uc.failInvestigation(ctx, investigation, operation, finishedAt, runErr)
	}

	completedInvestigation := *investigation
	if err = completedInvestigation.Complete(
		finishedAt, result.Summary, result.RootCause, result.RootCauseStatus, result.ResolutionStatus,
		result.Confidence, result.Hypotheses, result.RelevantFiles, result.RelevantCommits,
		result.RecommendedActions, result.EvidenceIDs,
	); err != nil {
		return uc.failInvestigation(ctx, investigation, operation, finishedAt, err)
	}

	var incident domain.Incident
	err = uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		locked, lockErr := uc.incidents.Lock(txCtx, completedInvestigation.IncidentID)
		if lockErr != nil {
			return lockErr
		}
		incident = *locked
		if updateErr := incident.StartIssuePublication(result.RootCauseStatus, result.ResolutionStatus, result.Confidence, result.Summary); updateErr != nil {
			return updateErr
		}
		if updateErr := uc.investigations.Update(txCtx, &completedInvestigation); updateErr != nil {
			return updateErr
		}
		return uc.incidents.Update(txCtx, &incident)
	})
	if err != nil {
		log.Printf("investigation: completion persistence failed investigation_id=%s operation_id=%s incident_id=%s error=%v", investigation.ID, operation.ID, investigation.IncidentID, err)
		return uc.failInvestigation(ctx, investigation, operation, finishedAt, err)
	}

	if existing, err := uc.issueRefs.FindByInvestigation(ctx, completedInvestigation.ID); err != nil {
		return uc.failIssuePublication(ctx, &incident, operation, uc.now().UTC(), err)
	} else if existing != nil {
		return uc.finalizeIssuePublication(ctx, &completedInvestigation, operation, *existing)
	}

	project, err := uc.projects.Get(ctx, incident.ProjectID)
	if err != nil {
		return uc.failIssuePublication(ctx, &incident, operation, uc.now().UTC(), err)
	}
	account, err := uc.githubAccounts.Get(ctx, project.GitHubRepository.GitHubAccountID)
	if err != nil {
		return uc.failIssuePublication(ctx, &incident, operation, uc.now().UTC(), err)
	}
	persistedEvidence, err := uc.evidence.ListByInvestigation(ctx, completedInvestigation.ID, paging.Params{Limit: maximumPersistedEvidence})
	if err != nil {
		return uc.failIssuePublication(ctx, &incident, operation, uc.now().UTC(), err)
	}
	content, err := uc.issueContent.BuildIssueContent(portsout.IssueContentInput{
		Project: *project, Incident: incident, Investigation: completedInvestigation, Evidence: persistedEvidence.Items,
	})
	if err != nil {
		return uc.failIssuePublication(ctx, &incident, operation, uc.now().UTC(), err)
	}
	published, err := uc.issuePublisher.PublishIssue(ctx, *account, project.GitHubRepository, content)
	if err != nil {
		log.Printf("investigation: github issue publication failed investigation_id=%s operation_id=%s incident_id=%s repository=%s error=%v", completedInvestigation.ID, operation.ID, incident.ID, project.GitHubRepository.FullName, err)
		return uc.failIssuePublication(ctx, &incident, operation, uc.now().UTC(), err)
	}
	reference, err := domain.NewGitHubIssueReference(incident.ID, completedInvestigation.ID, published.Number, published.URL, project.GitHubRepository.FullName, published.CreatedAt)
	if err != nil {
		return uc.failIssuePublication(ctx, &incident, operation, uc.now().UTC(), err)
	}
	if err = uc.persistNewIssueReference(ctx, &completedInvestigation, operation, reference); err != nil {
		return uc.failIssuePublication(ctx, &incident, operation, uc.now().UTC(), err)
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
	log.Printf("investigation: failing investigation_id=%s operation_id=%s incident_id=%s public_message=%q failure_code=%s cause=%v root_cause=%v", investigation.ID, operation.ID, investigation.IncidentID, message, failureCodeValue(failureCode), cause, rootCause(cause))
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

func (uc *RunUseCase) persistNewIssueReference(ctx context.Context, investigation *domain.Investigation, operation *domain.Operation, reference domain.GitHubIssueReference) error {
	return uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		incident, err := uc.incidents.Lock(txCtx, investigation.IncidentID)
		if err != nil {
			return err
		}
		if err = uc.issueRefs.Create(txCtx, &reference); err != nil {
			return err
		}
		return uc.finalizeIssuePublicationInTransaction(txCtx, incident, operation, reference)
	})
}

func (uc *RunUseCase) finalizeIssuePublication(ctx context.Context, investigation *domain.Investigation, operation *domain.Operation, reference domain.GitHubIssueReference) error {
	return uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		incident, err := uc.incidents.Lock(txCtx, investigation.IncidentID)
		if err != nil {
			return err
		}
		return uc.finalizeIssuePublicationInTransaction(txCtx, incident, operation, reference)
	})
}

func (uc *RunUseCase) finalizeIssuePublicationInTransaction(ctx context.Context, incident *domain.Incident, operation *domain.Operation, reference domain.GitHubIssueReference) error {
	at := uc.now().UTC()
	if err := incident.AttachGitHubIssue(reference); err != nil {
		return err
	}
	if incident.ResolutionStatus != nil && *incident.ResolutionStatus == domain.ResolutionRequiresHuman {
		if err := incident.CompleteRequiresHuman(); err != nil {
			return err
		}
	}
	if err := operation.Succeed(at, "La Issue de GitHub fue publicada correctamente."); err != nil {
		return err
	}
	if err := uc.incidents.Update(ctx, incident); err != nil {
		return err
	}
	return uc.operations.Update(ctx, operation)
}

func (uc *RunUseCase) failIssuePublication(ctx context.Context, incident *domain.Incident, operation *domain.Operation, at time.Time, cause error) error {
	message, failureCode := publicFailure(cause)
	if message == publicInvestigationFailure {
		message = publicIssuePublicationFailure
	}
	log.Printf("investigation: failing issue publication operation_id=%s incident_id=%s public_message=%q failure_code=%s cause=%v root_cause=%v", operation.ID, incident.ID, message, failureCodeValue(failureCode), cause, rootCause(cause))
	persistErr := uc.transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		locked, err := uc.incidents.Lock(txCtx, incident.ID)
		if err != nil {
			return err
		}
		if locked.Phase == domain.IncidentPhasePublishingIssue {
			if err = locked.FailIssuePublication(); err != nil {
				return err
			}
			if err = uc.incidents.Update(txCtx, locked); err != nil {
				return err
			}
		}
		if operation.Status == domain.OperationStatusRunning {
			if err = operation.Fail(at, message, failureCode); err != nil {
				return err
			}
			if err = uc.operations.Update(txCtx, operation); err != nil {
				return err
			}
		}
		return nil
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

func failureCodeValue(code *string) string {
	if code == nil {
		return ""
	}
	return *code
}

func rootCause(err error) error {
	for {
		unwrapped := errors.Unwrap(err)
		if unwrapped == nil {
			return err
		}
		err = unwrapped
	}
}
