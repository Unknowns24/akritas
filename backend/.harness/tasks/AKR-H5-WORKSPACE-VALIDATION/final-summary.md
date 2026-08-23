# Final Summary

## Task

AKR-H5-WORKSPACE-VALIDATION — PB-041/AKR-50 (dedicated branch creation),
PB-044/AKR-53 (execute validations), PB-045/AKR-54 (persist validation
results), scoped deliberately to minimize merge conflicts with the
parallel H4 (Incident → Investigation → Issue → GitHub Issue →
IssueReference) effort.

## What changed

New: `internal/core/ports/out/{repository_workspace,workspace_inspector,
validation_runner,remediation_store,validation_result_store}.go`,
`internal/core/ports/in/remediation.go`,
`internal/service/validationpolicy/*`, `internal/usecase/remediation/*`,
`internal/adapter/external/git/*` (filled the previously-empty
placeholder), `internal/adapter/external/validationrunner/*`,
`internal/adapter/db/postgres/repository/{remediation,validationresult}/*`,
two migrations (`20260823_06_add_remediations`,
`20260823_07_add_validation_results`), the harness task folder itself.

Modified (all additive): `internal/core/domain/errors.go` (+2 sentinels,
+1 catalog func), `internal/errorcatalog/catalog_test.go` (+3 catalogs
registered), `internal/adapter/db/postgres/errors/catalog.go` (+2
sentinels), `internal/adapter/db/postgres/migrations/migrate.go` (+2
migrations registered), `internal/adapter/db/postgres/migrations/migrate_test.go`
(+2 expected IDs — required by the pre-existing exhaustive-list test),
`internal/adapter/db/postgres/dbtest/dbtest.go` (+2 tables in TRUNCATE
list), `docs/errors/aaa-map.md` (+12 rows).

Not touched: `docs/openapi.yaml`, `cmd/main.go`,
`internal/bootstrap/integrations/module.go`, router/global wiring, any
Incident/Investigation/Issue/IssueReference file, any H4 migration.

## Tests run

New tests, all green:
- `internal/service/validationpolicy`: 4 test functions.
- `internal/adapter/external/git`: 3 test functions (`TestClientCreateBranch`
  has 7 subtests covering the full AKR-50 acceptance-criteria list).
- `internal/adapter/external/validationrunner`: 7 test functions.
- `internal/usecase/remediation`: 15 test functions (hand-written fakes,
  no mocking library).
- `internal/adapter/db/postgres/repository/remediation`: 3 test functions
  (real local Postgres via `dbtest.Connect`).
- `internal/adapter/db/postgres/repository/validationresult`: 3 test
  functions (real local Postgres).

`go build ./...`: clean. `go vet ./...`: clean. `go test ./...`: green
except 4 pre-existing failures (`administrator`, `administrator_session`,
`evidence`, `investigation` repository tests) — verified unrelated to this
task: none of those files were touched by this branch, and the same
failures reproduce in isolation (confirmed via `git status` showing zero
changes to those packages, and by re-running each failing package alone).

## Architecture review

See `architecture-review.md`. Result: **pass**. Two deliberate, documented
limitations noted (not defects): the closed `ValidationStatus` enum's
granularity, and a partial-failure idempotency gap in
`CreateRemediationBranch`.

## Security review

See `security-review.md`. Result: **pass**. No path exists for external
input (QVAC, Evidence, repository content) to become a shell command or
arbitrary git argv; verified by dedicated regression tests
(`TestClientCreateBranch/ref_name_injection_attempt_is_rejected_before_any_git_command_runs`,
`TestClientRunUnknownCommand`).

## Remaining risks

1. Migration numbering (`20260823_06`/`_07`) is provisional — coordinate
   with the H4 developer before merge, or expect a rename-on-rebase.
2. Six shared files received additive edits; H4 will likely touch the
   same files for its own errors/migrations — low conflict risk (pure
   appends) but worth flagging.
3. `ValidationStatus`'s closed enum can't distinguish a genuine failure
   from a runner malfunction from a timeout once persisted (only in
   free-text `Summary`/in-memory `ExecutionResult`).
4. Partial-failure idempotency gap between git branch creation and the DB
   write in `CreateRemediationBranch` (see architecture-review.md finding).
5. The minimal `remediations` table's shape (no `Changes`,
   `PullRequestReference`, or lifecycle-transition persistence) needs real
   design work in AKR-55.
6. `check-openapi.sh` could not run (PyYAML not installed in this
   environment) — moot for this task since `docs/openapi.yaml` was not
   touched, but worth fixing the environment before a task that does
   change the OpenAPI file.

## Ready for human review
yes
