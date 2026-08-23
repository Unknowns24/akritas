# Architecture Decisions

Record durable project decisions here.

## 2026-08-22 — API contract v1

- `backend/docs/openapi.yaml` is the canonical frontend/backend contract.
- The contract uses OpenAPI 3.1.0, API version 1.4.0 and base `/api/v1`. Version 1.4.0 completes the Project lifecycle, including hard deletion while monitoring is disabled; existing paths and payloads remain compatible.
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
- Persistible domain entities may contain passive GORM tags under ADR-012, but core never imports GORM or owns repository behavior. Domain/public types never contain integration/authentication secrets.
- Domain error components reserve `0x401` through `0x406` for auth, integrations, project, incidents, investigations and remediation respectively. `0x407` is reserved for `Operation`, the generic async-command entity shared by investigation and future remediation/pull_request flows.

## 2026-08-22 — PostgreSQL persistence

- PostgreSQL is the canonical relational database for the Akritas backend.
- GORM imports, queries and repository behavior remain confined to `internal/adapter/db/postgres`; tagged domain entities are persisted directly and schema evolution uses ordered gormigrate migrations with no global `AutoMigrate`.
- The encrypted Credential Store shares PostgreSQL with domain persistence but uses dedicated tables/repositories; `AKRITAS_MASTER_KEY` remains outside the database.
- Integration and migration tests use real PostgreSQL rather than SQLite as a semantic substitute.

## 2026-08-22 — Integration implementation conventions

- Runtime configuration is centralized in `config/config.go` with a local Viper instance; environment overrides optional `app.env` values and invalid security configuration fails closed.
- Uker v1.2.2 is the only pagination/cursor implementation. Core exposes its parameter alias, while REST owns parsing/signing and adapters translate provider boundaries.
- REST, database and external-adapter errors are declared by their owning layer. Domain retains domain and use-case errors through the common `domain.Error` contract.
- REST contract types use the `DTO` suffix, one structure per file and packages
  grouped by feature under `rest/dto/<feature>`; shared envelopes live in
  `rest/dto/common`, and mapping responsibilities live in `rest/mapper`.

## 2026-08-22 — Project lifecycle and integration snapshots

- Project names are unique case-insensitively and a Dokploy application can be associated with only one Project.
- Repository/application snapshots are resolved from the configured providers before create, association updates and monitoring activation; placeholders are not production data.
- The requested GitHub `default_branch` must match provider metadata.
- Association changes and deletion require monitoring disabled. Every enabled MonitoringConfiguration replacement revalidates both providers and returns the Project to `starting`.

## 2026-08-23 — GitHub Issue publication idempotency

- H4 creates one GitHub Issue per completed `Investigation`, not one per `Incident`.
- `GitHubIssueReference` is durable PostgreSQL state linked to both `Incident` and `Investigation`; a preexisting reference for the same Investigation is the idempotency boundary and prevents a second publish.
- `Incident.github_issue_reference` remains a singular REST projection of the most recent Issue for that Incident.
- `resolution_status = requires_human` completes the Incident after IssueReference persistence; `resolution_status = fixable` leaves the Incident in `publishing_issue` for H5 remediation.
- GitHub publication happens outside PostgreSQL transactions. The workflow uses a transaction before the external call to persist QVAC completion and a second transaction after the call to persist IssueReference and Operation completion.
