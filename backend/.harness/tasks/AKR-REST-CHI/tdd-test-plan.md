# TDD Test Plan

## Scope

Definir la migración del router REST a Chi sin modificar el contrato OpenAPI ni
los boundaries de seguridad actuales.

## Tests to add/update

1. El router expone exactamente los 23 pares método/path implementados.
2. Setup/login y callbacks GitHub permanecen públicos.
3. Session GET/DELETE requieren sesión; DELETE además requiere Origin válido.
4. Todas las rutas GitHub/Dokploy privadas requieren middleware administrador;
   sus mutaciones validan Origin después de autenticar.
5. `account_id` y `server_id` llegan a los handlers mediante
   `Request.PathValue`.
6. HEAD alcanza las rutas GET mediante `middleware.GetHead`.
7. Paths desconocidos responden 404; métodos inválidos responden 405 con
   `Allow`; trailing slashes no se redirigen.
8. Un `X-Request-ID` válido se conserva; uno inválido o ausente produce un ID
   seguro desde contexto Chi o fallback UUID.
9. Un panic se transforma en el envelope JSON 500 estable sin exponer su valor;
   `http.ErrAbortHandler` conserva su semántica.
10. `Config` inválida, middleware administrador nil o que produce nil fallan
    cerrados sin montar el router.
11. `handler.NewHandlers` construye Auth/GitHub/Dokploy desde un `in.UseCases`
    completo y falla con `ErrInvalidHandlersConfiguration` cuando falta el
    agregado, cualquier caso de uso requerido o paginación válida.
12. `router.New` recibe handlers preconstruidos y falla cerrado si el agregado o
    cualquiera de sus handlers es nil.
13. El bootstrap completa una copia del agregado auth con GitHub/Dokploy, sin
    mutar el valor recibido, y construye los handlers antes del router.

## Expected failing tests before implementation

- El handler actual no implementa `chi.Routes` ni expone un árbol recorrible.
- No existe middleware global RequestID/GetHead ni recuperación JSON de panics.
- El router sigue centralizado en `http.ServeMux`.
- Un middleware administrador que devuelve nil sólo se valida al envolver el
  mux actual y no existe el preflight requerido para grupos Chi.

## Acceptance criteria covered

Inventario HTTP, boundaries de auth/CSRF, path params, HEAD, errores de routing,
request IDs, panic recovery, fail-closed, modularidad y ausencia de cambios de
OpenAPI.

## TDD sequence

1. Agregar tests del middleware de recuperación y resolución de request ID;
   confirmar fallos.
2. Agregar tests de inventario y comportamiento del router Chi; confirmar
   fallos/errores de compilación relevantes.
3. Agregar Chi v5.3.2 e implementar el mínimo para poner esos tests en verde.
4. Ejecutar regresión completa y reviews.

## Revision sequence: handlers aggregation

1. Reabrir `AKR-REST-CHI` y agregar tests del factory y del nuevo contrato del
   router; confirmar fallos antes de crear los agregados.
2. Crear `in.UseCases`, `handler.Handlers` y `handler.NewHandlers`.
3. Mover el wiring de auth/main y GitHub/Dokploy/bootstrap al agregado; adaptar
   `router.Config` para consumir handlers ya construidos.
4. Reejecutar toda la regresión y actualizar los reviews/summaries.

## Final validations

- `go fmt ./...`
- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `.harness/kernel/scripts/check-backend-architecture.sh`
- `.harness/kernel/scripts/check-openapi.sh`
- `.harness/kernel/scripts/check-security.sh`

## Open questions / human approval notes

Sin preguntas abiertas. El usuario aprobó explícitamente este plan en la
conversación del 2026-08-22 al solicitar su implementación.

La revisión para separar construcción de handlers del router también fue
aprobada explícitamente por el usuario el 2026-08-22.
