# Security Review

## Summary

The primary security requirement for this task was preventing command
injection / arbitrary shell execution from external or untrusted input
(QVAC output, Evidence, repository content) reaching the git or validation
runner adapters. Both adapters were designed and tested against this
requirement directly.

## Auth / permissions

Not applicable — this task adds no REST endpoints, no authentication or
authorization surface.

## Input validation

- `internal/adapter/external/git`: branch/base ref names are validated by
  `validateRefName` (rejects empty, leading `-`, `..`, and control
  characters) BEFORE any `exec.Command` runs, closing the git-specific
  "flag-like ref name" argument-injection surface (e.g.
  `-Xupload-pack=/bin/sh`) that fixed-argv alone does not prevent.
  Regression test: `TestClientCreateBranch/ref_name_injection_attempt_is_rejected_before_any_git_command_runs`.
- `internal/adapter/external/validationrunner`: `ValidationCommand` is a
  closed, unexported-argv-mapped enum (`go_test`/`go_vet`/`go_build`);
  there is no `Run(command string)`-shaped method anywhere in the port or
  adapter. The only caller of `Run` is `usecase/remediation`, which only
  ever passes `ValidationCommand` values produced by
  `service/validationpolicy`'s own stack-detection logic — never a value
  derived from QVAC output, Evidence, or repository file content. An
  unrecognized command value returns an explicit error
  (`ErrValidationExecutionFailed`) rather than silently no-op-ing or
  falling back to a default command.
- Both adapters use `exec.CommandContext` with fixed argv slices — never
  `sh -c` or any shell string interpolation.

## Data exposure

- `ValidationResult.OutputRedacted` is always `true` (enforced by the
  domain constructor/transition methods themselves); this task's usecase
  never attempts to set it `false`.
- Output/summary are truncated (`truncateWithMarker`) to the domain's
  documented limits (5000/50000 bytes) before persistence, with a visible
  truncation marker — no unbounded stdout/stderr persisted.
- No secrets are introduced into logs, error messages, or persisted
  output by this task. Neither adapter touches `CredentialStore` or any
  GitHub/Dokploy token.

## Error leakage

Repository `mapError` functions wrap raw GORM/pgx errors behind
`dberrors.ErrRemediationPersistence`/`ErrValidationResultPersistence`
rather than returning driver errors directly (verified by
`TestRepositoryCreateForeignKeyViolation`/`TestRepositoryCreateDuplicateID`
asserting the raw error type does not leak).

## File uploads/downloads, path traversal, SSRF, SQL injection, XSS, CORS, tokens

Not applicable — no HTTP surface, no file upload/download, no outbound
HTTP calls, no new SQL string concatenation (all queries use GORM
parameter binding), no frontend code.

## Findings

None.

## Result
pass
