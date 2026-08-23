// Package validationrunner implements the outbound ValidationRunner port by
// invoking the local `go` toolchain with a fixed, hardcoded argv per
// portsout.ValidationCommand. It never accepts a caller-constructed command
// string: the set of runnable commands is closed and chosen entirely by
// Akritas's own policy code (internal/service/validationpolicy), never by
// QVAC output, Evidence, or repository content.
package validationrunner

import (
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

var _ out.ValidationRunner = (*Client)(nil)

type Client struct {
	goBinary string
}

// New returns a Client that invokes the given go binary (typically "go",
// resolved via PATH by exec.CommandContext). goBinary must be non-empty.
func New(goBinary string) (*Client, error) {
	if goBinary == "" {
		return nil, errors.New("validationrunner: go binary must not be empty")
	}
	return &Client{goBinary: goBinary}, nil
}

// commandArgs is the closed, hardcoded argv table. It is the only place a
// portsout.ValidationCommand is translated into process arguments.
var commandArgs = map[out.ValidationCommand][]string{
	out.ValidationCommandGoTest:  {"test", "./..."},
	out.ValidationCommandGoVet:   {"vet", "./..."},
	out.ValidationCommandGoBuild: {"build", "./..."},
}
