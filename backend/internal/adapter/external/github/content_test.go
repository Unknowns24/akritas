package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func TestReadFileAndSearchCodeScopedToRepository(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/search/code":
			if !strings.Contains(r.URL.Query().Get("q"), "repo:Unknowns24/akritas") {
				t.Fatalf("query = %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"items": []map[string]any{{
					"path": "main.go", "html_url": "https://github.com/Unknowns24/akritas/blob/main/main.go",
					"repository": map[string]any{"full_name": "Unknowns24/akritas"},
				}},
			})
		case strings.HasPrefix(r.URL.Path, "/repos/Unknowns24/akritas/contents/"):
			content := base64.StdEncoding.EncodeToString([]byte("package main\n"))
			_ = json.NewEncoder(w).Encode(map[string]any{
				"type": "file", "path": "main.go", "sha": "abc", "encoding": "base64", "content": content,
			})
		default:
			t.Fatalf("unexpected path %s", r.URL.Path)
		}
	}))
	t.Cleanup(server.Close)

	account := githubAccount(t)
	client := newTestClient(t, server.URL, credentialStoreFake{values: map[string][]byte{
		credentialKey(account.ID, portsout.SecretKindGitHubPAT): []byte("tok"),
	}})
	matches, err := client.SearchCode(context.Background(), account, "Unknowns24", "akritas", "TODO")
	if err != nil || len(matches) != 1 || matches[0].Path != "main.go" {
		t.Fatalf("search = %#v err=%v", matches, err)
	}
	file, err := client.ReadFile(context.Background(), account, "Unknowns24", "akritas", "main.go", "main")
	if err != nil || file.Content != "package main\n" {
		t.Fatalf("file=%#v err=%v", file, err)
	}
}

func TestListAndReadCommitDiff(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/repos/Unknowns24/akritas/commits":
			_ = json.NewEncoder(w).Encode([]map[string]any{{
				"sha": "deadbeef", "html_url": "https://github.com/Unknowns24/akritas/commit/deadbeef",
				"commit": map[string]any{"message": "fix", "author": map[string]any{"name": "dev", "date": "2026-01-01T00:00:00Z"}},
			}})
		case r.URL.Path == "/repos/Unknowns24/akritas/commits/deadbeef":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"sha": "deadbeef", "html_url": "https://github.com/Unknowns24/akritas/commit/deadbeef",
				"commit": map[string]any{"message": "fix", "author": map[string]any{"name": "dev", "date": "2026-01-01T00:00:00Z"}},
				"files":  []map[string]any{{"filename": "main.go", "status": "modified", "patch": "@@ -1 +1 @@\n-a\n+b"}},
			})
		default:
			t.Fatalf("unexpected %s?%s", r.URL.Path, r.URL.RawQuery)
		}
	}))
	t.Cleanup(server.Close)

	account := githubAccount(t)
	client := newTestClient(t, server.URL, credentialStoreFake{values: map[string][]byte{
		credentialKey(account.ID, portsout.SecretKindGitHubPAT): []byte("tok"),
	}})
	commits, err := client.ListRecentCommits(context.Background(), account, "Unknowns24", "akritas", "main", 5)
	if err != nil || len(commits) != 1 || commits[0].SHA != "deadbeef" {
		t.Fatalf("commits=%#v err=%v", commits, err)
	}
	detail, err := client.ReadCommit(context.Background(), account, "Unknowns24", "akritas", "deadbeef")
	if err != nil || detail.Files[0].Filename != "main.go" {
		t.Fatalf("detail=%#v err=%v", detail, err)
	}
	diff, err := client.ReadDiff(context.Background(), account, "Unknowns24", "akritas", "deadbeef")
	if err != nil || !strings.Contains(diff, "main.go") {
		t.Fatalf("diff=%q err=%v", diff, err)
	}
}
