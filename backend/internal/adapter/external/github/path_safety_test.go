package github

import "testing"

func TestSanitizeRepositoryPathRejectsTraversal(t *testing.T) {
	t.Parallel()
	for _, raw := range []string{"../secret", "/etc/passwd", `..\windows`, "a/../../b"} {
		if _, err := sanitizeRepositoryPath(raw); err == nil {
			t.Fatalf("expected error for %q", raw)
		}
	}
}

func TestSanitizeRepositoryPathAcceptsNestedFile(t *testing.T) {
	t.Parallel()
	got, err := sanitizeRepositoryPath("cmd/main.go")
	if err != nil || got != "cmd/main.go" {
		t.Fatalf("got %q err %v", got, err)
	}
}
