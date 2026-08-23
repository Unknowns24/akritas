package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func TestCreateOrFindPullRequestReusesStrictExistingHeadBaseAndRepository(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 8, 23, 11, 0, 0, 0, time.UTC)
	createCalls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.URL.Path == "/repos/Unknowns24/akritas/pulls" {
			if r.URL.Query().Get("base") != "main" || r.URL.Query().Get("head") != "Unknowns24:akritas/remediation/123" {
				t.Fatalf("unexpected reconciliation query: %s", r.URL.RawQuery)
			}
			_ = json.NewEncoder(w).Encode([]pullRequestResponseDTO{
				pullRequestDTO(8, "main", "other", "Unknowns24/akritas", createdAt),
				pullRequestDTO(9, "main", "akritas/remediation/123", "Unknowns24/akritas", createdAt),
			})
			return
		}
		if r.Method == http.MethodPost {
			createCalls++
		}
		t.Fatalf("unexpected request = %s %s", r.Method, r.URL.String())
	}))
	defer server.Close()
	account := githubAccount(t)
	credentials := credentialStoreFake{values: map[string][]byte{credentialKey(account.ID, portsout.SecretKindGitHubPAT): []byte("secret")}}
	client := newTestClient(t, server.URL, credentials)
	repository := pullRequestRepository(t, account)

	result, err := client.CreateOrFindPullRequest(context.Background(), account, repository, portsout.PullRequestInput{
		BaseBranch: "main", HeadBranch: "akritas/remediation/123",
		Content: portsout.PullRequestContent{Title: "AKR-H6", Body: "body"},
	})
	if err != nil {
		t.Fatalf("CreateOrFindPullRequest: %v", err)
	}
	if result.Number != 9 || createCalls != 0 {
		t.Fatalf("unexpected result=%+v createCalls=%d", result, createCalls)
	}
}

func TestCreateOrFindPullRequestCreatesWhenMissing(t *testing.T) {
	t.Parallel()
	createdAt := time.Date(2026, 8, 23, 11, 30, 0, 0, time.UTC)
	var payload createPullRequestRequestDTO
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			_ = json.NewEncoder(w).Encode([]pullRequestResponseDTO{})
		case http.MethodPost:
			if r.URL.Path != "/repos/Unknowns24/akritas/pulls" {
				t.Fatalf("path = %s", r.URL.Path)
			}
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatal(err)
			}
			_ = json.NewEncoder(w).Encode(pullRequestDTO(10, "main", "akritas/remediation/456", "Unknowns24/akritas", createdAt))
		default:
			t.Fatalf("unexpected method %s", r.Method)
		}
	}))
	defer server.Close()
	account := githubAccount(t)
	credentials := credentialStoreFake{values: map[string][]byte{credentialKey(account.ID, portsout.SecretKindGitHubPAT): []byte("secret")}}
	client := newTestClient(t, server.URL, credentials)
	repository := pullRequestRepository(t, account)

	result, err := client.CreateOrFindPullRequest(context.Background(), account, repository, portsout.PullRequestInput{
		BaseBranch: "main", HeadBranch: "akritas/remediation/456",
		Content: portsout.PullRequestContent{Title: "AKR-H6", Body: "safe body"},
	})
	if err != nil {
		t.Fatalf("CreateOrFindPullRequest: %v", err)
	}
	if payload.Base != "main" || payload.Head != "akritas/remediation/456" || payload.Title != "AKR-H6" || payload.Body != "safe body" {
		t.Fatalf("unexpected payload: %+v", payload)
	}
	if result.Number != 10 || result.URL != "https://github.com/Unknowns24/akritas/pull/10" || !result.CreatedAt.Equal(createdAt) {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func pullRequestRepository(t *testing.T, account domain.GitHubAccount) domain.GitHubRepository {
	t.Helper()
	repository, err := domain.NewGitHubRepository(account.ID, "42", "Unknowns24", "akritas", "main", true, "https://github.com/Unknowns24/akritas")
	if err != nil {
		t.Fatal(err)
	}
	return repository
}

func pullRequestDTO(number int, base, head, fullName string, createdAt time.Time) pullRequestResponseDTO {
	dto := pullRequestResponseDTO{Number: number, HTMLURL: fmt.Sprintf("https://github.com/Unknowns24/akritas/pull/%d", number), CreatedAt: createdAt}
	dto.Base.Ref = base
	dto.Base.Repo.FullName = fullName
	dto.Head.Ref = head
	dto.Head.Repo.FullName = fullName
	return dto
}
