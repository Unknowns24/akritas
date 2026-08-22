# Frontend Runtime Env Policy

Dashboards may be built in CI, published as Docker images and deployed later in different environments. Environment variables that change between dev, staging, production or client installations must be resolved at runtime, not frozen at build time.

## Rule

Use `next-runtime-env` for frontend/runtime configuration when the project uses this deployment model.

## Required pattern

Centralize env access in:

```text
src/core/config/env.ts
```

or the existing equivalent location.

Example:

```ts
import { env } from "next-runtime-env";

export const appEnv = {
  apiBaseUrl: env("NEXT_PUBLIC_API_BASE_URL"),
};
```

## Rules

- Do not spread `process.env.X` across features, views or services.
- Features, services and API clients must consume config through `src/core/config/env.ts` or the existing centralized config module.
- Validate critical envs at app startup or config initialization when practical.
- If a value can change between deployments of the same already-built Docker image, it must be runtime env.
- Do not introduce build-time-only env usage unless the variable is truly build-time-only and documented.
