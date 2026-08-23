package git

import (
	"bytes"
	"context"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func requireGit(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("skipping: git binary not found on PATH")
	}
}

// newFixtureRepo creates a local git repository with one commit on branch
// "main" and returns its path and the SHA of that commit.
func newFixtureRepo(t *testing.T) (string, string) {
	t.Helper()
	dir := t.TempDir()

	run := func(args ...string) string {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=akritas-test", "GIT_AUTHOR_EMAIL=akritas-test@example.com",
			"GIT_COMMITTER_NAME=akritas-test", "GIT_COMMITTER_EMAIL=akritas-test@example.com",
		)
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out
		if err := cmd.Run(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out.String())
		}
		return strings0(out.String())
	}

	run("init", "--initial-branch=main")
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte("fixture\n"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}
	run("add", "README.md")
	run("commit", "-m", "initial commit")
	sha := run("rev-parse", "HEAD")

	return dir, sha
}

func strings0(s string) string {
	for len(s) > 0 && (s[len(s)-1] == '\n' || s[len(s)-1] == '\r') {
		s = s[:len(s)-1]
	}
	return s
}

func currentBranch(t *testing.T, dir string) string {
	t.Helper()
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("rev-parse: %v", err)
	}
	return strings0(string(out))
}

func TestClientCreateBranch(t *testing.T) {
	requireGit(t)

	t.Run("valid create switches HEAD and returns base commit", func(t *testing.T) {
		dir, baseSHA := newFixtureRepo(t)
		client, err := New("git")
		if err != nil {
			t.Fatalf("New: %v", err)
		}

		output, err := client.CreateBranch(context.Background(), portsout.CreateBranchInput{
			WorkspacePath: dir, BaseBranch: "main", BranchName: "akritas/remediation/test-1",
		})
		if err != nil {
			t.Fatalf("CreateBranch: %v", err)
		}
		if output.BranchName != "akritas/remediation/test-1" {
			t.Fatalf("unexpected branch name: %q", output.BranchName)
		}
		if output.BaseBranch != "main" {
			t.Fatalf("unexpected base branch: %q", output.BaseBranch)
		}
		if output.BaseCommit != baseSHA {
			t.Fatalf("expected base commit %q, got %q", baseSHA, output.BaseCommit)
		}
		if output.CreatedAt.IsZero() {
			t.Fatal("expected CreatedAt to be set")
		}
		if got := currentBranch(t, dir); got != "akritas/remediation/test-1" {
			t.Fatalf("expected HEAD on new branch, got %q", got)
		}
	})

	t.Run("branch name equal to base branch is rejected", func(t *testing.T) {
		dir, _ := newFixtureRepo(t)
		client, _ := New("git")

		_, err := client.CreateBranch(context.Background(), portsout.CreateBranchInput{
			WorkspacePath: dir, BaseBranch: "main", BranchName: "main",
		})
		if !errors.Is(err, ErrProtectedBranchTarget) {
			t.Fatalf("expected ErrProtectedBranchTarget, got %v", err)
		}
		if got := currentBranch(t, dir); got != "main" {
			t.Fatalf("expected HEAD to remain on main, got %q", got)
		}
	})

	t.Run("nonexistent base branch", func(t *testing.T) {
		dir, _ := newFixtureRepo(t)
		client, _ := New("git")

		_, err := client.CreateBranch(context.Background(), portsout.CreateBranchInput{
			WorkspacePath: dir, BaseBranch: "does-not-exist", BranchName: "akritas/remediation/test-2",
		})
		if !errors.Is(err, ErrBaseBranchNotFound) {
			t.Fatalf("expected ErrBaseBranchNotFound, got %v", err)
		}
	})

	t.Run("branch name collision", func(t *testing.T) {
		dir, _ := newFixtureRepo(t)
		client, _ := New("git")
		ctx := context.Background()
		input := portsout.CreateBranchInput{WorkspacePath: dir, BaseBranch: "main", BranchName: "akritas/remediation/test-3"}

		if _, err := client.CreateBranch(ctx, input); err != nil {
			t.Fatalf("first CreateBranch: %v", err)
		}
		run := exec.Command("git", "checkout", "main")
		run.Dir = dir
		if err := run.Run(); err != nil {
			t.Fatalf("checkout main: %v", err)
		}

		_, err := client.CreateBranch(ctx, input)
		if !errors.Is(err, ErrBranchAlreadyExists) {
			t.Fatalf("expected ErrBranchAlreadyExists, got %v", err)
		}
		if got := currentBranch(t, dir); got != "main" {
			t.Fatalf("expected HEAD to remain on main after collision, got %q", got)
		}
	})

	t.Run("ref name injection attempt is rejected before any git command runs", func(t *testing.T) {
		dir, _ := newFixtureRepo(t)
		client, _ := New("git")

		_, err := client.CreateBranch(context.Background(), portsout.CreateBranchInput{
			WorkspacePath: dir, BaseBranch: "main", BranchName: "-Xupload-pack=/bin/sh",
		})
		if !errors.Is(err, ErrInvalidWorkspace) {
			t.Fatalf("expected ErrInvalidWorkspace, got %v", err)
		}
	})

	t.Run("non-git directory", func(t *testing.T) {
		dir := t.TempDir()
		client, _ := New("git")

		_, err := client.CreateBranch(context.Background(), portsout.CreateBranchInput{
			WorkspacePath: dir, BaseBranch: "main", BranchName: "akritas/remediation/test-4",
		})
		if !errors.Is(err, ErrInvalidWorkspace) {
			t.Fatalf("expected ErrInvalidWorkspace, got %v", err)
		}
	})

	t.Run("context timeout leaves the workspace in a valid git state", func(t *testing.T) {
		dir, _ := newFixtureRepo(t)
		client, _ := New("git")

		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()
		time.Sleep(time.Millisecond)

		_, err := client.CreateBranch(ctx, portsout.CreateBranchInput{
			WorkspacePath: dir, BaseBranch: "main", BranchName: "akritas/remediation/test-5",
		})
		if err == nil {
			t.Fatal("expected an error from a canceled context")
		}

		status := exec.Command("git", "status")
		status.Dir = dir
		if err := status.Run(); err != nil {
			t.Fatalf("workspace left in an invalid git state: %v", err)
		}
	})
}
