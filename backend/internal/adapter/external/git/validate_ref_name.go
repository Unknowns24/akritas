package git

import "strings"

// validateRefName rejects branch/base names that could be interpreted as a
// git command-line flag (e.g. "-Xupload-pack=...") or that are otherwise
// unsafe as a git ref, BEFORE any exec.Command runs. exec.CommandContext
// with a fixed argv already prevents shell injection; this closes the
// separate git-specific argument-injection surface where a ref name
// starting with "-" is read as an option rather than a literal ref.
func validateRefName(name string) error {
	if name == "" {
		return ErrInvalidWorkspace.Wrap(errInvalidRefName)
	}
	if strings.HasPrefix(name, "-") {
		return ErrInvalidWorkspace.Wrap(errInvalidRefName)
	}
	if name == ".." || strings.Contains(name, "..") {
		return ErrInvalidWorkspace.Wrap(errInvalidRefName)
	}
	for _, r := range name {
		if r <= ' ' || r == 0x7f {
			return ErrInvalidWorkspace.Wrap(errInvalidRefName)
		}
	}
	return nil
}
