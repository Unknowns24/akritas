# Frontend Data Contract Policy

The frontend must not imagine the backend.

## Source of truth

Use:

```text
docs/openapi.yaml
```

as the source of truth for:

- endpoints
- request DTOs
- response DTOs
- errors
- authentication
- permissions when documented
- pagination shape

## Rules

- Do not invent endpoints.
- Do not invent field names.
- Do not assume response shapes not documented in OpenAPI.
- If backend behavior and OpenAPI conflict, report the inconsistency and do not hide it with hacks.
- If backend changes, update OpenAPI first, then update frontend.
- Services must use `src/core/libs/api-client` or the existing shared API client.
- Do not create ad-hoc `fetch`/`axios` clients inside features.
- Validate relevant backend responses with Zod schemas when the contract matters.
- DTO types should use the `Dto` suffix.

## Services

`services/` should:

- define backend calls only.
- use the shared API client.
- return typed DTOs.
- validate responses when relevant.
- contain no UI logic.

For large features, services may be split by operation to preserve SRP.

Examples:

```text
solicitudes.service.ts
create-solicitud.service.ts
get-solicitudes.service.ts
```

## Schemas

`schemas/` should:

- validate backend payloads.
- validate forms when applicable.
- stay close to the feature that uses them.

## Types

`types/` should include:

- DTOs
- UI types
- enums
- derived types

Avoid duplicating types without a reason.

## Errors

- Centralize error normalization in `core/errors` or `core/libs/api-client`.
- Show human-friendly messages.
- Do not expose raw technical errors to users.
- Preserve technical details only where appropriate for logging/debugging.
- Respect the backend error envelope when it exists.
