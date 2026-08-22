# Project Summary — Akritas

Akritas is a hackathon MVP that monitors configured projects, detects incidents from Dokploy logs, investigates them locally with QVAC, always publishes a GitHub Issue, and optionally produces a validated remediation Pull Request when the incident is fixable.

Primary loop:

```text
Dokploy → Incident → QVAC → GitHub Issue → Optional Fix → Pull Request
```

Architecture guidance:

- Backend: Go, hexagonal architecture.
- `backend_api`: HTTP/API/control-plane and query/configuration boundaries.
- `backend_service`: log collection, detection, incident grouping, investigation and remediation workflows.
- Frontend: Next.js App Router, feature-based architecture.
- API contract: OpenAPI is authoritative for frontend/backend integration.
- Canonical API contract: `backend/docs/openapi.yaml`, OpenAPI 3.1.0, API v1.0.0 under `/api/v1`.
- Persistence: GORM repositories + gormigrate/v2 migrations against domain structs (`gorm` tags). The database adapter must not introduce duplicate table models.
- Backend foundation: Go 1.26 module with an MVP domain under `internal/core/domain`. HTTP JSON stays in REST DTOs; the domain does not import GORM.
- QVAC inference must remain local for the hackathon requirements.
- GitHub, Dokploy, QVAC, Git and filesystem access are infrastructure adapters behind output ports.
- Integration credentials are backend-only and must never reach browser code.
- An installation has one administrator enrolled with password + TOTP; opaque sessions are server-side and recovery requires the deployment bootstrap token.
- GitHub connections support either a write-only PAT or the GitHub App Manifest flow. Account ownership (`personal|organization`) is independent from authentication method.

MVP execution order:

1. Control Plane
2. Detection + Incidents
3. QVAC Investigation
4. GitHub Issue
5. Remediation + Pull Request
6. Hardening + Demo

Prefer deterministic behavior and documented domain decisions over agent-invented abstractions.
