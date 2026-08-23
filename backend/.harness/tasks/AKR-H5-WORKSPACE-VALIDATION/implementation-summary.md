# Implementation Summary

## Task

AKR-H5-WORKSPACE-VALIDATION — PB-041 (dedicated branch), PB-044 (execute
validations), PB-045 (persist validation results).

## What was built

**Ports** (`internal/core/ports/out`): `RepositoryWorkspace` (branch
creation over an existing local git working tree), `WorkspaceInspector`
(closed-set file-existence check), `ValidationRunner` (closed
`ValidationCommand` enum: `go_test`/`go_vet`/`go_build` — no free-string
`Run`), `RemediationStore` (Create+Get), `ValidationResultStore`
(Create+ListByRemediation). Inbound port `RemediationUseCase`
(`internal/core/ports/in/remediation.go`).

**Service** (`internal/service/validationpolicy`): stack detection (Go via
`go.mod`) → closed `ValidationPlan`/`ValidationStep`. Pure logic, depends
only on `WorkspaceInspector`.

**Usecase** (`internal/usecase/remediation`): `CreateRemediationBranch`
(idempotent by `RemediationID`, never depends on `IssueReference`) and
`ExecuteRemediationValidations` (runs every planned step regardless of
earlier failures, persists incrementally, never calls
`Remediation.MarkValidated`/`Fail` since `MarkValidated` requires
`len(Changes) > 0` — impossible without AKR-51 — and the failure decision
is AKR-55's).

**Adapters**:
- `internal/adapter/external/git` — implements `RepositoryWorkspace` +
  `WorkspaceInspector` via the `git` CLI, `exec.CommandContext` with fixed
  argv only, ref-name allowlist validation before any process spawn.
- `internal/adapter/external/validationrunner` — implements
  `ValidationRunner` via a closed `map[ValidationCommand][]string` argv
  table invoking the `go` binary.
- `internal/adapter/db/postgres/repository/{remediation,validationresult}`
  — private record structs (domain entities have no GORM tags), following
  the `credentialstore` pattern.

**Persistence**: two new append-only migrations,
`20260823_06_add_remediations` (minimal: id, incident_id FK→incidents
RESTRICT, status, branch_name, changes_summary, failure_user_message,
timestamps) and `20260823_07_add_validation_results` (id, remediation_id
FK→remediations RESTRICT, type, name, status, timestamps, summary,
output_excerpt, output_redacted, a status/timestamp consistency CHECK).

**Error catalog**: `domain.RemediationErrors()` (component `506`:
`ErrRemediationNotFound`, `ErrValidationStackUnsupported`); DB persistence
sentinels (components `209`/`210`); adapter sentinels in
`internal/adapter/external/git/errors.go` (component `303`) and
`internal/adapter/external/validationrunner/errors.go` (component `304`).
All registered in `internal/errorcatalog/catalog_test.go` and documented in
`docs/errors/aaa-map.md`.

## Tests run

- All new packages: green (see `final-summary.md` for the full list and
  counts).
- Full `go test ./...`: green except 4 pre-existing failures in
  `administrator`, `administrator_session`, `evidence` and `investigation`
  repository tests — confirmed unrelated to this task (none of those
  files were touched; failures reproduce identically in isolation on an
  otherwise-clean tree).
- `go build ./...`, `go vet ./...`: clean.
- `.harness/kernel/scripts/check-backend-architecture.sh`: passed.
- `.harness/kernel/scripts/check-security.sh`: passed.
- `.harness/kernel/scripts/check-openapi.sh`: could not run (PyYAML not
  installed in this environment) — not installed, since `docs/openapi.yaml`
  was not modified by this task, so the gate would be a no-op regardless.

## Deliberately not implemented

AKR-49 (trigger), AKR-51/52 (change/test generation), AKR-55 (failure
decision), AKR-56/57 (commit/PR), any `cmd/main.go`/bootstrap wiring, any
REST endpoint. See `task.md` → "Fuera de alcance".
