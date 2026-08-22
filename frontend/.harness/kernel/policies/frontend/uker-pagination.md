# Frontend Uker Pagination Policy

For paginated lists, consume cursor-based pagination provided by `uker/pagination`.

Reference documentation should live in one of:

```text
docs/pagination_frontend.md
docs/uker_pagination.md
.harness/kernel/docs/uker/pagination.md
```

Use whichever exists in the project.

## Critical rule

Do not implement alternative `offset`, `page` or `per_page` pagination unless OpenAPI explicitly documents that endpoint using that style.

## Query params

First page may use:

- `limit`
- `sort`
- filters with `<field>_<op>`
- OR filters with `<field1>,<field2>_<op>`
- no `cursor`

Supported operators:

```text
eq, neq, lt, lte, gt, gte, like, in, nin
```

Example:

```http
GET /solicitudes?limit=20&sort=created_at:desc&estado_eq=pendiente&empleado_nombre_like=juan
```

Next/previous pages must use only the cursor received from backend:

```http
GET /solicitudes?cursor=eyJ...
```

## Cursor rule

When using `cursor`, resend it exactly as received.

Do not resend or modify:

- `limit`
- `sort`
- filters
- search
- other params encoded in the cursor

The backend validates the cursor signature and can reject inconsistent requests.

## Expected response shape

```json
{
  "data": [],
  "paging": {
    "limit": 20,
    "total": 120,
    "has_more": true,
    "next_cursor": "eyJ...",
    "prev_cursor": ""
  }
}
```

## Frontend abstraction

Centralize pagination helpers in:

```text
src/core/libs/pagination/
  buildPaginationQuery.ts
  pagination.types.ts
  cursorPaging.ts
```

or the existing equivalent.

Responsibilities:

- Serialize `limit`, `sort` and filters consistently.
- Resend `next_cursor` and `prev_cursor` without altering them.
- Avoid combining `cursor` with filters or sort.
- Type the shared `paging` structure.
- Centralize operator names to avoid magic strings.

## Feature usage

Each feature should use common pagination helpers.

Views should not manually build complex query strings or know signed cursor internals.
