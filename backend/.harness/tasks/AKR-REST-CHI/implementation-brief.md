# Implementation Brief

## Task

AKR-REST-CHI — Migración interna del router REST desde `http.ServeMux` a Chi
v5.3.2.

## Current project context

El router actual concentra 23 endpoints en un único archivo. Auth público,
sesión, callbacks GitHub e integraciones privadas ya tienen boundaries de
seguridad distintos que deben conservarse. Los handlers con parámetros usan
`Request.PathValue`.

## Proposed approach

- Mantener `router.New` como composition point y separar registradores privados
  para auth, GitHub y Dokploy.
- Crear un `chi.Router` raíz con prefijo `/api/v1`.
- Aplicar globalmente RequestID, recuperación JSON y GetHead.
- Validar el middleware administrador antes de registrar rutas.
- Mantener middlewares de sesión y Origin en sus scopes actuales.
- Agregar `in.UseCases` como contenedor de wiring de todos los input ports.
- Construir `handler.Handlers` en el bootstrap mediante `handler.NewHandlers` y
  entregar al router únicamente handlers listos, AuthenticateSession y
  middleware/configuración de routing.

## Architecture impact

Chi permanece dentro de `internal/adapter/rest`. `in.UseCases` no implementa
lógica: agrupa interfaces existentes para composición. El bootstrap continúa
siendo dueño de la construcción y el router deja de construir handlers.

## API/OpenAPI impact

Ninguno. Paths, métodos, payloads, estados y versión OpenAPI permanecen iguales.
El matching de paths será exacto y no habrá redirects de trailing slash.

## Data/persistence impact

Ninguno. No se requieren migraciones.

## Error handling impact

Los panics anteriores a un response committed se convertirán en el envelope
interno existente. El request ID se resolverá desde un header válido, luego el
contexto Chi y finalmente `req-<uuid>`. `http.ErrAbortHandler` no se recuperará.

## Test strategy

Tests TDD del router y middleware al nivel REST, seguidos por regresión completa,
race detector, vet y gates del harness.

## Risks

- Alterar sin intención boundaries públicos/privados u orden sesión/Origin.
- Perder la semántica HEAD de ServeMux.
- Responder panics fuera del envelope común.
- Introducir un request ID no validado desde headers externos.

## Files likely to change

`internal/core/ports/in`, `internal/adapter/rest/handler`,
`internal/adapter/rest/router`, `internal/bootstrap/integrations`, `cmd/main.go`
y artefactos de la tarea.
