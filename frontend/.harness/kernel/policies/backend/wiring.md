# Backend Wiring Policy

## Entry point

The expected application entry point is:

```text
cmd/main.go
```

## Wiring order

For Go backends using this harness, preserve this initialization order unless the existing project clearly uses another one:

```text
config → DB + migrations + seeding → repositories → external adapters → services → usecases → REST router
```

## Rules

- Wiring belongs in the composition root, usually `cmd/main.go` or a dedicated bootstrap package if the project already has one.
- Do not instantiate concrete DB repositories inside handlers or usecases.
- Do not instantiate external SDK clients inside usecases unless the existing architecture explicitly does so. Prefer ports/adapters.
- Constructors should receive dependencies explicitly.
- Keep constructor return types aligned with ports when applicable.

## Typical constructor conventions

```go
func NewXUseCase(...) in.XUseCase
func NewXRepository(...) out.XRepository
func NewXAdapter(...) out.XAdapter
func NewXService(...) in.XService
```

## When adding a new dependency

Update the wiring in the same layer order:

1. Config/env if needed.
2. Concrete adapter/repository/service constructor.
3. Usecase constructor.
4. Handler/router constructor.
5. Tests/mocks/fakes.

If the new dependency affects deployment, update env examples and docs.
