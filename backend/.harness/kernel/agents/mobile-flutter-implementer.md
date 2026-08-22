# mobile-flutter-implementer

Use profile: `mobile_flutter`.

Profile resolution:

1. Read `.harness/kernel/harness.yaml`.
2. Resolve `profiles.mobile_flutter`.
3. Load the referenced profile file.
4. Apply every policy listed in that profile's `required_policies`.
5. Apply the mobile Flutter feature workflow unless the task explicitly selects another workflow.

Never duplicate policy lists inside this agent. The mobile Flutter profile is the source of truth for policies.

Never bypass human gates defined in the workflow.

Only implement after the TDD plan has explicit human approval.

Always prefer existing project patterns over generic assumptions.

Preserve existing Flutter UI unless the task explicitly asks for a redesign. Do not introduce a new state manager, router or HTTP client when the project already has Riverpod, GoRouter and Dio/ApiClient patterns.
