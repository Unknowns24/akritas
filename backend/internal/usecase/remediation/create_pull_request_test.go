package remediation

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/Unknowns24/akritas/backend/internal/service/validationpolicy"
	"github.com/google/uuid"
)

func TestCreateRemediationPullRequestCommitsPushesCreatesPRAndStops(t *testing.T) {
	now := fixedNow()()
	remediations := newFakeRemediationStore()
	value := validatedRemediation(t, uuid.New(), uuid.New(), uuid.New(), now)
	remediations.byID[value.ID] = *value
	workspace := &fakeRepositoryWorkspace{commitOutput: portsout.CommitAllOutput{SHA: "cafebabe", Summary: "M\tmain.go", CreatedAt: now}}
	prs := &fakePullRequestPublisher{result: portsout.PublishedPullRequest{Number: 33, URL: "https://github.com/Unknowns24/akritas/pull/33", CreatedAt: now}}
	uc := newPullRequestUseCase(t, remediations, workspace, prs, value.IncidentID)

	got, err := uc.CreateRemediationPullRequest(context.Background(), portsin.CreateRemediationPullRequestCommand{
		RemediationID: value.ID, WorkspacePath: "/workspace",
	})
	if err != nil {
		t.Fatalf("CreateRemediationPullRequest: %v", err)
	}
	if got.Status != domain.RemediationStatusPullRequestCreated || got.PullRequestReference == nil || got.PullRequestReference.Number != 33 {
		t.Fatalf("unexpected remediation: %+v", got)
	}
	if len(workspace.commitCalls) != 1 || len(workspace.pushCalls) != 1 || len(prs.calls) != 1 {
		t.Fatalf("expected commit, push and PR once; commit=%d push=%d pr=%d", len(workspace.commitCalls), len(workspace.pushCalls), len(prs.calls))
	}
	if prs.calls[0].input.BaseBranch != "main" || prs.calls[0].input.HeadBranch != value.BranchName {
		t.Fatalf("unexpected PR input: %+v", prs.calls[0].input)
	}
	if strings.Contains(prs.calls[0].input.Content.Body, "github_pat_") || !strings.Contains(prs.calls[0].input.Content.Body, "STOP: do not merge") {
		t.Fatalf("PR body must be safe and explicitly stop after creation: %q", prs.calls[0].input.Content.Body)
	}
}

func TestCreateRemediationPullRequestIdempotentWhenAlreadyPersisted(t *testing.T) {
	now := fixedNow()()
	remediations := newFakeRemediationStore()
	value := validatedRemediation(t, uuid.New(), uuid.New(), uuid.New(), now)
	pr, err := domain.NewPullRequestReference(10, "https://github.com/Unknowns24/akritas/pull/10", "Unknowns24/akritas", value.BranchName, now.Add(time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if err := value.CreatePullRequest(pr, now.Add(2*time.Minute)); err != nil {
		t.Fatal(err)
	}
	remediations.byID[value.ID] = *value
	workspace := &fakeRepositoryWorkspace{}
	prs := &fakePullRequestPublisher{}
	uc := newPullRequestUseCase(t, remediations, workspace, prs, value.IncidentID)

	got, err := uc.CreateRemediationPullRequest(context.Background(), portsin.CreateRemediationPullRequestCommand{
		RemediationID: value.ID, WorkspacePath: "/workspace",
	})
	if err != nil {
		t.Fatalf("CreateRemediationPullRequest: %v", err)
	}
	if got.PullRequestReference == nil || got.PullRequestReference.Number != 10 {
		t.Fatalf("unexpected replay result: %+v", got)
	}
	if len(workspace.commitCalls) != 0 || len(workspace.pushCalls) != 0 || len(prs.calls) != 0 {
		t.Fatalf("idempotent replay must not mutate external systems")
	}
}

func TestCreateRemediationPullRequestPushFailureFailsRemediationAndSkipsPR(t *testing.T) {
	now := fixedNow()()
	remediations := newFakeRemediationStore()
	value := validatedRemediation(t, uuid.New(), uuid.New(), uuid.New(), now)
	remediations.byID[value.ID] = *value
	workspace := &fakeRepositoryWorkspace{pushErr: errors.New("push failed")}
	prs := &fakePullRequestPublisher{}
	uc := newPullRequestUseCase(t, remediations, workspace, prs, value.IncidentID)

	err := firstErr(uc.CreateRemediationPullRequest(context.Background(), portsin.CreateRemediationPullRequestCommand{
		RemediationID: value.ID, WorkspacePath: "/workspace",
	}))
	if err == nil {
		t.Fatal("expected push error")
	}
	if len(prs.calls) != 0 {
		t.Fatalf("PR must not be created after push failure")
	}
	updated := remediations.updated[len(remediations.updated)-1]
	if updated.Status != domain.RemediationStatusFailed || updated.FailureUserMessage == "" {
		t.Fatalf("expected failed remediation, got %+v", updated)
	}
}

func validatedRemediation(t *testing.T, remediationID, incidentID, investigationID uuid.UUID, now time.Time) *domain.Remediation {
	t.Helper()
	value, err := domain.NewRemediation(remediationID, incidentID, now)
	if err != nil {
		t.Fatal(err)
	}
	if err := value.AttachInvestigation(investigationID, now.Add(time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := value.Start(remediationBranchName(remediationID), now.Add(2*time.Second)); err != nil {
		t.Fatal(err)
	}
	change, err := domain.NewCodeChange("main.go", domain.CodeChangeModified, "diff --git a/main.go b/main.go")
	if err != nil {
		t.Fatal(err)
	}
	if err := value.AddChange(change, now.Add(3*time.Second)); err != nil {
		t.Fatal(err)
	}
	result, err := domain.NewValidationResult(uuid.New(), remediationID, domain.ValidationTypeTest, "go test", now.Add(4*time.Second))
	if err != nil {
		t.Fatal(err)
	}
	if err := result.Start(now.Add(5 * time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := result.PassWithOutputRedacted(now.Add(6*time.Second), "ok", "safe", false); err != nil {
		t.Fatal(err)
	}
	if err := value.AddValidationResult(*result, now.Add(7*time.Second)); err != nil {
		t.Fatal(err)
	}
	if err := value.MarkValidated(now.Add(8 * time.Second)); err != nil {
		t.Fatal(err)
	}
	return value
}

func newPullRequestUseCase(t *testing.T, remediations *fakeRemediationStore, workspace *fakeRepositoryWorkspace, prs *fakePullRequestPublisher, incidentID uuid.UUID) portsin.RemediationUseCase {
	t.Helper()
	project, account := remediationProject(t)
	return NewWithPullRequests(
		workspace, newFakeValidationRunner(), remediations, &fakeValidationResultStore{},
		fakeIncidentGetter{incident: &domain.Incident{ID: incidentID, ProjectID: project.ID}},
		fakeProjectGetter{project: project},
		fakeGitHubAccountGetter{account: account},
		prs,
		validationpolicy.New(&fakeWorkspaceInspector{}),
		uuid.New, func() time.Time { return fixedNow()().Add(30 * time.Second) },
	)
}

func remediationProject(t *testing.T) (*domain.Project, *domain.GitHubAccount) {
	t.Helper()
	now := fixedNow()()
	account, err := domain.NewGitHubAccount(uuid.New(), "Akritas", domain.GitHubAccountOrganization, domain.GitHubAuthenticationPersonalAccessToken, "Unknowns24", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	repository, err := domain.NewGitHubRepository(account.ID, "42", "Unknowns24", "akritas", "main", true, "https://github.com/Unknowns24/akritas")
	if err != nil {
		t.Fatal(err)
	}
	application, err := domain.NewDokployApplication(uuid.New(), "app", "api", "API", "production", domain.DokployApplicationRunning)
	if err != nil {
		t.Fatal(err)
	}
	source, err := domain.SourceFromApplication(application)
	if err != nil {
		t.Fatal(err)
	}
	project, err := domain.NewProject(uuid.New(), "API", "demo", repository, source, domain.DefaultMonitoringConfiguration(), now)
	if err != nil {
		t.Fatal(err)
	}
	return project, account
}

func firstErr(_ *domain.Remediation, err error) error { return err }

type fakeIncidentGetter struct{ incident *domain.Incident }

func (f fakeIncidentGetter) Get(context.Context, uuid.UUID) (*domain.Incident, error) {
	return f.incident, nil
}

type fakeProjectGetter struct{ project *domain.Project }

func (f fakeProjectGetter) Create(context.Context, *domain.Project) error { return nil }
func (f fakeProjectGetter) Get(context.Context, uuid.UUID) (*domain.Project, error) {
	return f.project, nil
}
func (f fakeProjectGetter) FindByNormalizedName(context.Context, string) (*domain.Project, error) {
	return nil, domain.ErrProjectNotFound
}
func (f fakeProjectGetter) FindByDokploySource(context.Context, domain.DokploySourceSelector) (*domain.Project, error) {
	return nil, domain.ErrProjectNotFound
}
func (f fakeProjectGetter) List(context.Context, paging.Params) (paging.Slice[domain.Project], error) {
	return paging.Slice[domain.Project]{}, nil
}
func (f fakeProjectGetter) Update(context.Context, *domain.Project, time.Time) error { return nil }
func (f fakeProjectGetter) Delete(context.Context, uuid.UUID) error                  { return nil }

type fakeGitHubAccountGetter struct{ account *domain.GitHubAccount }

func (f fakeGitHubAccountGetter) Get(context.Context, uuid.UUID) (*domain.GitHubAccount, error) {
	return f.account, nil
}

type pullRequestCall struct {
	input portsout.PullRequestInput
}

type fakePullRequestPublisher struct {
	result portsout.PublishedPullRequest
	err    error
	calls  []pullRequestCall
}

func (f *fakePullRequestPublisher) CreateOrFindPullRequest(ctx context.Context, account domain.GitHubAccount, repository domain.GitHubRepository, input portsout.PullRequestInput) (portsout.PublishedPullRequest, error) {
	f.calls = append(f.calls, pullRequestCall{input: input})
	if f.err != nil {
		return portsout.PublishedPullRequest{}, f.err
	}
	return f.result, nil
}
