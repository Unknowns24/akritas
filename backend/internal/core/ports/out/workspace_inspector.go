package out

import "context"

// WorkspaceInspector is a narrow, read-only local-filesystem capability
// used exclusively for closed-set stack-marker detection (e.g. "does
// go.mod exist"). It is deliberately NOT a general file-read/browse/list
// capability: it must never be used to feed file contents into command
// selection, which would reopen the injection surface ValidationRunner is
// designed to close.
type WorkspaceInspector interface {
	HasFile(ctx context.Context, workspacePath, relativePath string) (bool, error)
}
