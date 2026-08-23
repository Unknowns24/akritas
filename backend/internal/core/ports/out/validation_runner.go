package out

import (
	"context"
	"time"
)

// ValidationCommand is a small, closed enum. Akritas's own policy logic
// (internal/service/validationpolicy) is the ONLY thing that selects a
// value from this enum; it is never derived from Evidence, Investigation
// output, repository content or any other untrusted input. Each value maps
// internally, in the adapter, to exactly one fixed, hardcoded argv array.
type ValidationCommand string

const (
	ValidationCommandGoTest  ValidationCommand = "go_test"
	ValidationCommandGoVet   ValidationCommand = "go_vet"
	ValidationCommandGoBuild ValidationCommand = "go_build"
)

// ExecutionOutcome distinguishes "the process ran to completion" (the exit
// code may still be non-zero — that is validation data, not a runner
// error) from "Akritas killed the process because the deadline passed".
type ExecutionOutcome string

const (
	ExecutionOutcomeCompleted ExecutionOutcome = "completed"
	ExecutionOutcomeTimedOut  ExecutionOutcome = "timed_out"
)

type ExecutionResult struct {
	Outcome  ExecutionOutcome
	ExitCode int
	Stdout   string
	Stderr   string
	Duration time.Duration
}

// ValidationRunner executes exactly one ValidationCommand against
// workspacePath, bounded by ctx. It intentionally has NO
// Run(command string)-shaped method: command is a closed enum value, never
// a caller-constructed string, so untrusted input can never become argv.
//
// Run returns a non-nil error ONLY when Akritas's own execution machinery
// failed to observe an outcome at all (unknown command, unreadable
// workspace, could not start the process) — never for a validation that
// ran and failed, and never merely because ctx's deadline was reached
// (that is ExecutionOutcomeTimedOut with err == nil).
type ValidationRunner interface {
	Run(ctx context.Context, command ValidationCommand, workspacePath string) (ExecutionResult, error)
}
