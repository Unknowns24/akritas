# Flutter Mobile TDD Policy

Flutter tasks managed by the harness must follow the same human-gated TDD flow as backend and frontend tasks.

## Flow

1. Analyze the task and current project patterns.
2. Produce `implementation-brief.md`.
3. Produce `tdd-test-plan.md`.
4. Wait for explicit human approval.
5. Implement tests first or alongside the smallest production change needed to make them pass.
6. Implement production code.
7. Run validation commands.
8. Produce summary/review artifacts.

Do not implement production code before the TDD plan is approved.

## Test priorities

Prefer the lowest useful test level.

### Domain tests

Use for:

- Entity/value object rules.
- Status mapping.
- Calculations.
- Pure business decisions.
- Permission/scope decisions represented in app logic.

### Application/provider/notifier tests

Use for:

- Use case orchestration.
- Riverpod providers/notifiers.
- Loading/data/error transitions.
- Retry behavior.
- Auth/session bootstrap decisions.

### Data tests

Use for:

- DTO parsing.
- Mapper behavior.
- API envelope parsing.
- Pagination parsing.
- API error normalization.
- Repository mapping from DTO to domain.

### Widget tests

Use for:

- Important user flows.
- Loading/empty/error/data UI states.
- Forms and validation.
- Navigation triggers when practical.
- Retry buttons.

Avoid brittle pixel-perfect widget tests unless the project already uses golden tests or the task is specifically visual.

## Recommended tooling

Use existing project tooling first.

Default Flutter commands:

```text
flutter test
flutter analyze
```

If the project uses code generation:

```text
dart run build_runner build --delete-conflicting-outputs
```

only when required by the project.

If the project has custom scripts, use those instead of inventing new ones.

## Test plan content

`tdd-test-plan.md` must include:

- Scope.
- Test files to add/update.
- Cases covered.
- Expected failing tests before implementation.
- Fixtures/mocks/fakes needed.
- Validation commands.
- Open questions / human approval notes.

## Fakes and fixtures

- Prefer fake repositories for application/widget tests.
- Prefer fake `ApiClient` or mocked transport for data tests.
- Keep fixtures in `test/fixtures` or the existing test fixture location.
- Do not use production mocks as test fixtures unless moved/renamed clearly.

## Human gate

The TDD plan requires explicit human approval before implementation.

Approval examples:

```text
aprobado
ok implementá
seguí con ese plan
```

Without approval, stop after the plan.
