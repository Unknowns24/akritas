# AKR-H6-BACKEND-HARDENING-DEMO - H6 backend hardening and reproducible demo

## Estado

implemented_with_limitations

## Tipo de tarea

backend-service-feature with backend-api/OpenAPI verification

## Modo de proyecto

existing_project

## Contexto

Task: `AKR-H6-BACKEND-HARDENING-DEMO`

Proyecto Linear: `Akritas MVP`

Milestone: `H6 - Hardening + Demo`

Included backend issues:

- `AKR-60 / PB-051` - Guarantee QVAC never receives credentials.
- `AKR-61 / PB-052` - Guarantee absence of secrets in logs, Evidence, GitHub Issue and PR.
- `AKR-62 / PB-053` - Avoid duplicate Issues and PRs under retries.
- `AKR-63 / PB-054` - Handle integration and QVAC failures explicitly.
- `AKR-64 / PB-055` - Execute a reproducible demo E2E scenario.
- `AKR-65 / PB-056` - Correlate Incident with recent commits.

Explicitly excluded:

- `AKR-66` and `AKR-67`: frontend/UI.
- `AKR-68` and `AKR-69`: already complete.
- Stretch/Post-MVP.
- Automatic merge, automatic deploy, production rollback, and mutations after Pull Request creation.

Harness resolution:

- Primary profile: `backend_service` because the core scope is workers, QVAC, GitHub, persistence, retries, recovery, filesystem/Git/process and E2E orchestration.
- Workflow: `.harness/kernel/workflows/backend-service-feature.yaml`.
- Additional API/OpenAPI policy loaded because the current OpenAPI exposes remediation, PR and QVAC lifecycle surfaces that must be verified against real handlers/routes/bootstrap.
- Required policies loaded: backend project structure, wiring, modularity/SRP, domain errors, migrations, external adapters, background processes, Uker, testing, security, architecture decisions, plus OpenAPI policy.
- Architecture/memory sources loaded: harness memory, backend architecture, incident lifecycle, demo story, MVP/product backlog, ADRs 001/002/003/004/014/015 and RFC-002 as future-only context.

Baseline source of truth:

- `git fetch origin main` completed successfully after sandbox escalation.
- Local `HEAD` equals `origin/main`: `c325c8a5bd3d7a7508a049b56006d93cef281f83`.

## Objetivo

Leave the backend demonstrable, retry-safe and security-hardened for this flow:

```text
real repeated Dokploy error
-> LogEvents
-> grouped Incident
-> safe Evidence
-> potentially relevant recent commits
-> local QVAC Investigation
-> unique GitHub Issue
-> unique fixable Remediation
-> isolated workspace
-> dedicated branch
-> regression test + proposed change
-> validations
-> commit
-> push
-> unique Pull Request
-> STOP
```

The system must never merge, deploy, rollback or modify production after creating the Pull Request.

## Requerimiento funcional

### Baseline blockers to address before H6 can be proven

1. `internal/adapter/db/postgres/errors/catalog.go` does not compile. `ErrValidationResultPersistence` is not closed before `ErrGitHubIssueReferencePersistence`, causing `go build ./...`, `go test -p 1 ./...` and `go vet ./...` to fail.
2. Migration IDs/order are conflict-prone and currently encode H5 migrations before H4 migrations on the same date/order slots:
   - `20260823_06_add_remediations`
   - `20260823_07_add_validation_results`
   - `20260823_06_add_github_issue_references`
   - `20260823_07_enforce_issue_reference_investigation_incident`
   - `20260823_08_add_dokploy_compose_sources`
3. The existing migration registry test hardcodes the conflicted order instead of detecting same-day slot conflicts or dependency ordering.
4. OpenAPI declares QVAC configuration/status, Remediation and Pull Request endpoints, but the actual REST router/handler/DTO wiring does not expose QVAC, Remediation or Pull Request routes.
5. `portsin.UseCases` does not expose Remediation.
6. `internal/bootstrap/integrations/module.go` does not construct remediation repositories, Git workspace adapter, validation runner, remediation usecase or handlers.
7. H5 Remediation usecase intentionally implements only branch creation and validations. It receives `RemediationID`, `IncidentID`, `WorkspacePath` and `BaseBranch` from the caller rather than starting from persisted `Investigation -> IssueReference`.
8. Current Remediation persistence is minimal and does not cover Investigation, IssueReference, changes, validation lifecycle decision, commit, push, PR reference, checkpoints or recovery.
9. Current validation failure persistence does not automatically move Remediation to `failed`.
10. Existing `ValidationResult.OutputRedacted` is forced to `true` by the domain and DB check constraint, which cannot honestly distinguish "sanitization applied" from "no redaction was needed".

### H5 prerequisites indispensable for H6

- Create/reconcile a Remediation from a completed `fixable` Investigation with durable IssueReference.
- Resolve repository/account/base branch/workspace from persisted Project/Investigation/Issue state, not caller-supplied paths.
- Persist full Remediation state transitions and external artifacts: branch, changes, validation results, commit, push and PR reference.
- Add narrow output ports/adapters for change generation, workspace isolation, commit, push and PR creation/reconciliation.
- Wire Remediation through bootstrap, input ports, REST handlers/routes and OpenAPI only where observation or manual commands are required.
- Ensure validation failure persists all available results, marks Remediation as `failed`, and prevents commit/push/PR.

### H6 hardening scope

- AKR-60: explicit safe QVAC data boundary; credentials resolved only inside owner adapters at the last possible moment; QVAC tools remain read-only and allowlisted; tool inputs/outputs are bounded, valid UTF-8 and sanitized.
- AKR-61: one content safety policy applied consistently to logs, causes, Evidence, command excerpts, validation stdout/stderr, QVAC payloads/results, GitHub Issue title/body, Pull Request title/body, change summaries, timeline and REST projections.
- AKR-62: idempotency and reconciliation for Investigation IssueReference, Remediation, branch, commit and Pull Request, including concurrent retries and external-success/local-persistence-failure scenarios.
- AKR-63: explicit failure matrix, stable layer-owned errors, retryable vs terminal classification, state transitions, safe user messages, persisted audit information and startup recovery for orphaned operations.
- AKR-64: controlled reproducible demo fixture/runbook/scripts with no committed secrets and retry-safe external artifacts.
- AKR-65: deterministic bounded recent commit correlation integrated into the real Investigation pipeline as safe Evidence, without treating temporal correlation as causality and without blocking Investigation when GitHub is unavailable.

## Criterios de aceptación

- Baseline build/test/vet blockers are fixed and preexisting failures are not hidden.
- Migration registry has safe, explicit handling of the H4/H5 ordering conflict, with tests for order, uniqueness, upgrade and rollback.
- QVAC request bodies, prompts, tool args, tool results, Evidence and persisted Investigation data never contain GitHub, GitHub App, Dokploy, PostgreSQL, session, bootstrap, TOTP or master-key secrets.
- Every externally visible text sink uses the same sanitizer/redaction contract and has negative sentinel-secret tests.
- Retrying the same completed Investigation yields the same GitHub IssueReference.
- Retrying the same eligible Investigation yields one Remediation, one logical branch, one successful commit and one Pull Request.
- GitHub external success followed by local persistence failure is reconciled by deterministic markers/correlation IDs without duplicate Issues or PRs.
- No PostgreSQL transaction is held while calling GitHub, Dokploy, QVAC, Git, filesystem or validation processes.
- Failure states cannot remain ambiguous forever; startup recovery handles orphaned queued/running operations where applicable.
- Validation failure marks Remediation `failed` and prevents commit/push/PR.
- Commit correlation persists bounded safe commit Evidence and degrades safely when GitHub is unavailable.
- The demo scenario can be re-run without creating duplicate Issues or PRs for the same logical execution.
- The automated backend flow stops immediately after creating the Pull Request.

## Restricciones tecnicas

- TDD approval gate is required before code implementation.
- Maintain hexagonal architecture and SRP. Do not create a giant H6 manager/service.
- Use narrow input/output ports by capability.
- Domain must not import external clients, persistence behavior, HTTP, GORM or process execution.
- Adapters own GitHub/Dokploy/QVAC/Git/process/filesystem details and credential resolution.
- No generic shell execution. Validation and Git operations must use closed enums/argv.
- Workspaces must be confined under an Akritas-managed root; cleanup only inside that root.
- Do not keep DB transactions open while calling external systems or slow processes.
- Do not rename already-applied migration IDs silently. Investigate shared-environment exposure first and use an explicit safe migration strategy.
- OpenAPI must be updated before REST implementation if the contract changes. Do not invent endpoints if existing contract is sufficient.
- REST/API projections must never expose secrets, absolute internal paths or unsanitized stdout/stderr.

## Archivos o zonas probablemente afectadas

- `.harness/tasks/AKR-H6-BACKEND-HARDENING-DEMO/`
- `.harness/tasks/index.md`
- `.harness/memory/project-summary.md`
- `.harness/memory/decisions.md`
- `docs/openapi.yaml` if lifecycle observability/commands need contract corrections.
- `docs/errors/aaa-map.md`
- `docs/demo.md` and/or backend demo runbook/fixture docs.
- `internal/core/domain/*remediation*`, `*evidence*`, `*operation*`, `*timeline*`, `*pull_request*`, `*errors*`
- `internal/core/ports/in/*remediation*`, `use_cases.go`
- `internal/core/ports/out/*github*`, `*qvac*`, `*workspace*`, `*validation*`, `*sanitization*`, `*commit*`, `*pull_request*`
- `internal/service/evidencesafety`, plus likely focused new services for safe QVAC context, artifact sanitization, commit correlation, reconciliation and failure classification.
- `internal/usecase/investigation`
- `internal/usecase/remediation`
- `internal/adapter/external/qvac`
- `internal/adapter/external/github`
- `internal/adapter/external/git`
- `internal/adapter/external/validationrunner`
- `internal/adapter/db/postgres/migrations`
- `internal/adapter/db/postgres/repository/remediation`
- `internal/adapter/db/postgres/repository/validationresult`
- new/extended PostgreSQL repositories for PR/reference/checkpoint state as needed.
- `internal/adapter/rest/dto/<feature>`, `handler/<feature>`, `mapper`, `router`
- `internal/bootstrap/integrations/module.go`
- demo fixture/scripts under an appropriate docs/test/demo location.

## Fuera de alcance

- Frontend/UI implementation (`AKR-66`, `AKR-67`).
- Auth recovery/session hardening already completed (`AKR-68`, `AKR-69`).
- Stretch/Post-MVP.
- Automatic merge, deploy, rollback or any production mutation after PR creation.
- Consuming admin-only endpoints from public/user-final app surfaces.
- Secret rotation automation beyond ensuring no leakage.

## Instruccion para el harness

This first execution stops at the Human/TDD approval gate. Do not implement code until `tdd-test-plan.md` is explicitly approved by the user.
