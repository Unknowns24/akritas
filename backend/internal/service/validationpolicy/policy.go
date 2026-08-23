// Package validationpolicy decides WHAT validations to run against a
// workspace (stack detection → closed ValidationPlan). It never executes
// anything itself: execution is out.ValidationRunner's responsibility, kept
// separate so the set of runnable commands stays a closed, Akritas-chosen
// enum rather than something this policy's caller could smuggle input
// into.
package validationpolicy

import (
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

// ValidationStep is one planned, closed-enum validation command paired
// with the domain.ValidationType it will be recorded as.
type ValidationStep struct {
	Type    domain.ValidationType
	Name    string
	Command portsout.ValidationCommand
}

// ValidationPlan is the policy's decision output. Supported=false means the
// repository's stack was not recognized: Steps is empty and the caller
// must not fabricate a passed result.
type ValidationPlan struct {
	Supported bool
	Steps     []ValidationStep
}

type Policy struct {
	inspector portsout.WorkspaceInspector
}

func New(inspector portsout.WorkspaceInspector) *Policy {
	return &Policy{inspector: inspector}
}
