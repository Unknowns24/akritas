# Akritas — Runtime and deployment configuration

## Purpose

This document defines runtime values that affect the public API, authentication
and integration callbacks. Secrets shown as placeholders must be supplied through
the deployment secret mechanism and never committed.

## Required values

### `AKRITAS_DATABASE_URL`

PostgreSQL connection URL for domain metadata, GitHub App registration attempts and the encrypted Credential Store. PostgreSQL is the only persistent engine; schema changes run exclusively through the ordered gormigrate registry.

The URL is secret deployment configuration and must never be logged.

### `AKRITAS_PUBLIC_URL`

Canonical public origin used for redirects, Origin validation and GitHub
manifest URLs. Production deployments must use HTTPS. Local development may use
HTTP only with loopback hosts such as `localhost`, `127.0.0.1` or `::1`.

Example:

```text
https://akritas.example.com
http://localhost:8080
```

The value must not include credentials, query parameters or fragments.

### `AKRITAS_MASTER_KEY`

High-entropy encryption root for the Credential Store and encrypted TOTP seeds.
It is never persisted by Akritas or exposed through API/diagnostics.

The value must be standard Base64 encoding of exactly 32 bytes. Startup fails closed when it is missing or malformed.

Losing this key makes stored integration and TOTP secrets inaccessible. Replacing
it requires an explicit secret-rotation migration, which is outside the MVP.

### `AKRITAS_BOOTSTRAP_TOKEN`

High-entropy administrator setup/recovery capability. It must be independent from
`AKRITAS_MASTER_KEY` and every generated TOTP seed.

It is read at runtime, compared in constant time, rate limited at the HTTP boundary
and never stored or logged.

### `AKRITAS_PAGINATION_SECRET`

High-entropy key used to sign cursor pagination payloads. Rotating it invalidates
outstanding cursors but does not alter persisted application data.

### Optional configuration file and defaults

Runtime configuration is loaded centrally through Viper. An optional `app.env`
file may provide local defaults; environment variables always take precedence.
Secret values are never printed when validation fails.

```text
AKRITAS_PAGINATION_TTL=15m
AKRITAS_DB_MAX_OPEN_CONNECTIONS=10
AKRITAS_DB_MAX_IDLE_CONNECTIONS=5
AKRITAS_DB_CONNECTION_MAX_LIFETIME=30m
```

`AKRITAS_PAGINATION_TTL` limits signed cursor lifetime. Database pool values must
be positive and the idle limit cannot exceed the open-connection limit.

## Session defaults

```text
AKRITAS_SESSION_IDLE_TTL=12h
AKRITAS_SESSION_ABSOLUTE_TTL=168h
AKRITAS_SESSION_COOKIE_SECURE=true
AKRITAS_SESSION_COOKIE_SAME_SITE=lax
AKRITAS_AUTH_RATE_LIMIT_ATTEMPTS=5
AKRITAS_AUTH_RATE_LIMIT_WINDOW=15m
AKRITAS_AUTH_RATE_LIMIT_MAX_KEYS=4096
```

Production cookie attributes:

```text
HttpOnly; Secure; SameSite=Lax; Path=/
```

`AKRITAS_SESSION_COOKIE_SAME_SITE` accepts `lax`, `strict` or `none` (case
insensitive) and defaults to `lax`. The selected mode is used both when issuing
and expiring the session cookie. Startup fails closed for any other value or
when `Secure` is disabled; therefore `none` still requires HTTPS for browsers to
store and send the session cookie. Local development may configure
`AKRITAS_PUBLIC_URL` with HTTP loopback, but authenticated browser flows that
need `SameSite=None` still require HTTPS or an equivalent secure reverse proxy.

Authentication entry points use separate in-memory fixed-window limiters for
setup, enrollment verification, login and recovery. The MVP runtime is a single
backend process; these limits do not coordinate across replicas or survive a
restart. `MAX_KEYS` bounds memory and new keys fail closed while all current
buckets are active. Client identity uses the direct peer address and never
trusts forwarding headers.

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
