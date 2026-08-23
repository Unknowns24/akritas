package git

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestClientHasFile(t *testing.T) {
	client, err := New("git")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example\n"), 0o644); err != nil {
		t.Fatalf("write fixture: %v", err)
	}

	present, err := client.HasFile(context.Background(), dir, "go.mod")
	if err != nil {
		t.Fatalf("HasFile: %v", err)
	}
	if !present {
		t.Fatal("expected go.mod to be reported present")
	}

	absent, err := client.HasFile(context.Background(), dir, "package.json")
	if err != nil {
		t.Fatalf("HasFile: %v", err)
	}
	if absent {
		t.Fatal("expected package.json to be reported absent")
	}

	_, err = client.HasFile(context.Background(), filepath.Join(dir, "does-not-exist"), "go.mod")
	if err == nil {
		t.Fatal("expected an error for a nonexistent workspace path")
	}
}
