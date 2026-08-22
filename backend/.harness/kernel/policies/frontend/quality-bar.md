# Frontend Quality Bar

Before considering a frontend task complete, verify:

- TypeScript types are correct; avoid `any` by default.
- No duplicated API client; use `src/core/libs/api-client`.
- No hardcoded routes; use `src/core/routes`.
- UI is consistent with `DESIGN.md` / `design.md`.
- Colors are used through tokens where available.
- `src/app` remains thin and imports views from features.
- Services align with `docs/openapi.yaml`.
- Relevant payloads/responses are validated with schemas when appropriate.
- Paginated lists follow `docs/pagination_frontend.md` / uker cursor pagination.
- No imports to internals of another feature.
- Components, hooks, services, mappers and schemas respect SRP.
- Loading, empty and error states are covered for remote data.
- Permission denied / restricted access is handled when applicable.
- Basic accessibility is considered.
- Client/server boundaries follow `frontend/rendering-boundaries.md`.
- No integration credentials or secrets are exposed to browser code.
- No unnecessary monolithic components or views.
- No premature global abstractions.
