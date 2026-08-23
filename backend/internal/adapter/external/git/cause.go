package git

import "errors"

var (
	errInvalidRefName      = errors.New("ref name is empty, starts with '-', contains '..' or a control character")
	errWorkspaceUnreadable = errors.New("workspace path does not exist or is not a directory")
	errNotAGitRepository   = errors.New("workspace path is not a git repository")
)
