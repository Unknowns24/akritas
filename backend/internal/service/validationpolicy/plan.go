package validationpolicy

import (
	"context"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

// Plan detects the repository's stack and returns a closed ValidationPlan.
// MVP scope: Go, detected by the presence of go.mod at the workspace root.
// Any other (or absent) stack yields Supported=false with no steps — never
// a fabricated success.
func (p *Policy) Plan(ctx context.Context, workspacePath string) (ValidationPlan, error) {
	if err := ctx.Err(); err != nil {
		return ValidationPlan{}, err
	}

	hasGoMod, err := p.inspector.HasFile(ctx, workspacePath, "go.mod")
	if err != nil {
		return ValidationPlan{}, err
	}
	if !hasGoMod {
		return ValidationPlan{Supported: false}, nil
	}

	return ValidationPlan{
		Supported: true,
		Steps: []ValidationStep{
			{Type: domain.ValidationTypeBuild, Name: "go build ./...", Command: portsout.ValidationCommandGoBuild},
			{Type: domain.ValidationTypeStaticAnalysis, Name: "go vet ./...", Command: portsout.ValidationCommandGoVet},
			{Type: domain.ValidationTypeTest, Name: "go test ./...", Command: portsout.ValidationCommandGoTest},
		},
	}, nil
}
