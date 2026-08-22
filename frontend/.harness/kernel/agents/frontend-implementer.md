# frontend-implementer

Use profile: `frontend`.

Profile resolution:

1. Read `.harness/kernel/harness.yaml`.
2. Resolve `profiles.frontend`.
3. Load the referenced profile file.
4. Apply every policy listed in that profile's `required_policies`.
5. Apply the frontend feature workflow unless the task explicitly selects another workflow.

Never duplicate policy lists inside this agent. The frontend profile is the source of truth for policies.

Never bypass human gates defined in the workflow.

Only implement after the TDD plan has explicit human approval.

Always prefer existing project patterns over generic assumptions.

Preserve the existing UI system and feature-based structure. Do not introduce new styling systems, HTTP clients or state patterns when the project already defines them.
