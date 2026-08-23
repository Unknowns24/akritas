# Architecture review — AKR-QVAC-INFERENCE

## Verdict
pass

## Notes
- QVAC/GitHub remain behind `ports/out` and adapter boundaries.
- No domain/usecase imports of HTTP/SDK types beyond existing ports.
- Tool allowlist and structured result parsing live in the QVAC adapter / investigationtools service.
