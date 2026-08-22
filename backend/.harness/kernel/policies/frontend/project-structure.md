# Frontend Project Structure Policy

Expected frontend stack by default:

- Next.js App Router inside `src/app/`.
- TypeScript.
- Feature-based architecture inside `src/`.
- CSS Modules for custom styles (`*.module.css`).
- Base UI when useful.
- Lucide Icons or a visually consistent modern icon set.
- Zod strongly recommended for runtime validation of DTOs and relevant backend responses.

For existing projects, always preserve the current stack and conventions. Do not introduce Tailwind, shadcn, Mantine, Radix wrappers or another UI kit unless the project already uses it or the task explicitly requests it.

## Sources of truth

Respect these sources in this order:

1. `docs/openapi.yaml` — backend contract: endpoints, DTOs, responses, errors and auth.
2. `DESIGN.md` or `design.md` — visual and UX contract.
3. `AGENTS.md` and `.harness/kernel/**` — technical architecture, workflow and quality rules.
4. Existing code patterns.

## Base structure

```text
src/
  app/
  core/
  features/
```

## `src/app/` — Next.js routing

`src/app/` contains only Next.js routing and framework files:

- `layout.tsx`
- `page.tsx`
- `loading.tsx`
- `error.tsx`
- `not-found.tsx`
- route groups such as `(public)` and `(private)`

### Rule: app must be thin

- Do not put business logic in `app`.
- Do not do complex requests directly in pages when the logic belongs to a feature.
- Do not define complex forms directly in `page.tsx`.
- Pages should import ready-to-use views from `features/*/views`.

Recommended pattern:

```ts
import { SolicitudesListView } from "@/features/solicitudes";

export default function Page() {
  return <SolicitudesListView />;
}
```

## `src/core/` — shared cross-feature code

`core` contains code reused by multiple features.

It must not contain logic that belongs to only one feature.

Suggested structure:

```text
src/core/
  config/
  libs/
  routes/
  ui/
  auth/
  errors/
```

Important rules:

- Shared UI goes in `src/core/ui`.
- Global helpers go in `src/core/libs`.
- Auth/session helpers go in `src/core/auth` or the existing auth location.
- Route builders go in `src/core/routes`.
- A shared API client belongs in `src/core/libs/api-client`.
- `core` must not import from `features`.

## `src/features/` — feature modules

Each feature should be autonomous and expose a clear public API.

Standard feature structure:

```text
src/features/<feature>/
  index.ts
  views/
  components/
  composites/
  hooks/
  services/
  schemas/
  types/
  mappers/
  context/
  storage/
  utils/
```

Use only the folders needed by the feature.

Descriptions:

- `views/`: screens ready to be imported from `app/`.
- `components/`: small UI components specific to the feature.
- `composites/`: composed sections coordinating multiple components.
- `hooks/`: feature-specific behavior and state orchestration.
- `services/`: backend calls.
- `schemas/`: runtime validation with Zod.
- `types/`: TypeScript types and DTOs.
- `mappers/`: DTO ↔ UI model transformations.
- `context/`: feature-internal contextual state.
- `storage/`: localStorage/sessionStorage/preferences.
- `utils/`: strictly internal feature helpers.

## Public surface per feature

Every feature should expose a public surface through `index.ts`.

Correct:

```ts
import { SolicitudesListView } from "@/features/solicitudes";
```

Incorrect:

```ts
import { SolicitudesListView } from "@/features/solicitudes/views/SolicitudesListView";
```

Rules:

- Other parts of the system should import only from `@/features/<feature>`.
- Do not import internals of another feature.
- Export only what is intentionally public.

## Coupling rules

Allowed direction:

```text
app → features → core
features → core
core → no features
```

Forbidden:

- `core` importing from `features`.
- A feature importing internals of another feature.
- Services importing UI components.
- Generic UI components with business logic.
- Business logic inside `page.tsx`.
- Duplicated API clients.
- Duplicated hardcoded routes.

If two features need the same thing:

- Generic UI → move to `core/ui`.
- Global helper → move to `core/libs`.
- Shared domain concern → evaluate `core` or a more explicit feature boundary.
- Casual reuse only → avoid premature abstraction.

## Import paths

Use `@/` alias pointing to `src/`.

Allowed:

```ts
import { Button } from "@/core/ui/components/Button";
import { SolicitudesListView } from "@/features/solicitudes";
```

Avoid long relative imports:

```ts
import { something } from "../../../../core/libs/something";
```

## Naming conventions

- Components: `PascalCase.tsx`
- Hooks: `useSomething.ts`
- Services: `*.service.ts`
- Schemas: `*.schemas.ts`
- Types: `*.types.ts`
- Mappers: `*.mapper.ts`
- Constants: `*.constants.ts`
- CSS Modules: `ComponentName.module.css`
- Views: `SomethingView`
- DTOs: suffix `Dto`
- Props: suffix `Props`
