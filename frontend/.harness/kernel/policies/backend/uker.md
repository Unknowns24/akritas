# Backend Uker Policy

This project may use `github.com/unknowns24/uker` as a utility kernel.

Consult project docs first if they exist:

```text
docs/uker_bodyparser.md
docs/uker_pagination.md
.harness/kernel/docs/uker/bodyparser.md
.harness/kernel/docs/uker/pagination.md
```

## BodyParser

Use `httpx.BodyParser` for HTTP body parsing when the project uses that pattern.

Rules:

- Always pass a pointer to struct.
- Always define `json` tags for public API DTOs.
- Use `uker:"required"` only for transport-level required fields.
- Do not treat `uker:"required"` as business validation.
- Add explicit usecase/domain validation for business rules.
- Use `WithBase64Data()` only when the endpoint contract requires base64 inside `data`.

## Pagination

For list endpoints, prefer cursor-based signed pagination.

Heuristic:

- If a usecase receives `pagination.Params`, the handler should parse query params using the uker pagination flow.
- If a usecase returns `([]T, int64, error)`, the REST layer should build a paginated response.

Rules:

- First request may include `limit`, `sort` and filters.
- Subsequent requests using `cursor` must not modify filters/sort/limit.
- Use blocked filters for backend-owned scopes, such as `user_id`, tenant IDs or authenticated user IDs.
- Do not introduce offset/limit pagination in new endpoints unless explicitly required.
- Small controlled catalogs may return full lists only if the existing module already does that or the task explicitly requires it.
