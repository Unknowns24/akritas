# Implementation Brief

## Task

`AKR-H6-BACKEND-HARDENING-DEMO` - final backend integration for H6 Hardening + Demo, covering `AKR-60` through `AKR-65` only.

## Current project context

Profile resolution:

- Primary profile: `backend_service`.
- Workflow: `backend-service-feature`.
- OpenAPI/API policies also loaded because `docs/openapi.yaml` declares observable QVAC, Remediation and Pull Request lifecycle surfaces.

Remote baseline:

- `origin/main` fetched and verified.
- Local `HEAD` equals `origin/main` at `c325c8a5bd3d7a7508a049b56006d93cef281f83`.

Preflight from `backend/`:

| Check | Result |
| --- | --- |
| `go build ./...` | failed: syntax error in `internal/adapter/db/postgres/errors/catalog.go:40` |
| `go test -p 1 ./...` | failed: same compile break; independent packages passed where build graph allowed |
| `go vet ./...` | failed: same parse error |
| `check-backend-architecture.sh` | passed via Git Bash with Unix PATH |
| `check-security.sh` | passed |
| `check-openapi.sh` | passed via Git Bash with Python 3.12.5 shim |
| `git diff --check` | passed |

Confirmed baseline blockers:

- `catalog.go` has an unterminated `ErrValidationResultPersistence` declaration and nests `ErrGitHubIssueReferencePersistence` inside it.
- Migration registry and tests encode the conflicted H5/H4 ordering:
  - H5 `20260823_06_add_remediations`
  - H5 `20260823_07_add_validation_results`
  - H4 `20260823_06_add_github_issue_references`
  - H4 `20260823_07_enforce_issue_reference_investigation_incident`
  - then `20260823_08_add_dokploy_compose_sources`
- Existing migration tests verify the current order but do not reject same-day slot collisions or dependency inversion.
- `internal/core/ports/in/remediation.go` explicitly excludes trigger, change/test generation and failure-decision logic.
- `internal/usecase/remediation` implements branch creation and validation execution only.
- `internal/core/ports/in/use_cases.go` does not expose Remediation.
- `internal/bootstrap/integrations/module.go` does not wire remediation repositories, Git adapter, validation runner, remediation usecase, Remediation handler or PR handler.
- The REST router/handler/DTO tree has no `remediation`, `pull_request` or `qvac` handler package, while OpenAPI declares those paths.
- Current `RemediationStore` is Create+Get only.
- Current `remediations` table lacks Investigation, IssueReference, changes, commit, push, PR and checkpoint state.
- Current `validation_results.output_redacted` is constrained to `true`, which is not truthful enough for AKR-61.

Important existing strengths:

- QVAC tools are currently read-only and allowlisted for repository inspection.
- QVAC prompt/tool payloads already have size limits, untrusted-data framing and redaction in several paths.
- GitHub credentials are resolved in the GitHub adapter through Credential Store.
- H4 Issue publication uses short transactions around GitHub calls and includes an Investigation HTML marker in the Issue body.
- Git and validation adapters use fixed argv through `exec.CommandContext`, not shell strings.

## Proposed approach

Phase 0 - Baseline repair and regression tests:

- Add failing tests for the DB error catalog compile/registration shape where possible, then fix `catalog.go`.
- Add migration registry tests that reject same-day slot collisions, prove dependency ordering, and exercise upgrade/rollback on real PostgreSQL.
- Investigate whether H5 migration IDs may have run in any shared environment before deciding between preserving IDs with compensating forward migrations or explicitly renumbering only when proven safe. Do not silently rename deployed IDs.
- Add route/contract correspondence tests for OpenAPI paths vs registered routes and handler wiring, so missing QVAC/Remediation/PR lifecycle surfaces are explicit.

Phase 1 - Safe content boundary:

- Introduce a focused content safety service/value object that returns `(sanitized, redacted)` or equivalent metadata.
- Apply it consistently to Evidence, QVAC prompts, QVAC tool args/results, validation stdout/stderr, logs/errors intended for external/user surfaces, Issue content, PR content, change summaries, timeline and REST DTOs.
- Keep redaction as defense in depth; do not authorize sending known credentials to QVAC.
- Ensure credentials remain IDs/references in usecases and are resolved only inside owner adapters.

Phase 2 - Commit correlation:

- Add deterministic commit-correlation service that selects bounded recent commits from the configured/default branch around Incident time.
- Use timestamp window, branch, dedupe, file/location hints when available, strict count/size limits and sanitizer.
- Persist safe commit summaries/files as Evidence and pass only bounded safe Evidence to QVAC.
- Degrade safely if GitHub is unavailable and do not block Investigation creation/completion solely because commit correlation failed.

Phase 3 - Remediation orchestration and idempotency:

- Add narrow input ports/usecases for starting/recovering Remediation from persisted `Investigation -> IssueReference`.
- Add output ports for workspace provisioning, branch reconciliation, change generation, validation policy, commit, push, PR creation and PR reconciliation.
- Use deterministic keys and auditable external markers/correlation IDs:
  - Investigation marker for Issue.
  - Remediation/branch marker.
  - Commit marker.
  - PR marker referencing Remediation, Investigation and Issue.
- Enforce local uniqueness through PostgreSQL constraints/indexes and repository operations that are atomic locally.
- Reconcile external success/local failure by querying GitHub using repository, base branch, head branch and markers. Never reuse a PR from another repository or branch.
- Keep all GitHub/QVAC/Git/filesystem/process calls outside DB transactions.

Phase 4 - Failure matrix and recovery:

- Define stable layer-owned errors for Dokploy, GitHub auth/API/rate limit, QVAC unavailable/timeout/invalid output, workspace, tool loop, change generation, validation, branch, commit, push, Issue, PR and local persistence boundaries.
- Classify retryable vs terminal.
- Persist safe user messages and audit detail.
- Move Remediation out of ambiguous states on terminal failures.
- Add startup recovery for queued/running/orphaned work where product behavior is deterministic.
- Validation failure must persist all available results, mark Remediation `failed`, and prevent commit/push/PR.

Phase 5 - REST/OpenAPI observation and demo:

- Prefer existing OpenAPI paths if sufficient. If the contract is wrong or unverifiable, update `docs/openapi.yaml` first and bump version coherently.
- Implement missing DTOs/mappers/handlers/routes only for observable lifecycle and allowed manual commands.
- Add a controlled demo fixture/runbook/scripts with no secrets and closed arguments.
- Ensure repeat demo execution reuses the same logical Issue/PR instead of creating duplicates.
- Document real external dependencies vs fake/local integration-test mode.

## Architecture impact

Expected architecture changes are additive but cross-cutting:

- New focused services, not one large H6 manager:
  - safe QVAC context
  - artifact sanitization
  - commit correlation
  - retry/reconciliation
  - failure classification
  - remediation orchestration
  - demo fixture/runbook
- Usecases orchestrate application behavior through ports.
- External adapters retain provider/process/filesystem details.
- Repositories remain separated by aggregate/capability.
- Bootstrap composes concrete dependencies in the configured order.

Architecture risks:

- Remediation can easily become a monolith. Keep operations split by capability and file.
- Reconciliation can accidentally couple usecases to GitHub API details. Keep provider search/query mechanics in the GitHub adapter behind capability ports.
- Safe-content policy can become a generic dumping ground. Keep it cohesive around sanitization/redaction contracts.

## API/OpenAPI impact

`docs/openapi.yaml` currently passes the OpenAPI gate and declares:

- QVAC configuration/status paths.
- Remediation paths.
- Pull Request paths.
- Idempotency-Key requirements for remediation commands.

However, the real router/handler tree does not currently expose all those declared paths. The implementation must first add contract correspondence tests, then either:

- implement the already-declared contract exactly; or
- update `docs/openapi.yaml` before code if the declared lifecycle surface is not the desired one.

No API response may expose credentials, absolute workspace paths, full raw stdout/stderr, unsanitized tool payloads or provider error bodies.

## Data/persistence impact

Likely persistence additions:

- Safe migration strategy for H4/H5 ordering conflict.
- Full remediation lifecycle columns/tables:
  - investigation_id
  - github_issue_reference_id or equivalent durable link
  - repository/base_branch/head_branch
  - workspace logical ID/root-relative path only, not absolute internal paths in REST
  - branch/base commit/head commit
  - changes summary and bounded sanitized changes metadata
  - commit reference
  - push state
  - pull request reference
  - failure classification/user message
  - checkpoints/recovery metadata
- Constraints/indexes for:
  - one Remediation per eligible Investigation
  - one logical branch per Remediation
  - one successful commit per Remediation
  - one Pull Request per Remediation
  - repository/base/head uniqueness for PR reconciliation
  - unique external markers where persisted
- `validation_results.output_redacted` semantics must become truthful. Prefer storing both sanitized output and `output_redacted` based on whether redaction actually occurred, while preserving bounded UTF-8 output.

Migration strategy:

- Preserve H4 migrations already integrated.
- Determine whether H5 IDs `20260823_06_add_remediations` and `20260823_07_add_validation_results` have been applied outside ephemeral local DBs.
- If potentially applied, do not rename those IDs. Add forward migrations that repair/extend schema and registry ordering tests that document the historical conflict.
- If proven not applied anywhere shared, propose explicit renumbering with human approval before editing IDs.
- Cover registry order, ID uniqueness, schema upgrade and rollback with tests.

## Error handling impact

Add or extend layer-owned errors for:

- safe-content policy failures
- commit correlation degraded/unavailable states
- QVAC timeout/unavailable/invalid output/tool-limit/tool-output failures
- GitHub auth/API/rate-limit/reconciliation failures
- branch/commit/push/PR failures
- workspace checkout/provisioning/confinement failures
- validation execution vs validation failed result
- local persistence before and after external side effects
- orphan recovery failure

Every failure must define:

- stable error owner/code
- retryable vs terminal
- state transition
- safe user message
- persisted audit data
- reconciliation strategy
- cancellation/shutdown behavior

## Test strategy

The approved TDD plan must lead implementation. Tests should use:

- domain unit tests for redaction metadata, state transitions and invariants
- service unit tests for safe context, sanitization, commit selection, failure classification and reconciliation decisions
- usecase tests with manual fakes for orchestration, idempotency, concurrency and external-success/local-failure
- adapter tests with `httptest` for GitHub/QVAC/Dokploy behavior and sanitized provider errors
- real PostgreSQL migration/repository tests following existing `dbtest` conventions
- integration tests for the full fake-backed happy path
- opt-in E2E test/runbook for real Dokploy/GitHub/QVAC demo prerequisites

See `tdd-test-plan.md` for the RED list.

## Risks

- Baseline does not compile, so H6 cannot be proven until baseline blocker is fixed.
- H5 is incomplete relative to full Remediation + PR. H6 must not declare success on top of only branch+validation.
- Migration IDs may have been applied in a shared environment; silent renumbering could break deployments.
- OpenAPI and runtime route mismatch can mislead frontend/demo work.
- Sanitization gaps are high risk because H6 touches many sinks.
- Idempotency around external effects cannot rely on DB transactions alone.
- QVAC redaction must be defense in depth, not a substitute for safe data selection.
- Commit correlation can be over-interpreted as causality; wording and data model must keep it as "potentially related".

## Files likely to change

See `task.md` for the full likely areas. High-risk files include:

- `internal/adapter/db/postgres/errors/catalog.go`
- `internal/adapter/db/postgres/migrations/*`
- `internal/core/domain/remediation.go`
- `internal/core/domain/validation_result.go`
- `internal/core/ports/in/use_cases.go`
- `internal/core/ports/in/remediation.go`
- `internal/core/ports/out/*`
- `internal/service/evidencesafety/*`
- `internal/adapter/external/qvac/*`
- `internal/adapter/external/github/*`
- `internal/usecase/investigation/*`
- `internal/usecase/remediation/*`
- `internal/bootstrap/integrations/module.go`
- `internal/adapter/rest/**`
- `docs/openapi.yaml`
- `docs/errors/aaa-map.md`
- demo/runbook/fixture files

## Human gate

Stop here until `tdd-test-plan.md` is approved explicitly.
