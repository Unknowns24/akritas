# Implementation Summary

## Task

`AKR-DOKPLOY-COMPOSE-SOURCES` — descubrimiento y monitoreo de aplicaciones y
servicios Compose de Dokploy.

## Implemented

- Se agregó el modelo discriminado `DokploySource` con variantes `application`
  y `compose_service`, identidad estable y validación cruzada de servicio y
  runtime.
- Project create/update y sus respuestas reemplazan inmediatamente el contrato
  anterior por `dokploy_source`; la activación de monitoring vuelve a resolver
  el snapshot remoto.
- Se incorporaron `DokployCompose`, `DokployComposeService`, `compose_count` y
  los casos de uso para listar Composes y servicios.
- Se agregaron los endpoints paginados de Composes y el endpoint de servicios
  con `refresh=false` por caché y `refresh=true` por fetch explícito.
- El adapter integra `compose.search`, `compose.one`, `compose.loadServices`,
  `docker.getContainersByAppNameMatch` y `compose.readLogs`.
- La lectura Compose resuelve el contenedor en cada ciclo, filtra labels de
  docker-compose/stack, considera sólo candidatos running y selecciona el ID
  lexicográficamente menor.
- Se generalizaron checkpoints, LogEvents, monitoreo y evidencia nueva a la
  identidad de fuente. La evidencia histórica no se modifica.
- Se agregó una migración reversible que transforma las fuentes existentes en
  `application`, conserva sus relaciones, aplica unicidad por servicio y evita
  un rollback destructivo cuando existen Projects Compose.
- OpenAPI se actualizó a API `2.0.0` manteniendo las rutas `/api/v1`, y se
  documentaron errores y ADR-015.

## Tests added or updated

- Dominio: discriminante, combinaciones inválidas e identidad por servicio.
- Adapter Dokploy: paginación/mapeo Compose, cache/fetch de servicios,
  deduplicación, selección de réplica y lectura application/Compose.
- Projects y repositorios: resolución, conflictos y coexistencia de servicios
  distintos del mismo Compose.
- REST/router: inventario de rutas, envelope de servicios y booleano `refresh`
  estricto, incluido el valor vacío.
- Monitoreo/evidencia: propagación de identidad genérica.
- Registro de migraciones y compilación de fixtures con tag `integration`.

## Validation

- `go test ./...`: pass.
- `go vet ./...`: pass.
- `go fmt ./...`: pass.
- `go test -tags integration -run '^$' ./internal/adapter/db/postgres/...`: pass.
- `.harness/kernel/scripts/check-backend-architecture.sh`: pass.
- `.harness/kernel/scripts/check-openapi.sh`: pass — 62 operaciones, 123 schemas.
- `.harness/kernel/scripts/check-security.sh`: pass.
- `git diff --check`: pass.

## Notes

- Los tests PostgreSQL con tag `integration` se compilaron, pero no se ejecutó
  una base PostgreSQL real en esta sesión.
- Los archivos `.DS_Store` no relacionados permanecen intactos.
