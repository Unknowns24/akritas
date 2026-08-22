# Frontend Rendering Boundaries Policy

This policy defines Next.js App Router client/server boundaries for Akritas.

## Default boundary

Prefer Server Components when a component does not require browser-only state, effects, event handlers or browser APIs.

Use `"use client"` only at the smallest practical interactive boundary. Do not turn entire routes or large feature trees into Client Components without a concrete reason.

## API access

Use the project's shared API client and established data-access pattern. Do not introduce arbitrary `fetch`, Axios instances or Server Actions that bypass the documented API contract.

The frontend must consume the backend/OpenAPI contract and must not invent endpoints or payload fields.

## Credentials and secrets

The browser must never receive integration credentials or secrets for:

- GitHub;
- Dokploy;
- QVAC;
- repository access;
- internal service credentials.

Do not expose secrets through `NEXT_PUBLIC_*`, page props, serialized server-component data, client state, query parameters, localStorage or browser logs.

The UI may receive safe status/configuration projections such as connection state, display name, identifiers intended for display and non-secret capabilities.

## Environment variables

Read runtime configuration through the project's centralized runtime-env/config module. Do not scatter `process.env` access through features or app routes.

## Server-only code

Modules that handle server-only credentials or privileged operations must remain server-only and must not be imported by Client Components.
