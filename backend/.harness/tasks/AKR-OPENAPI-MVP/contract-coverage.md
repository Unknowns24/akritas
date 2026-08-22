# Contract Coverage — AKR-OPENAPI-MVP

## Product backlog P0

| ID | Contract coverage |
| --- | --- |
| PB-001 | GitHub account collection/item CRUD and `GitHubAccount` |
| PB-002 | PAT write-only request plus `credential_configured` safe projection |
| PB-003 | GitHub `connection-test` and normalized `ConnectionTestResult` |
| PB-004 | Cursor-paginated GitHub repositories |
| PB-005 | Dokploy server collection/item CRUD and `DokployServer` |
| PB-006 | Dokploy credential write-only requests and safe projection |
| PB-007 | Dokploy `connection-test` |
| PB-008 | Cursor-paginated Dokploy applications |
| PB-009 | Project collection/item create, list, get and update |
| PB-010 | Project GitHub account/repository selection fields |
| PB-011 | Project Dokploy server/application selection fields |
| PB-012 | Monitoring configuration get/replace and schema |
| PB-013 | `MonitoringConfiguration.enabled` and `MonitoringStatus` |
| PB-014 | Project/Application references and monitoring status needed by log acquisition |
| PB-015 | Operational monitoring remains backend-owned; no cursor internals exposed to UI |
| PB-016 | `LogEvent` represents assembled multiline events rather than raw line CRUD |
| PB-017 | `ignored_patterns` schema and documented precedence |
| PB-018 | Read-only `BuiltInDetectionRule` catalog |
| PB-019 | `error_patterns` schema |
| PB-020 | Stable normalized fingerprint exposed by `LogEvent`/`Incident` |
| PB-021 | Fingerprint fields on LogEvent and Incident |
| PB-022 | Sanitized `context_before` / `context_after` records |
| PB-023 | Cursor-paginated Incident LogEvents |
| PB-024 | Incident occurrence count and first/last seen projections |
| PB-025 | Incident list/detail/log-events/timeline operations |
| PB-026 | Investigation collection, command and detail operations |
| PB-027 | Typed, sanitized Evidence collection |
| PB-028 | QVAC configuration/status and local/private endpoint constraint |
| PB-029 | Optional `tool_used` timeline events without provider payload leakage |
| PB-030 | Evidence types for code locations and sanitized content |
| PB-031 | Evidence types for commits/diffs and relevant commit fields |
| PB-032 | Structured `Investigation` schema |
| PB-033 | `RootCauseStatus` enum |
| PB-034 | `ResolutionStatus` enum |
| PB-035 | Investigation confidence, root cause and evidence count |
| PB-036 | Issue reference required by completed investigation lifecycle documentation |
| PB-037 | `GitHubIssueReference` linked from Incident and PullRequest |
| PB-038 | Investigation/Issue safe evidence and classification projections |
| PB-039 | Incident detail and Investigation endpoints expose Issue result |
| PB-040 | Manual/automatic remediation resource for fixable incidents |
| PB-041 | `Remediation.branch_name` |
| PB-042 | Sanitized `CodeChange` collection and patch |
| PB-043 | Validation type `test` and code-change projection |
| PB-044 | Validation result collection and statuses |
| PB-045 | Persistable `ValidationResult` projection |
| PB-046 | `failed` validation/remediation and PR command conflict semantics |
| PB-047 | Validated remediation state before PR; commit remains internal orchestration |
| PB-048 | Idempotent async Pull Request command and PR resources |
| PB-049 | PullRequest cross-references Incident, Remediation and Issue |
| PB-050 | Remediation detail exposes branch, changes, validations and PR |
| PB-051 | QVAC DTOs exclude all integration/auth credentials |
| PB-052 | Evidence, logs, patches and output explicitly marked sanitized/redacted |
| PB-053 | `Idempotency-Key` plus async `Operation` for retryable mutations |
| PB-054 | Normalized integration statuses, operation failure and terminal outcomes |
| PB-055 | Overview plus full Project → Incident → Investigation → Issue → Remediation → PR contract |

## Authentication additions

| ID | Contract coverage |
| --- | --- |
| PB-061 | Setup status/start/verify and single Administrator projection |
| PB-062 | One-time TOTP enrollment and six-digit verification |
| PB-063 | Password+TOTP login, safe session and logout |
| PB-064 | Recovery start/verify and documented session revocation |
| PB-065 | 429 responses, generic errors and session expiry projection |
| PB-066 | GitHub App Manifest start, registration callback and installation callback |
| PB-067 | `docs/openapi.yaml` version 1.0.0 with 59 operations |

## Stitch interface compatibility

| Screen | Contract source | Deliberate differences |
| --- | --- | --- |
| Production Overview | `GET /overview`, IncidentSummary, ActivityEvent | Uses workflow-completed count, not “resolved” |
| Projects Inventory | `GET /projects` | Health sparkline remains presentation-only |
| Project Configuration | Project + monitoring endpoints, GitHub/Dokploy discovery | Built-in rules are read-only; no AI detector toggle; QVAC is global |
| Incident Detail | Incident, timeline, log events, Evidence, Remediation and PR | PR creation never means merge/deploy/resolution |
| Settings General | System status and diagnostics Operation | No Team or manual incident creation |
| Settings GitHub | GitHub account CRUD, PAT and App Manifest | Secrets never return to browser |
| Settings Dokploy | Server CRUD, discovery and connection test | Credential is write-only |
| Settings QVAC | Configuration/status/test | Endpoint restricted to local/private networks |
| Settings Automation | AutomationPolicy | Issue publication remains mandatory and is not a toggle |

## Explicit exclusions

- Team, invitations, RBAC and multi-tenancy.
- Manual incident creation.
- QVAC/AI-based log detection.
- Per-project QVAC model selection.
- Merge, deploy, rollback or a production `resolved` assertion.
