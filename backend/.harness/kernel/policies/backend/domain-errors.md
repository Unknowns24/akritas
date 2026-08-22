# Backend Domain Errors Policy

Backend errors must use explicit enriched domain errors.

## Code format

```text
DxAAABBBT
```

Where:

- `D`: scope, one uppercase hex nibble.
- `x`: literal separator.
- `AAA`: component, three uppercase hex nibbles.
- `BBB`: cause, three uppercase hex nibbles, incremental per component.
- `T`: type, one letter used to map HTTP status.

Example:

```text
1x5B2014N
```

## Scope `D`

- `0`: platform
- `1`: tenant
- `2`: shared
- `F`: reserved

## Layer encoded in `AAA[0]`

- `1`: REST
- `2`: DB
- `3`: external adapters
- `4`: domain
- `5`: usecase
- `6`: service

## Type `T`

- `V`: validation → 400
- `U`: unauthorized → 401
- `F`: forbidden → 403
- `N`: not found → 404
- `C`: conflict → 409
- `I`: internal → 500

## Rules

- Use explicit sentinel errors.
- Do not create ad-hoc errors in handlers.
- Do not dynamically build error codes at runtime.
- Preserve `domain.Error` or the project equivalent as the single error contract.
- Use wrapping only through the existing enriched error mechanism, for example `ErrXxx.Wrap(err)` if available.
- Do not leak infrastructure details to clients.
- New errors must be registered in the project error catalog and `docs/errors/aaa-map.md`.
- If `docs/errors/aaa-map.md` does not exist, create it.
- Follow the project language convention for `message` and `user_message`. Prefer Spanish emphasis when the project is Spanish-first.

## REST contract

REST responses should expose stable error information such as:

```json
{
  "code": "2x365001I",
  "message": "No se pudo procesar la operación.",
  "user_message": "No se pudo procesar la operación.",
  "request_id": "req-..."
}
```

Use the existing envelope if the project already defines one.
