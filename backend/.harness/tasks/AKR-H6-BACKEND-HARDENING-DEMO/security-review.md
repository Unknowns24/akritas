# Security Review

## Implemented Security Controls

- QVAC tool output is redacted and bounded before returning to model context or Evidence.
- `evidencesafety` now reports whether redaction and truncation actually happened.
- Validation stdout/stderr excerpts are sanitized before persistence.
- PR content is generated from IDs and bounded sanitized text only.
- GitHub PR and Issue adapters resolve credentials internally and never expose provider response bodies on errors.
- Git operations use fixed argv and validate refs before execution.
- Demo fixture/runbook includes no secrets.

## Secret-Sink Coverage Added

- Validation output sentinel redaction.
- QVAC repository tool-output sentinel redaction.
- Commit correlation sanitization of author/message/url/sha fields.
- GitHub PR publisher `httptest` coverage for strict reconciliation and creation payload.

## Remaining Security Gaps

- REST projections for Remediation/PR are still absent, so API-level secret-sink tests for those surfaces are pending.
- Full logs/timeline/change-summary sanitization coverage across every sink listed in AKR-61 is not exhaustive yet.
- Startup recovery must avoid logging raw external errors when implemented.
- Workspace cleanup confinement still needs a root-owned workspace adapter.

