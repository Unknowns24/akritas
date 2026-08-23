package investigationtools

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	githubexternal "github.com/Unknowns24/akritas/backend/internal/adapter/external/github"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

type fakeAPI struct {
	searchQuery string
	readPath    string
}

func (f *fakeAPI) SearchCode(ctx context.Context, account domain.GitHubAccount, owner, repo, query string) ([]githubexternal.CodeSearchMatch, error) {
	f.searchQuery = query
	return []githubexternal.CodeSearchMatch{{Path: "a.go", Repository: owner + "/" + repo}}, nil
}
func (f *fakeAPI) ReadFile(ctx context.Context, account domain.GitHubAccount, owner, repo, path, ref string) (githubexternal.FileContent, error) {
	f.readPath = path
	return githubexternal.FileContent{Path: path, Ref: ref, Content: "ok"}, nil
}
func (f *fakeAPI) ListRecentCommits(ctx context.Context, account domain.GitHubAccount, owner, repo, branch string, limit int) ([]githubexternal.CommitSummary, error) {
	return nil, nil
}
func (f *fakeAPI) ReadCommit(ctx context.Context, account domain.GitHubAccount, owner, repo, sha string) (githubexternal.CommitDetail, error) {
	return githubexternal.CommitDetail{}, nil
}
func (f *fakeAPI) ReadDiff(ctx context.Context, account domain.GitHubAccount, owner, repo, sha string) (string, error) {
	return "", nil
}

func TestRegistryExposesOnlyAllowlistedTools(t *testing.T) {
	t.Parallel()
	api := &fakeAPI{}
	scope := Scope{Account: domain.GitHubAccount{ID: uuid.New()}, Owner: "acme", Name: "app", Branch: "main"}
	registry := Registry(api, scope)
	names := map[string]bool{}
	for _, def := range registry.DefinitionsForTest() {
		names[def] = true
	}
	for _, expected := range []string{"search_code", "read_file", "list_recent_commits", "read_commit", "read_diff"} {
		if !names[expected] {
			t.Fatalf("missing tool %s in %#v", expected, names)
		}
	}
	raw, err := registry.ExecuteForTest(context.Background(), "search_code", `{"query":"TODO"}`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(raw, "a.go") {
		t.Fatalf("result=%s", raw)
	}
	if api.searchQuery != "TODO" {
		t.Fatalf("query=%s", api.searchQuery)
	}
	_ = json.RawMessage(raw)
}
