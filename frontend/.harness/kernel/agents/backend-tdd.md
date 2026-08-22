# backend-tdd

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

Always prefer existing project patterns and documented Akritas decisions over generic assumptions.

Before implementation, produce or validate `implementation-brief.md` and `tdd-test-plan.md`. The TDD plan must define tests first and must be explicitly approved before implementation starts.

For integration-heavy flows, tests should target behavior through ports/fakes rather than requiring live GitHub, Dokploy, QVAC or shell execution.
