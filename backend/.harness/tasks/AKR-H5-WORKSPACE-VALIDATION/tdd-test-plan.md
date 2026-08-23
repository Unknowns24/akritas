# TDD Test Plan

## Scope

PB-041 (dedicated branch creation), PB-044 (execute validations), PB-045
(persist validation results). See `task.md` for exclusions.

## Tests to add/update

**`internal/service/validationpolicy`** (fake `WorkspaceInspector`):
`go.mod` present → 3-step Go plan (test/build/static_analysis mapped to
`go_test`/`go_build`/`go_vet`); `go.mod` absent → `Supported=false`, empty
steps; inspector error propagates; context cancellation respected.

**`internal/adapter/external/git`** (real local git repo via `t.TempDir()`
+ `git init` + one commit; `t.Skip` if `git` binary missing):
valid create → branch exists, HEAD switched, `BaseCommit` matches base SHA;
`BranchName == BaseBranch` → `ErrProtectedBranchTarget`, no mutation;
nonexistent base branch → `ErrBaseBranchNotFound`; branch name collision →
`ErrBranchAlreadyExists`, state unchanged; ref-name argument-injection
attempt (leading `-`) → rejected by `validateRefName` before any
`exec.Command` runs; non-git directory → `ErrInvalidWorkspace`; context
timeout → error surfaces, workspace left in a valid git state; `HasFile`
true/false/nonexistent cases; `validateRefName` table-driven unit tests.

**`internal/adapter/external/validationrunner`** (real `go` binary, inline
fixture modules): `go_test` passing/failing → `ExecutionOutcomeCompleted`,
correct exit code, `err == nil` either way; `go_build` compile error,
`go_vet` flagged issue → same shape; nonexistent workspace path →
`err != nil`, no fabricated result; short timeout vs. slow fixture →
`ExecutionOutcomeTimedOut`, `err == nil`; unknown `ValidationCommand` →
`err != nil`.

**`internal/usecase/remediation`** (hand-written fakes for every port):
`CreateRemediationBranch` — new ID creates+starts+persists; existing ID
replays idempotently without calling the workspace adapter again; workspace
failure leaves no persisted row; store failure after a successful branch
creation surfaces the divergence. `ExecuteRemediationValidations` — not
found → `ErrRemediationNotFound`; wrong status → `domain.ErrRemediationTransition`
via `AddValidationResult`'s own guard; unsupported stack →
`ErrValidationStackUnsupported`, zero persistence calls; all steps run
regardless of earlier failures; runner execution error persisted as `Fail`
with a summary distinguishable from a genuine test-failure summary; timeout
persisted as `Fail` with a distinguishable summary; oversized
summary/output truncated with a visible marker; `OutputRedacted` always
`true`; mid-plan `store.Create` failure surfaces immediately while
already-persisted prior results remain. Pure unit tests for
`remediationBranchName`/`truncateWithMarker`.

**`internal/adapter/db/postgres/repository/remediation` and
`.../validationresult`** (real Postgres via `dbtest.Connect(t)`): seed an
`incidents` row, `Create`+`Get` round-trip a `Remediation`; `Get` on
unknown ID → `ErrRemediationNotFound`; duplicate-ID `Create` → mapped
conflict, no leaked driver error; seed a `remediations` row, `Create`+
`ListByRemediation` round-trip several `ValidationResult`s across all
statuses with stable ordering; FK violation → `ErrValidationResultPersistence`.

## Expected failing tests before implementation

All of the above fail to compile (packages/types/functions don't exist yet)
until each corresponding production file is added — standard TDD red state
per package, verified incrementally as each package is built.

## Acceptance criteria covered

See `task.md` → "Criterios de aceptación" (maps 1:1 to the AKR-50/53/54
acceptance criteria the user specified: branch creation invariants,
validation execution determinism/security, validation result persistence
and reconstruction).

## Open questions / human approval notes

Resolved during plan-mode review (see approved plan and `AskUserQuestion`
answers): add a minimal `remediations` table now; run all validation steps
regardless of earlier failures; `CreateDedicatedBranch` operates on an
existing local workspace path (no cloning in this task).
