// Package git implements the outbound RepositoryWorkspace and
// WorkspaceInspector ports against a local git working tree, by invoking
// the system `git` binary with fixed argv arrays (never a shell string).
// It never clones or fetches a remote repository: it operates only on a
// workspace path already checked out on disk.
package git

import (
	"errors"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

var _ out.RepositoryWorkspace = (*Client)(nil)
var _ out.WorkspaceInspector = (*Client)(nil)

type Client struct {
	binary string
}

// New returns a Client that invokes the given git binary (typically "git",
// resolved via PATH by exec.CommandContext). binary must be non-empty.
func New(binary string) (*Client, error) {
	if binary == "" {
		return nil, errors.New("git: binary must not be empty")
	}
	return &Client{binary: binary}, nil
}
