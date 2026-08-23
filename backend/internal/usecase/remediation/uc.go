// Package remediation implements dedicated branch creation, validation
// execution and explicit pull-request creation as standalone use cases. It
// deliberately does not implement automatic change generation, merge, deploy
// or rollback.
package remediation

import (
	"context"
	"errors"
	"path/filepath"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
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
	operations        portsout.OperationStore
	investigations    portsout.InvestigationStore
	projections       portsout.ProjectionStore
	policy            *validationpolicy.Policy
	newID             func() uuid.UUID
	now               func() time.Time
}

type RuntimeDependencies struct {
	Operations     portsout.OperationStore
	Investigations portsout.InvestigationStore
	Projections    portsout.ProjectionStore
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
	runtime ...RuntimeDependencies,
) portsin.RemediationUseCase {
	uc := New(workspace, runner, remediations, validationResults, policy, newID, now).(*UseCase)
	uc.incidents = incidents
	uc.projects = projects
	uc.githubAccounts = githubAccounts
	uc.pullRequests = pullRequests
	if len(runtime) > 0 {
		uc.operations = runtime[0].Operations
		uc.investigations = runtime[0].Investigations
		uc.projections = runtime[0].Projections
	}
	return uc
}

func (uc *UseCase) GetIncidentRemediation(ctx context.Context, incidentID uuid.UUID) (*domain.Remediation, error) {
	remediation, err := uc.remediations.FindByIncident(ctx, incidentID)
	if err != nil {
		return nil, err
	}
	return uc.withValidationResults(ctx, remediation)
}

func (uc *UseCase) GetRemediation(ctx context.Context, remediationID uuid.UUID) (*domain.Remediation, error) {
	remediation, err := uc.remediations.Get(ctx, remediationID)
	if err != nil {
		return nil, err
	}
	return uc.withValidationResults(ctx, remediation)
}

func (uc *UseCase) ListValidationResults(ctx context.Context, remediationID uuid.UUID, params paging.Params) (paging.Slice[domain.ValidationResult], error) {
	values, err := uc.validationResults.ListByRemediation(ctx, remediationID)
	if err != nil {
		return paging.Slice[domain.ValidationResult]{}, err
	}
	limit := params.Limit
	if limit < 1 || limit > 100 {
		limit = 25
	}
	total := int64(len(values))
	if len(values) > limit {
		values = values[:limit]
	}
	return paging.Slice[domain.ValidationResult]{Items: values, Total: total}, nil
}

func (uc *UseCase) withValidationResults(ctx context.Context, remediation *domain.Remediation) (*domain.Remediation, error) {
	if remediation == nil {
		return nil, domain.ErrRemediationNotFound
	}
	results, err := uc.validationResults.ListByRemediation(ctx, remediation.ID)
	if err != nil {
		return nil, err
	}
	remediation.ValidationResults = results
	return remediation, nil
}

func (uc *UseCase) ListPullRequests(ctx context.Context, params paging.Params) (paging.Slice[domain.PullRequestProjection], error) {
	if uc.projections == nil {
		return paging.Slice[domain.PullRequestProjection]{}, domain.ErrIntegrationUnavailable
	}
	return uc.projections.ListPullRequests(ctx, params)
}

func (uc *UseCase) GetPullRequest(ctx context.Context, id uuid.UUID) (*domain.PullRequestProjection, error) {
	if uc.projections == nil {
		return nil, domain.ErrIntegrationUnavailable
	}
	return uc.projections.GetPullRequest(ctx, id)
}

func (uc *UseCase) StartIncidentRemediation(ctx context.Context, cmd portsin.StartIncidentRemediationCommand) (*domain.Operation, error) {
	if uc.operations == nil || uc.investigations == nil || uc.projects == nil {
		return nil, domain.ErrIntegrationUnavailable
	}
	key := cmd.IdempotencyKey.String()
	existing, err := uc.operations.FindByIdempotencyKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	incident, err := uc.incidents.Get(ctx, cmd.IncidentID)
	if err != nil {
		return nil, err
	}
	investigation, err := uc.investigations.FindLatestByIncident(ctx, cmd.IncidentID)
	if err != nil {
		return nil, err
	}
	if investigation.Status != domain.InvestigationStatusCompleted || investigation.ResolutionStatus == nil || *investigation.ResolutionStatus != domain.ResolutionFixable {
		return nil, domain.ErrRemediationTransition
	}
	project, err := uc.projects.Get(ctx, incident.ProjectID)
	if err != nil {
		return nil, err
	}
	if cmd.WorkspaceRoot == "" {
		return nil, domain.ErrInvalidRemediation
	}
	workspacePath := filepath.Join(cmd.WorkspaceRoot, project.GitHubRepository.Owner+"__"+project.GitHubRepository.Name)
	remediationID := uc.newID()
	now := uc.now().UTC()
	resourceType := domain.OperationResourceRemediation
	operation, err := domain.NewOperation(uc.newID(), domain.OperationTypeRemediation, &resourceType, &remediationID, &key, "La remediación fue encolada.", now)
	if err != nil {
		return nil, err
	}
	if err := uc.operations.Create(ctx, operation); err != nil {
		return nil, err
	}
	go uc.executeRemediation(context.Background(), operation.ID, portsin.CreateRemediationBranchCommand{
		RemediationID: remediationID, IncidentID: cmd.IncidentID, InvestigationID: investigation.ID,
		WorkspacePath: workspacePath, BaseBranch: project.GitHubRepository.DefaultBranch,
	})
	return operation, nil
}

func (uc *UseCase) QueueRemediationPullRequest(ctx context.Context, cmd portsin.CreateRemediationPullRequestCommand) (*domain.Operation, error) {
	if uc.operations == nil || uc.incidents == nil || uc.projects == nil {
		return nil, domain.ErrIntegrationUnavailable
	}
	key := cmd.IdempotencyKey.String()
	existing, err := uc.operations.FindByIdempotencyKey(ctx, key)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return existing, nil
	}
	remediation, err := uc.remediations.Get(ctx, cmd.RemediationID)
	if err != nil {
		return nil, err
	}
	incident, err := uc.incidents.Get(ctx, remediation.IncidentID)
	if err != nil {
		return nil, err
	}
	project, err := uc.projects.Get(ctx, incident.ProjectID)
	if err != nil {
		return nil, err
	}
	if cmd.WorkspaceRoot == "" {
		return nil, domain.ErrInvalidRemediation
	}
	workspacePath := filepath.Join(cmd.WorkspaceRoot, project.GitHubRepository.Owner+"__"+project.GitHubRepository.Name)
	now := uc.now().UTC()
	resourceType := domain.OperationResourcePullRequest
	operation, err := domain.NewOperation(uc.newID(), domain.OperationTypePullRequest, &resourceType, &cmd.RemediationID, &key, "La creación de Pull Request fue encolada.", now)
	if err != nil {
		return nil, err
	}
	if err := uc.operations.Create(ctx, operation); err != nil {
		return nil, err
	}
	go uc.executePullRequest(context.Background(), operation.ID, cmd.RemediationID, workspacePath)
	return operation, nil
}

func (uc *UseCase) executeRemediation(ctx context.Context, operationID uuid.UUID, cmd portsin.CreateRemediationBranchCommand) {
	operation := uc.startOperation(ctx, operationID)
	if operation == nil {
		return
	}
	remediation, err := uc.CreateRemediationBranch(ctx, cmd)
	if err == nil {
		_, _, err = uc.ExecuteRemediationValidations(ctx, portsin.ExecuteRemediationValidationsCommand{RemediationID: remediation.ID, WorkspacePath: cmd.WorkspacePath})
	}
	uc.finishOperation(ctx, operation, err, "La remediación quedó preparada y validada.", "No se pudo completar la remediación.")
}

func (uc *UseCase) executePullRequest(ctx context.Context, operationID, remediationID uuid.UUID, workspacePath string) {
	operation := uc.startOperation(ctx, operationID)
	if operation == nil {
		return
	}
	_, err := uc.CreateRemediationPullRequest(ctx, portsin.CreateRemediationPullRequestCommand{RemediationID: remediationID, WorkspacePath: workspacePath})
	uc.finishOperation(ctx, operation, err, "La Pull Request fue creada.", "No se pudo crear la Pull Request.")
}

func (uc *UseCase) startOperation(ctx context.Context, id uuid.UUID) *domain.Operation {
	operation, err := uc.operations.FindByID(ctx, id)
	if err != nil || operation.Start(uc.now().UTC()) != nil {
		return nil
	}
	_ = uc.operations.Update(ctx, operation)
	return operation
}

func (uc *UseCase) finishOperation(ctx context.Context, operation *domain.Operation, err error, success, failure string) {
	if err != nil {
		code := domain.ErrIntegrationUnavailable.Code
		var stable *domain.Error
		if errors.As(err, &stable) {
			code = stable.Code
		}
		_ = operation.Fail(uc.now().UTC(), failure, &code)
		_ = uc.operations.Update(ctx, operation)
		return
	}
	_ = operation.Succeed(uc.now().UTC(), success)
	_ = uc.operations.Update(ctx, operation)
}
