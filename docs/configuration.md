# Akritas — Runtime and deployment configuration

## Purpose

This document defines runtime values that affect the public API, authentication
and integration callbacks. Secrets shown as placeholders must be supplied through
the deployment secret mechanism and never committed.

## Required values

### `AKRITAS_PUBLIC_URL`

Canonical HTTPS origin used for redirects, Origin validation and GitHub manifest
URLs.

Example:

```text
https://akritas.example.com
```

The value must not include credentials, query parameters or fragments.

### `AKRITAS_MASTER_KEY`

High-entropy encryption root for the Credential Store and encrypted TOTP seeds.
It is never persisted by Akritas or exposed through API/diagnostics.

Losing this key makes stored integration and TOTP secrets inaccessible. Replacing
it requires an explicit secret-rotation migration, which is outside the MVP.

### `AKRITAS_BOOTSTRAP_TOKEN`

High-entropy administrator setup/recovery capability. It must be independent from
`AKRITAS_MASTER_KEY` and every generated TOTP seed.

It is read at runtime, compared in constant time, rate limited at the HTTP boundary
and never stored or logged.

### `AKRITAS_POSTGRES_DSN`

PostgreSQL connection string used by the persistence adapter.

Example:

```text
postgres://akritas:akritas@127.0.0.1:5432/akritas?sslmode=disable
```

The value must not be logged. Tests may use `AKRITAS_POSTGRES_TEST_DSN` to reuse an
already running instance instead of starting an embedded Postgres.

### `AKRITAS_PAGINATION_SECRET`

High-entropy key used to sign cursor pagination payloads. Rotating it invalidates
outstanding cursors but does not alter persisted application data.

## Session defaults

```text
AKRITAS_SESSION_IDLE_TTL=12h
AKRITAS_SESSION_ABSOLUTE_TTL=168h
AKRITAS_SESSION_COOKIE_SECURE=true
```

Production cookie attributes:

```text
HttpOnly; Secure; SameSite=Lax; Path=/
```

Local HTTP development may disable `Secure`, but production validation must reject
that configuration when `AKRITAS_PUBLIC_URL` uses HTTPS.

## Browser origins

Production serves frontend and API under the same site. Development may define:

```text
AKRITAS_ALLOWED_ORIGINS=http://localhost:3000
```

Rules:

- exact origins only;
- no wildcard when credentials are enabled;
- mutations require a matching `Origin`;
- callbacks validate one-time `state` independently of cookies;
- preflight responses never expose provider credentials.

## GitHub App callbacks

The manifest is generated from `AKRITAS_PUBLIC_URL`:

```text
redirect_url = <public-url>/api/v1/integrations/github/app-manifest/callback
setup_url    = <public-url>/api/v1/integrations/github/app-installations/callback
```

The App is private, webhooks are inactive and requested permissions are limited to
metadata read, contents read/write, issues write and pull requests write.

GitHub.com is the only allowed GitHub API/HTML host in the MVP. Redirect targets
must be selected by the backend rather than accepted from browser input.

## QVAC

QVAC endpoint and optional credentials are managed through the authenticated API.
The backend accepts only loopback or private-network destinations after DNS/IP
validation and rejects redirects to public addresses.

The UI receives endpoint, auth type, `credential_configured`, runtime/model/version
and status. It never receives bearer/basic credentials.

## Logging and diagnostics

Configuration validation may report missing variable names, but must never print
their values. Diagnostics expose normalized component states and request IDs only;
raw GitHub, Dokploy or QVAC responses remain internal.
