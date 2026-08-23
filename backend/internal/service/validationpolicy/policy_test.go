package validationpolicy

import (
	"context"
	"errors"
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

type fakeInspector struct {
	has map[string]bool
	err error
}

func (f *fakeInspector) HasFile(ctx context.Context, workspacePath, relativePath string) (bool, error) {
	if f.err != nil {
		return false, f.err
	}
	return f.has[relativePath], nil
}

func TestPolicyPlanGoModPresent(t *testing.T) {
	inspector := &fakeInspector{has: map[string]bool{"go.mod": true}}
	policy := New(inspector)

	plan, err := policy.Plan(context.Background(), "/workspace")
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if !plan.Supported {
		t.Fatal("expected plan to be supported")
	}
	if len(plan.Steps) != 3 {
		t.Fatalf("expected 3 steps, got %d: %+v", len(plan.Steps), plan.Steps)
	}

	want := map[domain.ValidationType]portsout.ValidationCommand{
		domain.ValidationTypeBuild:          portsout.ValidationCommandGoBuild,
		domain.ValidationTypeStaticAnalysis: portsout.ValidationCommandGoVet,
		domain.ValidationTypeTest:           portsout.ValidationCommandGoTest,
	}
	seen := map[domain.ValidationType]bool{}
	for _, step := range plan.Steps {
		command, ok := want[step.Type]
		if !ok {
			t.Fatalf("unexpected validation type in plan: %v", step.Type)
		}
		if step.Command != command {
			t.Fatalf("type %v: expected command %v, got %v", step.Type, command, step.Command)
		}
		if step.Name == "" {
			t.Fatalf("step %+v missing a Name", step)
		}
		seen[step.Type] = true
	}
	for wantType := range want {
		if !seen[wantType] {
			t.Fatalf("plan missing step for type %v", wantType)
		}
	}
}

func TestPolicyPlanGoModAbsent(t *testing.T) {
	inspector := &fakeInspector{has: map[string]bool{}}
	policy := New(inspector)

	plan, err := policy.Plan(context.Background(), "/workspace")
	if err != nil {
		t.Fatalf("Plan: %v", err)
	}
	if plan.Supported {
		t.Fatal("expected plan to be unsupported")
	}
	if len(plan.Steps) != 0 {
		t.Fatalf("expected no steps, got %d", len(plan.Steps))
	}
}

func TestPolicyPlanInspectorError(t *testing.T) {
	wantErr := errors.New("boom")
	inspector := &fakeInspector{err: wantErr}
	policy := New(inspector)

	_, err := policy.Plan(context.Background(), "/workspace")
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected inspector error to propagate, got %v", err)
	}
}

func TestPolicyPlanRespectsContextCancellation(t *testing.T) {
	inspector := &fakeInspector{has: map[string]bool{"go.mod": true}}
	policy := New(inspector)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := policy.Plan(ctx, "/workspace")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
}
