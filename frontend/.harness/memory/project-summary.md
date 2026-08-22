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
- Persistence: GORM repositories + gormigrate/v2 migrations against domain structs (`gorm` tags). The database adapter must not introduce duplicate table models.
- QVAC inference must remain local for the hackathon requirements.
- GitHub, Dokploy, QVAC, Git and filesystem access are infrastructure adapters behind output ports.
- Integration credentials are backend-only and must never reach browser code.

MVP execution order:

1. Control Plane
2. Detection + Incidents
3. QVAC Investigation
4. GitHub Issue
5. Remediation + Pull Request
6. Hardening + Demo

Prefer deterministic behavior and documented domain decisions over agent-invented abstractions.
