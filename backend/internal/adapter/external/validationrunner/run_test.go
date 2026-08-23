package validationrunner

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func writeModule(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for name, content := range files {
		path := filepath.Join(dir, name)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
	return dir
}

const passingModule = `module fixture

go 1.21
`

func TestClientRunGoTestPassing(t *testing.T) {
	dir := writeModule(t, map[string]string{
		"go.mod": passingModule,
		"ok_test.go": `package fixture

import "testing"

func TestOK(t *testing.T) {}
`,
	})
	client, err := New("go")
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	result, err := client.Run(context.Background(), portsout.ValidationCommandGoTest, dir)
	if err != nil {
		t.Fatalf("Run: %v", err)
	}
	if result.Outcome != portsout.ExecutionOutcomeCompleted {
		t.Fatalf("expected completed outcome, got %v", result.Outcome)
	}
	if result.ExitCode != 0 {
		t.Fatalf("expected exit code 0, got %d\nstdout: %s\nstderr: %s", result.ExitCode, result.Stdout, result.Stderr)
	}
}

func TestClientRunGoTestFailing(t *testing.T) {
	dir := writeModule(t, map[string]string{
		"go.mod": passingModule,
		"fail_test.go": `package fixture

import "testing"

func TestFails(t *testing.T) { t.Fatal("intentional failure") }
`,
	})
	client, _ := New("go")

	result, err := client.Run(context.Background(), portsout.ValidationCommandGoTest, dir)
	if err != nil {
		t.Fatalf("Run should not return an error for a validation failure, got %v", err)
	}
	if result.Outcome != portsout.ExecutionOutcomeCompleted {
		t.Fatalf("expected completed outcome, got %v", result.Outcome)
	}
	if result.ExitCode == 0 {
		t.Fatal("expected a non-zero exit code")
	}
}

func TestClientRunGoBuildCompileError(t *testing.T) {
	dir := writeModule(t, map[string]string{
		"go.mod":     passingModule,
		"broken.go": `package fixture

func broken() { this does not compile }
`,
	})
	client, _ := New("go")

	result, err := client.Run(context.Background(), portsout.ValidationCommandGoBuild, dir)
	if err != nil {
		t.Fatalf("Run should not return an error for a compile failure, got %v", err)
	}
	if result.ExitCode == 0 {
		t.Fatal("expected a non-zero exit code for a compile error")
	}
}

func TestClientRunGoVetFlagsIssue(t *testing.T) {
	dir := writeModule(t, map[string]string{
		"go.mod": passingModule,
		"vet.go": `package fixture

import "fmt"

func BadPrintf() {
	fmt.Printf("%d\n", "not a number")
}
`,
	})
	client, _ := New("go")

	result, err := client.Run(context.Background(), portsout.ValidationCommandGoVet, dir)
	if err != nil {
		t.Fatalf("Run should not return an error for a vet finding, got %v", err)
	}
	if result.ExitCode == 0 {
		t.Fatal("expected a non-zero exit code for a vet finding")
	}
}

func TestClientRunNonexistentWorkspace(t *testing.T) {
	client, _ := New("go")

	_, err := client.Run(context.Background(), portsout.ValidationCommandGoBuild, filepath.Join(t.TempDir(), "does-not-exist"))
	if err == nil {
		t.Fatal("expected an error for a nonexistent workspace path")
	}
}

func TestClientRunTimeout(t *testing.T) {
	dir := writeModule(t, map[string]string{
		"go.mod": passingModule,
		"slow_test.go": `package fixture

import (
	"testing"
	"time"
)

func TestSlow(t *testing.T) { time.Sleep(5 * time.Second) }
`,
	})
	client, _ := New("go")

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	result, err := client.Run(ctx, portsout.ValidationCommandGoTest, dir)
	if err != nil {
		t.Fatalf("Run should not return an error for a timeout, got %v", err)
	}
	if result.Outcome != portsout.ExecutionOutcomeTimedOut {
		t.Fatalf("expected timed out outcome, got %v", result.Outcome)
	}
}

func TestClientRunUnknownCommand(t *testing.T) {
	dir := writeModule(t, map[string]string{"go.mod": passingModule})
	client, _ := New("go")

	_, err := client.Run(context.Background(), portsout.ValidationCommand("not_a_real_command"), dir)
	if err == nil {
		t.Fatal("expected an error for an unknown ValidationCommand")
	}
	if !errors.Is(err, ErrValidationExecutionFailed) {
		t.Fatalf("expected ErrValidationExecutionFailed, got %v", err)
	}
}
