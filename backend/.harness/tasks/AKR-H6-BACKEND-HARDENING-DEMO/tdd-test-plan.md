# TDD Test Plan

## Scope

Generate RED tests first for `AKR-60` through `AKR-65`, including baseline blockers and indispensable H5 prerequisites. Do not implement production code until this plan is approved.

## Baseline and migration tests

### RED-001 - DB error catalog compiles and registers DB errors

Package: `internal/adapter/db/postgres/errors`

- Test that `Catalog()` includes `ErrValidationResultPersistence` and `ErrGitHubIssueReferencePersistence` as distinct non-nil entries.
- Test that every catalog entry has a unique code and no nil pointer.
- Expected initial state: package does not compile because `catalog.go:40` has an invalid nested `var`.

### RED-002 - Global error catalog includes all layer-owned errors

Package: `internal/errorcatalog`

- Extend catalog tests to include any new H6 errors.
- Assert DB/external/usecase/domain errors remain owned by their layer.
- Assert no core catalog declares REST/DB/external adapter codes.

### RED-003 - Migration registry rejects date/order slot collisions

Package: `internal/adapter/db/postgres/migrations`

- Add a registry test that parses `YYYYMMDD_NN_*`.
- Fail if two migrations share the same date/order slot unless an explicit documented historical exception is listed.
- Fail if registry order is not monotonic by date/order.
- Expected initial state: fails on `20260823_06_*` and `20260823_07_*` conflicts.

### RED-004 - Migration dependency order is correct

Package: `internal/adapter/db/postgres/migrations`

- Assert `github_issue_references` exists before any migration that adds constraints to it.
- Assert `incidents` and `investigations` exist before issue references.
- Assert `remediations` exists before `validation_results`.
- Assert any future full-remediation migration runs after H4/H5 historical tables.

### RED-005 - Migration upgrade and rollback on real PostgreSQL

Package: `internal/adapter/db/postgres/migrations`

- Use existing `dbtest` conventions.
- Apply all migrations to empty DB and assert expected H6 tables/indexes/constraints.
- Roll back reversible migrations where supported and assert tables/constraints are removed safely.
- Include legacy-state test for a DB with historical H4/H5 migration IDs already applied.
- Assert migration errors propagate.

### RED-006 - OpenAPI/route correspondence

Packages: `internal/adapter/rest/router`, possibly `internal/adapter/rest`

- Parse `docs/openapi.yaml`.
- Register real router with fake handlers/usecases.
- Assert every non-callback API path declared for QVAC, Remediation and Pull Request has a concrete route or is explicitly marked not implemented by contract decision.
- Assert routes requiring `Idempotency-Key` in OpenAPI enforce it at handler/usecase boundary.
- Expected initial state: fails for QVAC/Remediation/Pull Request route gaps.

## AKR-60 - QVAC never receives credentials

### RED-010 - Safe QVAC context excludes credentials by construction

Package: new or existing safe-context service, plus `internal/usecase/investigation`

- Build Investigation context from Project/GitHub/Dokploy/admin/session fixtures containing unique sentinel secrets:
  - GitHub PAT classic and fine-grained
  - GitHub App private key
  - GitHub App JWT/installation token
  - Dokploy API credential
  - PostgreSQL DSN
  - session token/cookie
  - bootstrap token
  - TOTP seed
  - master key
- Assert QVAC run context contains only IDs and safe metadata.
- Assert no sentinel appears in prompt, tool definitions, tool arguments, tool outputs, Evidence, timeline or persisted Investigation fields.

### RED-011 - QVAC adapter request body and headers are inspected safely

Package: `internal/adapter/external/qvac`

- Use `httptest` to capture QVAC `/chat/completions` requests.
- Assert request JSON body contains no sentinels.
- Assert QVAC Authorization header, if configured for local QVAC itself, is never echoed into errors, Evidence or logs.
- Assert outbound model messages never include provider headers, cookies, DSNs or raw external error bodies.

### RED-012 - QVAC tools remain read-only and allowlisted

Package: `internal/adapter/external/qvac`

- Assert tool registry exposes exactly expected read-only tools.
- Unknown/mutating tools (`write_file`, `create_branch`, `commit`, `push`, `create_issue`, `create_pr`, `shell`) return stable errors and produce no side effect.
- Assert tool schemas disallow model-selected owner/repo/account and bind to configured RepositoryScope.

### RED-013 - Tool outputs are bounded, valid UTF-8 and sanitized

Package: `internal/adapter/external/qvac`

- Fake RepositoryInspector returns oversized content, invalid UTF-8, provider errors with secrets, code containing tokens and multiline payloads.
- Assert tool message content is valid UTF-8, below max size, redacted, framed as untrusted data and never leaks sentinels.
- Assert persisted discovered Evidence uses the sanitized payload, not raw tool output.

## AKR-61 - Absence of secrets in every sink

### RED-020 - Central sanitizer covers required patterns

Package: `internal/service/evidencesafety` or new content safety package

Table-driven cases for:

- Bearer Authorization
- Basic Authorization
- Cookie
- Set-Cookie
- JSON strings
- quoted assignments
- unquoted assignments
- PostgreSQL DSN
- MySQL DSN
- Redis DSN
- Mongo DSN
- AMQP/RabbitMQ DSN
- GitHub PAT classic
- GitHub fine-grained token
- GitHub App private key PEM
- GitHub JWT
- GitHub installation token
- AWS-style access keys
- `*_TOKEN`
- `*_PASSWORD`
- `*_SECRET`
- `*_API_KEY`
- values with spaces
- multiline content
- invalid UTF-8
- truncation after redaction

Assertions:

- No original secret fragment remains.
- Output is valid UTF-8.
- Truncation marker is visible when truncation occurs.
- Return metadata correctly reports whether redaction occurred.

### RED-021 - `output_redacted` is truthful

Packages: `internal/core/domain`, `internal/usecase/remediation`, `internal/adapter/db/postgres/repository/validationresult`

- Validation output with no secret yields `OutputRedacted=false` after sanitization.
- Validation output with a secret yields `OutputRedacted=true`.
- DB persists and returns the exact boolean.
- Expected initial state: fails because domain/DB force `true`.

### RED-022 - Evidence persistence rejects/normalizes unsafe content

Packages: `internal/core/domain`, `internal/usecase/investigation`, `internal/adapter/db/postgres/repository/evidence`

- Attempt to persist Evidence containing each sentinel format.
- Assert persisted Content/Summary/FilePath/Patch are sanitized and metadata indicates redaction.
- Assert Evidence size limits still apply after redaction.

### RED-023 - GitHub Issue content has no secrets

Package: `internal/service/issuecontent`

- Build Issue content from Incident, Investigation and Evidence containing all sentinel types.
- Assert title/body contain no sentinels, no raw provider headers and no raw DSNs.
- Assert Investigation marker remains intact.
- Assert model-generated conclusions are labeled as conclusions, not verified facts.

### RED-024 - Pull Request content has no secrets

Package: future PR content service or GitHub PR adapter tests

- Build PR title/body from Remediation, changes, validation outputs and Issue reference with sentinel values.
- Assert no sentinels, no raw stdout/stderr, no absolute workspace path.
- Assert PR marker/correlation IDs remain intact.

### RED-025 - Logs, errors, timeline and REST projections are sanitized

Packages: relevant REST handlers/mappers/usecases

- Inject errors/causes/timeline events with sentinels.
- Assert REST responses, operation user messages, timeline DTOs and any public error details do not leak.
- Assert internal logs do not log full external request/response bodies.

## AKR-65 - Deterministic commit correlation

### RED-030 - Commit selection is deterministic

Package: new `internal/service/commitcorrelation`

- Given Incident first/last seen timestamps, branch, relevant files and a list of commit summaries/details, select a stable bounded set.
- Sort deterministically by time proximity, file/path hint score and SHA tie-breaker.
- Deduplicate repeated commits.
- Enforce max count and max field sizes.

### RED-031 - Commit correlation persists safe Evidence

Package: `internal/usecase/investigation`

- Fake GitHub commit reader returns commits with secret-bearing messages/authors/patches.
- Assert persisted Evidence has type `commit`, safe SHA/timestamp/author/message/files, sanitized content, and no sentinel.
- Assert Evidence says "potentially related" or equivalent, never "confirmed cause".

### RED-032 - GitHub unavailable degrades safely

Package: `internal/usecase/investigation`

- Fake commit reader returns auth/API/rate-limit/unavailable errors.
- Assert Investigation still starts/runs with a safe audit note or absent commit Evidence according to design.
- Assert QVAC receives no raw provider error body.
- Assert Operation is not failed solely due to optional commit correlation failure unless policy marks the dependency mandatory.

### RED-033 - Commit data passed to QVAC is bounded and safe

Package: `internal/adapter/external/qvac`

- Run QVAC with commit Evidence containing large patch/message/author data.
- Assert prompt and tool payload limits hold and all content is redacted.

## AKR-62 - Idempotency and reconciliation

### RED-040 - Same Investigation yields same GitHubIssueReference

Package: `internal/usecase/investigation`

- Run completed Investigation with existing IssueReference.
- Assert publisher is not called and existing reference is reused.

### RED-041 - GitHub Issue external success/local failure reconciles

Packages: `internal/usecase/investigation`, `internal/adapter/external/github`

- First attempt: GitHub creates Issue, local `IssueReference.Create` fails.
- Retry: GitHub adapter finds existing Issue by Investigation marker and repository.
- Assert no second Issue is created and local reference is persisted.
- Assert repository mismatch or missing marker is not reused.

### RED-042 - One Remediation per eligible Investigation

Package: `internal/usecase/remediation`

- Concurrent calls start Remediation for the same completed `fixable` Investigation.
- Assert one Remediation row and all callers receive/reconcile the same logical Remediation.
- Assert non-fixable or missing IssueReference rejects with stable error and creates nothing.

### RED-043 - One logical branch per Remediation

Packages: `internal/usecase/remediation`, `internal/adapter/external/git`, GitHub branch adapter if remote branch is used

- Concurrent/retried branch stage uses deterministic branch name.
- Existing local/remote branch with matching marker is reconciled.
- Existing branch without marker or wrong base is terminal conflict.
- No DB transaction is open during Git/GitHub branch calls.

### RED-044 - One successful commit per Remediation

Package: `internal/usecase/remediation`

- Retry after successful commit but failed local persistence.
- Reconcile by commit marker and head branch.
- Assert no second commit is created for the same Remediation.
- Wrong repository/branch commit is not reused.

### RED-045 - One Pull Request per Remediation

Packages: `internal/usecase/remediation`, `internal/adapter/external/github`

- Retry after PR creation but failed local persistence.
- Reconcile by repository, base branch, head branch and PR marker.
- Assert no duplicate PR.
- Assert closed/wrong-base/wrong-head/wrong-repository PR is not reused unless policy explicitly allows.

### RED-046 - Idempotency keys and operations

Packages: REST handlers/usecases/repositories

- `StartIncidentRemediation` and `CreateRemediationPullRequest` require `Idempotency-Key`.
- Same key returns same Operation.
- Different key for same already-in-progress resource reconciles to existing Operation/Remediation without duplicate external effects.

## AKR-63 - Failures and recovery

### RED-050 - Failure matrix classification

Package: new failure classification service

Table-driven cases:

- Dokploy unavailable/auth/API
- GitHub auth/API/rate limit
- QVAC unavailable/timeout/invalid output/tool limit
- checkout/clone/workspace
- tool loop
- change generation
- validation failed
- validation runner failed
- branch
- commit
- push
- Issue creation
- Pull Request creation
- local persistence before external effect
- local persistence after external effect
- context cancellation
- shutdown

Assertions:

- stable owner error
- retryable vs terminal
- safe user message
- audit category
- expected state transition
- reconciliation strategy

### RED-051 - Validation failure fails Remediation and stops before commit/push/PR

Package: `internal/usecase/remediation`

- Validation runner returns failed validation result.
- Assert all available ValidationResults are persisted.
- Assert Remediation status becomes `failed` with safe user message.
- Assert commit, push and PR ports are never called.

### RED-052 - Runner/tool infrastructure failure is distinguishable from validation failure

Package: `internal/usecase/remediation`

- Runner cannot start process or workspace unreadable.
- Assert failure classification is not treated as test failure.
- Assert status and operation state follow failure matrix.

### RED-053 - No indefinite running/in_progress states after recovery

Package: recovery usecases/services

- Seed queued/running operations and in-progress remediations older than configured timeout.
- On startup recovery:
  - safe queued work is requeued
  - interrupted running work is failed or reconciled according to checkpoint
  - confirmed external artifacts are preserved
- Assert Evidence, ValidationResults, IssueReference and PRReference are not deleted.

### RED-054 - Cancellation and shutdown are safe

Packages: usecases/adapters

- Cancel context during QVAC, GitHub, Git, validation and persistence boundaries.
- Assert external calls receive context cancellation.
- Assert DB transaction is not left open.
- Assert final state is retryable/failed according to matrix.

## Remediation happy path and STOP

### RED-060 - Full fake-backed happy path

Package: integration-style usecase test

With fakes for Dokploy/GitHub/QVAC/Git/process:

```text
LogEvents
-> grouped Incident
-> safe Evidence
-> commit correlation Evidence
-> QVAC fixable Investigation
-> unique Issue
-> unique Remediation
-> workspace
-> branch
-> regression test red before fix
-> minimal change
-> validations passed
-> commit
-> push
-> PR
-> STOP
```

Assertions:

- No sentinel secret in any captured request/prompt/tool/result/Evidence/Issue/PR/REST projection.
- IDs and markers correlate Incident, Investigation, Issue, Remediation, commit and PR.
- No merge/deploy/rollback port exists or is called.

### RED-061 - Requires-human path stops after Issue

Package: `internal/usecase/investigation` / orchestration

- QVAC returns `requires_human`.
- Assert Issue is created.
- Assert no Remediation/workspace/branch/change/commit/push/PR calls occur.

### RED-062 - STOP after PR

Package: `internal/usecase/remediation`

- After PR creation, assert status is `pull_request_created`.
- Assert no further mutating provider capability is invoked.
- Assert no merge/deploy/rollback dependencies are wired.

## AKR-64 - Reproducible demo E2E

### RED-070 - Demo fixture reproduces repeated error

Location: demo fixture package/script tests

- Fixture app has a deterministic bug.
- Test or script can trigger the same error multiple times.
- Assert generated logs produce same fingerprint and one Incident.
- No secrets are required in committed config.

### RED-071 - Demo runbook has safe closed arguments

Location: demo script tests or docs validation

- Scripts reject arbitrary shell fragments.
- Required args are enumerated and validated.
- Cleanup is confined to Akritas-managed demo/workspace root.
- No tokens, cookies, DSNs or credentials are committed.

### RED-072 - Demo retry reuses Issue and PR

Integration test with fakes/local servers:

- Run demo logical scenario twice with same deterministic IDs/markers.
- Assert same Issue and same PR are reused/reconciled.
- Assert no duplicate branch/commit/PR for the same logical execution.

### RED-073 - Opt-in real E2E test documents prerequisites

Build tag: `e2e` or explicit opt-in equivalent.

- Skips unless real Dokploy/GitHub/QVAC prerequisites are configured through environment.
- Asserts all acceptance stages and STOP after PR.
- Prints/saves only sanitized artifact references.

## Test levels required

- Domain unit tests:
  - state transitions
  - truthful redaction metadata
  - Remediation invariants
  - PullRequestReference/commit/reference invariants
- Service unit tests:
  - sanitizer/redactor
  - safe QVAC context
  - commit correlation
  - failure classification
  - reconciliation decisions
- Usecase tests with manual fakes:
  - Investigation -> Issue
  - Investigation -> commit Evidence -> QVAC
  - Remediation orchestration
  - idempotency
  - concurrency
  - external success/local persistence failure
  - recovery
- Adapter tests with `httptest`:
  - GitHub Issues/branches/commits/PRs/reconciliation
  - QVAC request/tool loop/error handling
  - Dokploy failure normalization if touched
- Adapter process/Git tests:
  - fixed argv only
  - workspace confinement
  - cancellation/timeouts
- PostgreSQL tests:
  - migrations
  - constraints/indexes
  - repository atomic operations
  - rollback
- Integration tests:
  - fake-backed full flow
  - opt-in external E2E/demo

## Validation gates after implementation

Run and report:

- `gofmt`
- `go build ./...`
- `go test -p 1 ./...`
- `go test -race -p 1 ./...`
- `go test -tags=integration -p 1 ./...`
- `go vet ./...`
- `.harness/kernel/scripts/check-backend-architecture.sh`
- `.harness/kernel/scripts/check-security.sh`
- `.harness/kernel/scripts/check-openapi.sh`
- migration tests
- `git diff --check`

If a baseline or preexisting test fails, classify it honestly. Do not hide failing tests.

## Human approval gate

Implementation is blocked until this TDD plan is explicitly approved by the user.
