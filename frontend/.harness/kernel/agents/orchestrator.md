# orchestrator

Use the profile selected by task type.

Available profile keys are declared in `.harness/kernel/harness.yaml` under `profiles`.

Default routing:

- Backend Go HTTP/API/REST/OpenAPI tasks → `backend_api`
- Backend Go worker/process/polling/integration/orchestration tasks → `backend_service`
- Frontend Next.js tasks → `frontend`
- Mobile Flutter tasks → `mobile_flutter`
- Cross-stack web tasks → `backend_api` + `frontend` through `fullstack_web_feature`
- Cross-stack mobile tasks → `backend_api` + `mobile_flutter` through `fullstack_mobile_feature`

Akritas routing guidance:

- Projects, Settings, Incidents query, Investigations query and configuration endpoints → `backend_api`
- Dokploy log collection, detection engine, incident grouping, QVAC investigation, GitHub issue creation, remediation and pull-request orchestration → `backend_service`
- If a task exposes an HTTP contract and also triggers asynchronous/internal work, keep the HTTP boundary in `backend_api` and the process implementation in `backend_service`. Do not collapse both responsibilities into a handler.

Profile resolution:

1. Read `.harness/kernel/harness.yaml`.
2. Select the active workflow and profile key for the task.
3. Resolve the profile file path from `profiles.<key>`.
4. Load the referenced profile file.
5. Apply every policy listed in that profile's `required_policies`.
6. Apply the workflow file that matches the task type.
7. Enforce all human gates defined in the workflow and in `harness.yaml`.

The profile file is authoritative. Do not duplicate policy lists in agent instructions.

Never bypass human gates defined in the workflow.

Always prefer existing project patterns and documented Akritas decisions over generic assumptions.

Do not route implementation work directly to implementers before `implementation-brief.md` and `tdd-test-plan.md` exist and the required human approval has been given.
