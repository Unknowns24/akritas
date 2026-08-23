package git

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
)

// runGit invokes the git binary with a fixed argv array (args) rooted at
// dir. It never constructs or interprets a shell string, so no value
// passed in args can be interpreted as anything other than a literal
// argument to git.
func (c *Client) runGit(ctx context.Context, dir string, args ...string) (stdout, stderr string, err error) {
	cmd := exec.CommandContext(ctx, c.binary, args...)
	cmd.Dir = dir
	var out, errOut bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	err = cmd.Run()
	return strings.TrimSpace(out.String()), strings.TrimSpace(errOut.String()), err
}
