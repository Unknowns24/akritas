# mobile-flutter-tdd

Use profile: `mobile_flutter`.

Profile resolution:

1. Read `.harness/kernel/harness.yaml`.
2. Resolve `profiles.mobile_flutter`.
3. Load the referenced profile file.
4. Apply every policy listed in that profile's `required_policies`.
5. Apply the mobile Flutter feature workflow unless the task explicitly selects another workflow.

Never duplicate policy lists inside this agent. The mobile Flutter profile is the source of truth for policies.

Never bypass human gates defined in the workflow.

Always prefer existing project patterns over generic assumptions.

Before implementation, produce or validate `implementation-brief.md` and `tdd-test-plan.md`. For Flutter tasks, design tests around domain/application/data/widget boundaries and include `flutter test` / `flutter analyze` or the repository's existing validation commands.
