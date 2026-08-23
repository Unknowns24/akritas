# Implementation Summary

## Task

AKR-REST-CHI — Migración modular del adaptador REST a Chi v5.3.2.

## Implemented changes

- Se reemplazó `http.ServeMux` por `chi.Router` bajo `/api/v1`.
- El registro quedó separado en rutas auth, GitHub y Dokploy.
- Se conservaron exactamente los 23 pares método/path existentes.
- El stack global aplica RequestID, recuperación JSON y GetHead.
- Los request IDs se validan desde header/contexto Chi antes del fallback UUID.
- El router falla cerrado si el middleware administrador es nil o produce nil.
- Chi quedó como dependencia directa del módulo.
- Se agregó `in.UseCases` como contenedor de wiring de los diez input ports
  usados por REST.
- Se agregó `handler.NewHandlers`, que valida el agregado y construye Auth,
  GitHub y Dokploy en un único lugar.
- `router.Config` recibe `*handler.Handlers`; el router dejó de importar o
  invocar constructores de handlers.
- `cmd/main.go` reúne los casos de auth y el bootstrap completa una copia con
  GitHub/Dokploy antes de construir handlers y router.

## Contract and persistence

No hubo cambios en OpenAPI, DTOs, usecases, bootstrap, configuración,
persistencia ni migraciones.

## Tests added

- Inventario completo del árbol Chi.
- Boundaries públicos, sesión y Origin.
- `Request.PathValue`, HEAD, 404/405 y matching exacto.
- Request IDs válidos/inválidos/ausentes.
- Panic recovery estable y preservación de `http.ErrAbortHandler`.
- Configuración fail-closed del middleware administrador.
- Factory completo y rechazo individual de cada caso de uso ausente,
  paginación inválida y agregado de handlers incompleto.
- Bootstrap fail-closed cuando no recibe `UseCases`.

## TDD result

La fase roja falló por dependencia Chi ausente, router no-Chi y
`RecoverPanics` inexistente. Luego de la implementación, los tests target y la
suite completa quedaron verdes.
