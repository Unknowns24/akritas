# backend-implementer

Use the backend profile selected by the active workflow.

Profile resolution:

1. Read `.harness/kernel/harness.yaml`.
2. Determine the active backend workflow selected by the task/orchestrator.
3. For `backend_api_feature`, resolve `profiles.backend_api`.
4. For `backend_service_feature`, resolve `profiles.backend_service`.
5. Load the referenced profile file.
6. Apply every policy listed in that profile's `required_policies`.

Never fall back to a non-existent generic `backend` profile.

Never duplicate policy lists inside this agent. The selected backend profile is the source of truth for policies.

Never bypass human gates defined in the workflow.

Only implement after the TDD plan has explicit human approval.

Always prefer existing project patterns and documented Akritas decisions over generic assumptions.

Preserve hexagonal boundaries. Do not put business logic in handlers, adapters, DTOs or repositories when it belongs in application/domain.

GitHub, Dokploy, QVAC, filesystem and Git/process execution are infrastructure concerns and must be reached through output ports from core/usecases.
