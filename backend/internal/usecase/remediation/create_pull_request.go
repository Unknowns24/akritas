package remediation

import (
	"context"
	"fmt"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/service/evidencesafety"
)

const publicPullRequestFailure = "No se pudo crear la Pull Request de remediación."

func (uc *UseCase) CreateRemediationPullRequest(ctx context.Context, cmd portsin.CreateRemediationPullRequestCommand) (*domain.Remediation, error) {
	if uc.incidents == nil || uc.projects == nil || uc.githubAccounts == nil || uc.pullRequests == nil {
		return nil, domain.ErrIntegrationUnavailable
	}
	remediation, err := uc.remediations.Get(ctx, cmd.RemediationID)
	if err != nil {
		return nil, err
	}
	if remediation.Status == domain.RemediationStatusPullRequestCreated {
		return remediation, nil
	}
	if remediation.Status != domain.RemediationStatusValidated || remediation.PullRequestReference != nil || remediation.BranchName == "" {
		return nil, domain.ErrRemediationTransition
	}

	incident, err := uc.incidents.Get(ctx, remediation.IncidentID)
	if err != nil {
		return nil, uc.failAfterValidation(ctx, remediation, err)
	}
	project, err := uc.projects.Get(ctx, incident.ProjectID)
	if err != nil {
		return nil, uc.failAfterValidation(ctx, remediation, err)
	}
	account, err := uc.githubAccounts.Get(ctx, project.GitHubRepository.GitHubAccountID)
	if err != nil {
		return nil, uc.failAfterValidation(ctx, remediation, err)
	}

	commit, err := uc.workspace.CommitAll(ctx, portsout.CommitAllInput{
		WorkspacePath: cmd.WorkspacePath,
		BranchName:    remediation.BranchName,
		Message:       commitMessage(remediation),
	})
	if err != nil {
		return nil, uc.failAfterValidation(ctx, remediation, err)
	}
	if _, err = uc.workspace.PushBranch(ctx, portsout.PushBranchInput{
		WorkspacePath: cmd.WorkspacePath,
		BranchName:    remediation.BranchName,
	}); err != nil {
		return nil, uc.failAfterValidation(ctx, remediation, err)
	}

	content := pullRequestContent(*remediation, commit.SHA)
	published, err := uc.pullRequests.CreateOrFindPullRequest(ctx, *account, project.GitHubRepository, portsout.PullRequestInput{
		BaseBranch: project.GitHubRepository.DefaultBranch,
		HeadBranch: remediation.BranchName,
		Content:    content,
	})
	if err != nil {
		return nil, uc.failAfterValidation(ctx, remediation, err)
	}
	reference, err := domain.NewPullRequestReference(
		published.Number, published.URL, project.GitHubRepository.FullName,
		remediation.BranchName, published.CreatedAt,
	)
	if err != nil {
		return nil, uc.failAfterValidation(ctx, remediation, err)
	}
	if err := remediation.CreatePullRequest(reference, uc.now().UTC()); err != nil {
		return nil, err
	}
	if err := uc.remediations.Update(ctx, remediation); err != nil {
		return nil, err
	}
	return remediation, nil
}

func (uc *UseCase) failAfterValidation(ctx context.Context, remediation *domain.Remediation, cause error) error {
	if remediation == nil {
		return cause
	}
	if err := remediation.Fail(publicPullRequestFailure, uc.now().UTC()); err != nil {
		return cause
	}
	if err := uc.remediations.Update(ctx, remediation); err != nil {
		return err
	}
	return cause
}

func commitMessage(remediation *domain.Remediation) string {
	return fmt.Sprintf("AKR-H6 remediation %s", remediation.ID.String())
}

func pullRequestContent(remediation domain.Remediation, commitSHA string) portsout.PullRequestContent {
	body := fmt.Sprintf(`Akritas remediation.

AKRITAS-REMEDIATION-ID: %s
AKRITAS-INCIDENT-ID: %s
AKRITAS-INVESTIGATION-ID: %s
AKRITAS-COMMIT-SHA: %s

STOP: do not merge, deploy or rollback automatically.
`, remediation.ID, remediation.IncidentID, remediation.InvestigationID, commitSHA)
	body = evidencesafety.RedactAndLimit(body, 12000)
	title := fmt.Sprintf("Akritas remediation %s", remediation.ID.String()[:8])
	return portsout.PullRequestContent{Title: title, Body: body}
}
