# frontend-tdd

Use profile: `frontend`.

Profile resolution:

1. Read `.harness/kernel/harness.yaml`.
2. Resolve `profiles.frontend`.
3. Load the referenced profile file.
4. Apply every policy listed in that profile's `required_policies`.
5. Apply the frontend feature workflow unless the task explicitly selects another workflow.

Never duplicate policy lists inside this agent. The frontend profile is the source of truth for policies.

Never bypass human gates defined in the workflow.

Always prefer existing project patterns over generic assumptions.

Before implementation, produce or validate `implementation-brief.md` and `tdd-test-plan.md`. The TDD plan must define tests first and must be explicitly approved before implementation starts.
