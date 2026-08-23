# Implementation Brief

## Task

AKR-H5-WORKSPACE-VALIDATION — PB-041 (dedicated branch), PB-044 (execute
validations), PB-045 (persist validation results). See `task.md`.

## Current project context

`domain.Remediation`/`domain.ValidationResult`/`domain.CodeChange`/
`domain.PullRequestReference` and their error sentinels (component `406`)
already exist and are fully tested in `internal/core/domain`. No port,
usecase, adapter or table exists above the domain layer for any of them.
`internal/adapter/external/git/` is an empty placeholder. No process
execution (`os/exec`) infrastructure exists anywhere. `docs/openapi.yaml`
already contracts the relevant schemas — untouched by this task.

## Proposed approach

Ports-and-adapters, mirroring `internal/usecase/investigation` and
`internal/adapter/db/postgres/repository/project`:

- `internal/core/ports/out`: `RepositoryWorkspace` (branch creation),
  `WorkspaceInspector` (closed-set file-existence check), `ValidationRunner`
  (closed `ValidationCommand` enum, no free-string `Run`), `RemediationStore`
  (Create+Get), `ValidationResultStore` (Create+ListByRemediation).
- `internal/service/validationpolicy`: stack detection (`go.mod` presence)
  → `ValidationPlan`/`ValidationStep`. Pure logic, depends only on
  `WorkspaceInspector`.
- `internal/usecase/remediation`: `CreateRemediationBranch` and
  `ExecuteRemediationValidations`, depending only on injected ports.
- `internal/adapter/external/git`: `git` CLI via `exec.CommandContext` with
  fixed argv arrays (never `sh -c`), ref-name allowlist validation before
  any process spawn.
- `internal/adapter/external/validationrunner`: closed
  `map[ValidationCommand][]string` argv table, `go` binary invocation with
  context timeout and bounded output capture.
- `internal/adapter/db/postgres/repository/{remediation,validationresult}`:
  private record structs (domain entities have no GORM tags), following the
  `credentialstore` pattern.

## Architecture impact

New ports/adapters/usecases only; no existing package's public API changes.
`internal/adapter/external/git` goes from an empty placeholder to a real
adapter. No REST layer changes.

## API/OpenAPI impact

None. No REST handlers added in this task; existing OpenAPI schemas for
`Remediation`/`ValidationResult` are reused as-is by the domain layer they
already govern.

## Data/persistence impact

Two new tables via two new append-only migrations:
`remediations` (minimal: id, incident_id FK→incidents RESTRICT, status,
branch_name, changes_summary, failure_user_message, timestamps) and
`validation_results` (id, remediation_id FK→remediations RESTRICT, type,
name, status, timestamps, summary, output_excerpt, output_redacted, a
status/timestamp consistency CHECK). `migrate.go` and `dbtest.go`'s
TRUNCATE list gain additive entries.

## Error handling impact

New domain usecase-layer sentinels `ErrRemediationNotFound` (`0x506001N`),
`ErrValidationStackUnsupported` (`0x506002V`) in a new
`domain.RemediationErrors()` catalog. New DB persistence sentinels
`ErrRemediationPersistence` (`0x209001I`), `ErrValidationResultPersistence`
(`0x210001I`) in the existing `db/postgres/errors` catalog. New adapter
sentinels in `internal/adapter/external/git/errors.go` (component `303`)
and `internal/adapter/external/validationrunner/errors.go` (component
`304`). All registered in `internal/errorcatalog/catalog_test.go`'s
`catalogs` slice and documented once each in `docs/errors/aaa-map.md`.

## Test strategy

Hand-written fakes (no mocking framework) for usecase tests, mirroring
`internal/usecase/investigation/fakes_test.go`. Real local `git` CLI
(`t.Skip` if unavailable) for the git adapter, real `go` binary for the
validation runner adapter, real local/testcontainers Postgres via
`dbtest.Connect(t)` for repository tests. See `tdd-test-plan.md`.

## Risks

See the approved plan's Risks section (migration renumbering on rebase with
H4; shared-file additive edits; closed `ValidationStatus` enum can't
distinguish runner-malfunction/timeout from genuine failure once persisted;
partial-failure idempotency gap between git branch creation and the DB
write in `CreateRemediationBranch`; minimal `remediations` table shape will
need extension work in AKR-55).

## Files likely to change

See `task.md` → "Archivos o zonas probablemente afectadas" and the approved
plan's file tree.
