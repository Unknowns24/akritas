package git

import (
	"context"
	"os"
	"path/filepath"
)

func (c *Client) HasFile(ctx context.Context, workspacePath, relativePath string) (bool, error) {
	info, err := os.Stat(workspacePath)
	if err != nil || !info.IsDir() {
		return false, ErrInvalidWorkspace.Wrap(errWorkspaceUnreadable)
	}
	if _, err := os.Stat(filepath.Join(workspacePath, relativePath)); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, ErrInvalidWorkspace.Wrap(err)
	}
	return true, nil
}
