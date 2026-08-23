# Architecture Review

Profile: `backend_service` with OpenAPI policy verification.

## Passed Architecture Points

- Domain remains free of external clients, GORM behavior, HTTP and process execution.
- Git and GitHub mutations are exposed as narrow output-port methods, not generic shell/provider escape hatches.
- Git commands use `exec.CommandContext` with fixed argv.
- GitHub credentials are still resolved only inside the GitHub adapter.
- Commit correlation is a focused service under `internal/service/commitcorrelation`.
- Sanitization remains centralized in `internal/service/evidencesafety`.
- Remediation PR orchestration is in the use case layer and calls adapters through ports.
- Bootstrap now constructs real Remediation persistence/adapters/use case.

## Architecture Risks Remaining

- REST/OpenAPI/runtime correspondence is still incomplete for Remediation, PR and QVAC lifecycle paths.
- Remediation orchestration is still split between manual commands and future automatic Investigation-driven trigger.
- Recovery/reconciliation should be extracted further before adding worker startup recovery.
- Workspace root confinement/provisioning is still not a complete owned adapter; current Git adapter operates on caller-provided workspace paths.

