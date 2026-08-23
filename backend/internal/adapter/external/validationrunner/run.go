package validationrunner

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (c *Client) Run(ctx context.Context, command out.ValidationCommand, workspacePath string) (out.ExecutionResult, error) {
	args, ok := commandArgs[command]
	if !ok {
		return out.ExecutionResult{}, ErrValidationExecutionFailed.Wrap(fmt.Errorf("unknown validation command: %s", command))
	}

	cmd := exec.CommandContext(ctx, c.goBinary, args...)
	cmd.Dir = workspacePath
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	runErr := cmd.Run()
	duration := time.Since(start)

	if runErr != nil {
		if ctx.Err() != nil {
			return out.ExecutionResult{
				Outcome: out.ExecutionOutcomeTimedOut, Stdout: stdout.String(), Stderr: stderr.String(), Duration: duration,
			}, nil
		}
		if exitErr, isExitErr := runErr.(*exec.ExitError); isExitErr {
			return out.ExecutionResult{
				Outcome: out.ExecutionOutcomeCompleted, ExitCode: exitErr.ExitCode(),
				Stdout: stdout.String(), Stderr: stderr.String(), Duration: duration,
			}, nil
		}
		return out.ExecutionResult{}, ErrValidationExecutionFailed.Wrap(runErr)
	}

	return out.ExecutionResult{
		Outcome: out.ExecutionOutcomeCompleted, ExitCode: 0, Stdout: stdout.String(), Stderr: stderr.String(), Duration: duration,
	}, nil
}
