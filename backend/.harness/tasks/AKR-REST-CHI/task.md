# AKR-REST-CHI - Migración del adaptador REST a Chi

## Estado

complete

## Tipo de tarea

backend-api-feature

## Modo de proyecto

existing_project

## Contexto

El adaptador REST registra actualmente todos los endpoints en un único
`router.go` mediante `http.ServeMux`. La arquitectura ya reserva Chi para el
boundary HTTP y exige modularidad por feature.

## Objetivo

Adoptar Chi v5 en el adaptador REST, separar el registro de rutas de auth,
GitHub y Dokploy, e incorporar un stack global seguro sin modificar el contrato
HTTP público.

Como revisión de wiring, centralizar la construcción de handlers en
`rest/handler/handlers.go` a partir de un agregado `in.UseCases`, dejando al
router exclusivamente responsable de registrar handlers ya construidos.

## Requerimiento funcional

- Mantener los 23 pares método/path implementados.
- Conservar callbacks públicos, autenticación de sesión y protección Origin.
- Preservar `HEAD` sobre rutas `GET` y `Request.PathValue`.
- Recuperar panics como el error JSON interno estable.
- Mantener construcción fail-closed ante configuración inválida.

## Criterios de aceptación

- `router.New(Config) (http.Handler, error)` conserva su firma; `Config` recibe
  `*handler.Handlers` en lugar de casos de uso, paginación y un auth handler
  separado.
- `handler.NewHandlers` construye y valida Auth, GitHub y Dokploy desde
  `*in.UseCases`.
- Chi v5.3.2 queda como dependencia directa exclusiva del adaptador REST.
- Los tests prueban inventario, boundaries, path params, HEAD, 404/405,
  request IDs, panic recovery y fail-closed.
- OpenAPI permanece sin cambios.
- Tests, race, vet y gates del profile `backend_api` pasan.

## Restricciones técnicas

- Profile `backend_api` y workflow `backend-api-feature`.
- No usar Logger, RealIP, redirects de trailing slash, timeout o CORS global.
- No exponer valores de panic ni secretos.
- No introducir dependencias Chi en core o usecases.

## Archivos o zonas probablemente afectadas

- `internal/adapter/rest/router/`
- `internal/adapter/rest/handler/`
- `internal/core/ports/in/`
- `internal/bootstrap/integrations/` y `cmd/main.go`
- `internal/adapter/rest/middleware/`
- `internal/adapter/rest/response/`
- `go.mod` y `go.sum`
- `.harness/tasks/AKR-REST-CHI/`

## Fuera de alcance

- Nuevos endpoints o cambios de OpenAPI.
- Cambios de DTOs, comportamiento de casos de uso, persistencia o configuración
  runtime.
- Logging HTTP, resolución de IP detrás de proxies y CORS.

## Instrucción para el harness

Primero generar implementation-brief.md y tdd-test-plan.md. No implementar
código hasta aprobación humana.

## Aprobación humana

El usuario aprobó explícitamente el plan TDD en la conversación del 2026-08-22
con la instrucción “PLEASE IMPLEMENT THIS PLAN”.

La revisión de agregación de handlers/usecases fue aprobada explícitamente por
el usuario el mismo día con una segunda instrucción “PLEASE IMPLEMENT THIS
PLAN”.
