# Architecture Decisions

Record durable project decisions here.

## 2026-08-22 — API contract v1

- `backend/docs/openapi.yaml` is the canonical frontend/backend contract.
- The contract uses OpenAPI 3.1.0, API version 1.0.0, base `/api/v1`, snake_case JSON, cursor pagination, reusable envelopes and opaque asynchronous operations.
- A GitHub Pull Request is the human-control boundary. The API does not merge, deploy or claim that production is resolved.

## 2026-08-22 — Single-administrator authentication

- Each installation has exactly one administrator.
- Setup and recovery are authorized by `AKRITAS_BOOTSTRAP_TOKEN`; it is never a TOTP seed.
- `AKRITAS_MASTER_KEY` encrypts TOTP and integration credentials and remains separate from the bootstrap token.
- Login requires email, Argon2id password verification and RFC 6238 TOTP. Sessions are opaque and stored server-side.

## 2026-08-22 — GitHub authentication methods

- GitHub.com connections support write-only Personal Access Tokens and GitHub App registration via the official Manifest flow.
- `github_account_type` (`personal|organization`) and `authentication_method` (`personal_access_token|github_app`) are independent concepts.
- Private keys, PATs, webhook secrets and installation tokens never appear in response DTOs.

## 2026-08-22 — Backend foundation

- The backend Go module is `github.com/Unknowns24/akritas/backend` and declares Go 1.26.
- The MVP domain starts as one flat `internal/core/domain` package with one file per cohesive concept.
- Akritas-owned identities use `github.com/google/uuid`; provider identifiers remain strings.
- Domain entities do not contain JSON/GORM tags or integration/authentication secrets.
- Domain error components reserve `0x401` through `0x406` for auth, integrations, project, incidents, investigations and remediation respectively.
