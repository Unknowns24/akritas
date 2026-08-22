# Flutter Mobile Data Contract Policy

The mobile app must not imagine the backend.

## Source of truth

Use, in order:

```text
docs/openapi.yaml
openapi.yaml
backend documentation explicitly referenced by the task
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
- If backend changes, update OpenAPI first, then update Flutter integration.
- Do not consume administrative endpoints from public/user apps.
- If an endpoint has a permission marker and the app is not administrative, treat it as forbidden unless explicitly whitelisted.
- Do not calculate derived data in Flutter if the backend exposes a specific endpoint for it.
- Do not wrap request bodies in `{ data: ... }` unless the contract explicitly requires it.

## ApiClient

`ApiClient` is the only HTTP gateway.

It should handle:

- Base URL.
- Common headers.
- Bearer token.
- CookieJar.
- Timeouts.
- Request/response interceptors.
- Envelopes.
- Pagination.
- API error normalization.
- Unauthorized/session cleanup.

Forbidden:

- Creating `Dio()` inside features.
- Duplicating base URLs in features.
- Calling external services directly from screens.
- Logging tokens, cookies or passwords.

## Envelopes

If the backend uses a simple envelope:

```json
{
  "data": {},
  "message": "..."
}
```

parse it in `core/api` or `data`, never in widgets.

If the backend uses cursor pagination:

```json
{
  "data": [],
  "paging": {}
}
```

use a shared `CursorPage<T>` / `CursorPaging` type or the project equivalent.

If the backend uses page/size pagination inside `data`:

```json
{
  "data": {
    "items": [],
    "page_number": 0,
    "page_size": 50,
    "total_elements": 100,
    "total_pages": 2,
    "has_next": true,
    "has_previous": false
  },
  "message": "..."
}
```

use a shared page type such as `Page<T>` / `ERPPage<T>` or the project equivalent.

## Errors

If the backend returns:

```json
{
  "code": "...",
  "message": "...",
  "user_message": "...",
  "request_id": "..."
}
```

show errors to the user in this order:

1. `user_message`
2. `message`
3. Generic Spanish message

Rules:

- Do not expose raw technical errors to end users.
- Preserve `request_id` in the normalized error when available.
- Keep debug details out of UI unless the project intentionally has a debug surface.

## DTOs and mappers

Each integrated feature should use:

```text
data/dto/
data/mapper/
data/remote/
data/<feature>_repository_impl.dart
```

Rules:

- DTO = exact JSON shape.
- Domain entity = model comfortable for app/business/UI.
- Mapper = explicit translation.
- UI never receives DTOs.
- Avoid `Map<String, dynamic>` outside `data`.
- Keep visual formatting in presentation/core UI utils, not DTOs.

## Remote datasources

Remote datasources may exist when useful for SRP.

They should:

- Call only `ApiClient`.
- Return DTOs or parsed API DTO pages.
- Contain no UI logic.
- Contain no business decisions beyond transport mapping.

Repositories should map datasource results into domain entities.

## Scope blocks per repository

Each mobile repository should declare in its root `AGENTS.md`:

```text
Allowed endpoints:
  - ...
Forbidden endpoints:
  - ...
```

When scope is missing, the agent must infer conservatively from the app type and OpenAPI permissions, and document uncertainty in the implementation brief.
