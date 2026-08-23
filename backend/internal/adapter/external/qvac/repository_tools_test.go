package qvac

import (
	"context"
	"errors"
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type fakeRepositoryInspector struct {
	owner, repo, branch string
	readFileBranch      string
	calls               []string
	fileContent         string
}

func (f *fakeRepositoryInspector) record(owner, repo, branch, call string) {
	f.owner, f.repo, f.branch = owner, repo, branch
	f.calls = append(f.calls, call)
}
func (f *fakeRepositoryInspector) SearchCode(_ context.Context, _ domain.GitHubAccount, owner, repo, query string) ([]portsout.RepositoryCodeMatch, error) {
	f.record(owner, repo, "", "search_code")
	return []portsout.RepositoryCodeMatch{{Path: "db.go", Repository: owner + "/" + repo}}, nil
}
func (f *fakeRepositoryInspector) ReadFile(_ context.Context, _ domain.GitHubAccount, owner, repo, path, ref string) (portsout.RepositoryFile, error) {
	f.record(owner, repo, ref, "read_file")
	f.readFileBranch = ref
	content := f.fileContent
	if content == "" {
		content = "package db"
	}
	return portsout.RepositoryFile{Path: path, Ref: ref, Content: content}, nil
}
func (f *fakeRepositoryInspector) ListRecentCommits(_ context.Context, _ domain.GitHubAccount, owner, repo, branch string, limit int) ([]portsout.RepositoryCommitSummary, error) {
	f.record(owner, repo, branch, "list_recent_commits")
	return []portsout.RepositoryCommitSummary{{SHA: "deadbeef"}}, nil
}
func (f *fakeRepositoryInspector) ReadCommit(_ context.Context, _ domain.GitHubAccount, owner, repo, sha string) (portsout.RepositoryCommit, error) {
	f.record(owner, repo, "", "read_commit")
	return portsout.RepositoryCommit{SHA: sha}, nil
}
func (f *fakeRepositoryInspector) ReadDiff(_ context.Context, _ domain.GitHubAccount, owner, repo, sha string) (string, error) {
	f.record(owner, repo, "", "read_diff")
	return "diff --git a/db.go b/db.go", nil
}

func TestRepositoryToolsExposeOnlyFiveReadsAndBindConfiguredScope(t *testing.T) {
	t.Parallel()
	inspector := &fakeRepositoryInspector{}
	registry := newRepositoryToolRegistry(inspector, portsout.RepositoryScope{Owner: "project-a", Name: "repo-a", Branch: "release"})
	want := []string{"list_recent_commits", "read_commit", "read_diff", "read_file", "search_code"}
	if got := registry.DefinitionsForTest(); len(got) != len(want) {
		t.Fatalf("tools=%v", got)
	} else {
		for index := range want {
			if got[index] != want[index] {
				t.Fatalf("tools=%v", got)
			}
		}
	}
	calls := map[string]string{
		"search_code":         `{"query":"connection refused"}`,
		"read_file":           `{"path":"internal/db.go"}`,
		"list_recent_commits": `{"limit":5}`,
		"read_commit":         `{"sha":"deadbeef"}`,
		"read_diff":           `{"sha":"deadbeef"}`,
	}
	for name, arguments := range calls {
		if _, err := registry.ExecuteForTest(context.Background(), name, arguments); err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if inspector.owner != "project-a" || inspector.repo != "repo-a" {
			t.Fatalf("%s escaped Project A scope: %s/%s", name, inspector.owner, inspector.repo)
		}
	}
	if inspector.readFileBranch != "release" {
		t.Fatalf("default branch was not bound: %q", inspector.readFileBranch)
	}
	if _, err := registry.ExecuteForTest(context.Background(), "write_file", `{}`); !errors.Is(err, ErrUnknownTool) {
		t.Fatalf("unknown mutation must fail closed: %v", err)
	}
}

func TestRepositoryToolsCannotCrossProjectRepositoryScopes(t *testing.T) {
	t.Parallel()
	projectAInspector := &fakeRepositoryInspector{}
	projectBInspector := &fakeRepositoryInspector{}
	projectATools := newRepositoryToolRegistry(projectAInspector, portsout.RepositoryScope{Owner: "owner-a", Name: "repo-a", Branch: "main"})
	projectBTools := newRepositoryToolRegistry(projectBInspector, portsout.RepositoryScope{Owner: "owner-b", Name: "repo-b", Branch: "stable"})

	if _, err := projectATools.ExecuteForTest(context.Background(), "search_code", `{"query":"failure","owner":"owner-b","repo":"repo-b"}`); err == nil {
		t.Fatal("Project A tool accepted a model-selected Project B repository")
	}
	if len(projectAInspector.calls) != 0 {
		t.Fatalf("Project A inspector was invoked after cross-scope arguments: %v", projectAInspector.calls)
	}
	if _, err := projectATools.ExecuteForTest(context.Background(), "search_code", `{"query":"failure"}`); err != nil {
		t.Fatal(err)
	}
	if _, err := projectBTools.ExecuteForTest(context.Background(), "read_file", `{"path":"internal/db.go"}`); err != nil {
		t.Fatal(err)
	}
	if projectAInspector.owner != "owner-a" || projectAInspector.repo != "repo-a" {
		t.Fatalf("Project A escaped scope: %s/%s", projectAInspector.owner, projectAInspector.repo)
	}
	if projectBInspector.owner != "owner-b" || projectBInspector.repo != "repo-b" || projectBInspector.readFileBranch != "stable" {
		t.Fatalf("Project B escaped scope: %s/%s@%s", projectBInspector.owner, projectBInspector.repo, projectBInspector.readFileBranch)
	}
}
