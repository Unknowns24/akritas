# Backend Project Structure Policy

## Architecture

This backend follows hexagonal architecture / ports and adapters.

Allowed dependency direction:

```text
adapter/rest|db|external → usecase → core(domain + ports)
```

`service` may implement core ports/services and be consumed by `usecase`.

## Forbidden dependencies

- `internal/core/**` importing `internal/adapter/**`.
- `internal/core/**` importing HTTP, Chi, GORM packages, SQL drivers or external SDKs.
- `internal/usecase/**` importing concrete adapters.
- REST handlers using GORM or concrete repositories directly.
- HTTP DTOs leaking into domain/usecase contracts.
- GORM models being returned as REST responses.

Domain structs MAY contain declarative `gorm` tags when an accepted ADR defines
that convention. Tags do not authorize GORM imports, hooks, queries or
persistence behavior in core.

## Expected structure

```text
cmd/main.go
config/
internal/core/domain/
internal/core/ports/in/
internal/core/ports/out/
internal/core/ports/paging/
internal/usecase/
internal/service/
internal/adapter/db/<technology>/
internal/adapter/external/
internal/adapter/rest/
```

## REST adapter

Expected substructure:

```text
internal/adapter/rest/dto/
internal/adapter/rest/handler/
internal/adapter/rest/middleware/
internal/adapter/rest/router/
internal/adapter/rest/swagger/ or internal/adapter/rest/docs/
```

REST rules:

- Handlers should be thin.
- Parse/validate transport input in REST.
- Map DTOs to usecase/domain input types.
- Call usecases through ports/interfaces.
- Map domain errors to HTTP through the existing mapper/middleware.
- Never put business rules in handlers.
- REST transport structs MUST use the `DTO` suffix and one DTO struct per file.
- DTO/domain/usecase conversion MUST live under `internal/adapter/rest/mapper/`;
  mapping files must have one conversion responsibility.

## Feature implementation order

When adding or extending a backend feature, use this order unless the existing module strongly suggests otherwise:

1. Domain model or domain behavior.
2. Ports/contracts.
3. Usecase.
4. Internal service if reusable orchestration/process logic is needed.
5. DB or external adapters.
6. REST DTOs.
7. REST handlers.
8. Router registration.
9. OpenAPI update.
10. Tests.
11. Architecture review.
12. Security review.

Do not start from the HTTP handler unless the task is purely about delivery/transport.
