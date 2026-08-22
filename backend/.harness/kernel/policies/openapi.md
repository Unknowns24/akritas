# OpenAPI Policy

OpenAPI is the source of truth for public API contracts.

## Possible spec locations

Check the existing project first. Common locations:

```text
internal/adapter/rest/swagger/openapi.yaml
internal/adapter/rest/docs/openapi.yaml
docs/openapi.yaml
```

Do not create a second spec if one already exists.

## Mandatory updates

Update OpenAPI whenever changing:

- Endpoint path or method.
- Request body.
- Query params.
- Path params.
- Response schema.
- Error schema.
- Auth/security requirements.
- Permissions or semantics that affect clients.

## SemVer

Update `info.version` when the API contract changes.

- major: breaking changes.
- minor: compatible additions.
- patch: documentation/fixes that do not change contract.

## Agent rules

- Do not invent endpoints or fields.
- Do not implement backend API changes without updating OpenAPI.
- Do not implement frontend API calls not present in OpenAPI.
- If OpenAPI is missing required backend functionality, report the mismatch instead of guessing.
