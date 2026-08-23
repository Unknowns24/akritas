# Architecture Review

## Summary

New capabilities are added strictly behind ports (`RepositoryWorkspace`,
`WorkspaceInspector`, `ValidationRunner`, `RemediationStore`,
`ValidationResultStore`); `internal/core` and `internal/usecase` contain no
`os/exec`, `gorm.io`, or `net/http` imports and declare no REST/DB/adapter
error codes (verified by `check-backend-architecture.sh`). Process
execution (`git` CLI, `go` CLI) lives exclusively in
`internal/adapter/external/*`, matching `backend/external-adapters.md`'s
boundary rule.

## Layering

`core/ports/out` → `service/validationpolicy` (pure decision logic,
depends only on `WorkspaceInspector`) → `usecase/remediation` (depends
only on injected ports + the policy service) → `adapter/external/git`,
`adapter/external/validationrunner`, `adapter/db/postgres/repository/*`
(infrastructure). No usecase depends on a concrete adapter type. No
adapter depends on another adapter.

## Modularity / SRP

`RepositoryWorkspace` exposes exactly one method (`CreateBranch`) — no
anticipatory giant interface for future H5 git operations (commit, push).
`ValidationRunner` exposes exactly one method with a closed
`ValidationCommand` enum, not a general `Run(command string)`. Policy
(what to run) is a separate package from the runner (how to run it) is
separate from persistence (`RemediationStore`/`ValidationResultStore`),
per the "policy ≠ execution ≠ persistence" goal in the task brief.
Repository packages follow the existing project/investigation convention:
one file per method, a private `errors.go` with `mapError`.

## OpenAPI consistency

No changes to `docs/openapi.yaml`. This task adds no REST handlers, so no
new endpoint/DTO surface needed reconciling against the spec. The
`Remediation`/`ValidationResult` domain shapes this task persists remain
consistent with the schemas already defined there (field names, enum
values, and length limits mirrored exactly in the new migrations).

## Findings

None blocking. Two deliberate, documented limitations (not defects):
1. The closed `ValidationStatus` domain enum (`pending/running/passed/failed`)
   cannot distinguish a genuine validation failure from a runner execution
   error or a timeout once persisted — preserved only in the non-persisted
   `ExecutionResult`/`ExecutionOutcome` and in free-text `Summary`.
2. `CreateRemediationBranch` has a partial-failure idempotency gap: if the
   git branch is created but the subsequent `RemediationStore.Create` call
   fails, a retry with the same `RemediationID` hits `ErrBranchAlreadyExists`
   from the adapter rather than adopting the existing branch. Documented in
   `implementation-brief.md`'s Risks section; not addressed in this task.

## Result
pass
