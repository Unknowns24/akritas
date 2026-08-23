# Architecture Review

## Summary

La implementación respeta el profile `backend_api`, el `tdd-test-plan.md` aprobado y la dirección de dependencias hexagonal: `adapter/rest|db → usecase → core`. Es la primera vez que el proyecto instancia ports, usecases, un adapter de persistencia y HTTP; el orden de wiring (`config → DB+migraciones → repositorios → adapters → usecases → router`) sigue `wiring.md` exactamente.

## Layering

- `internal/core/domain` y `internal/core/ports/{in,out}` no importan `internal/adapter`, GORM, Chi, `net/http` ni `os/exec` (`check-backend-architecture.sh` pasa).
- `internal/usecase/auth` sólo depende de `internal/core/domain` y `internal/core/ports/{in,out}` — ningún adapter concreto.
- `internal/adapter/db/postgres` y `internal/adapter/security` implementan los out ports sin filtrar tipos GORM ni de terceros hacia el usecase.
- `internal/adapter/rest` traduce en el borde: DTOs propios, nunca se serializan entidades de dominio ni modelos GORM directamente.

## Modularity / SRP

- Un archivo por operación pública en usecases, repositorios y handlers (`get_setup_status.go`, `start_administrator_setup.go`, `exists_active.go`, `save.go`), con `handler.go`/`repo.go`/`uc.go`-equivalentes limitados a struct + constructor.
- Ocho out ports en ocho archivos separados, cada uno con una sola responsabilidad.
- `internal/adapter/rest/response` es una excepción deliberada: agrupa helpers de escritura JSON/error compartidos por ambos handlers (no es una "operación de negocio" en el sentido de la política, sino la implementación directa de la regla ya documentada en `aaa-map.md` — "los adapters deben mapear el tipo final del código a HTTP").

## OpenAPI consistency

- No se modificó `docs/openapi.yaml`; `check-openapi.sh` confirma 59 operaciones y 112 schemas sin cambios.
- Los DTOs (`SetupStatus`, `SetupStatusResponse`, `SetupRequest`, `TotpEnrollment`, `TotpEnrollmentResponse`, `ErrorResponse`) espejan los schemas campo por campo, incluyendo el header `Cache-Control: no-store` inline tal como está en el spec (no el `$ref` que usan los endpoints hermanos).
- Verificado manualmente contra el servidor real: los cuatro status codes documentados (200, 201, 400, 409, 429) se comportan como el contrato especifica.

## Findings

- `internal/adapter/rest/handler/auth` importa `internal/usecase/auth` directamente (no sólo `internal/core/ports/in`) para reconocer `authusecase.ErrSetupRateLimited` — el rate limit no tiene representación en `ports/in` porque el `ErrorCode` del OpenAPI no tiene letra para 429. Es una dependencia adapter→usecase concreto en vez de adapter→puerto, permitida por `project-structure.md` (`adapter/rest → usecase → core`) pero menos limpia que el resto del wiring. No bloqueante; si una futura tarea de rate limiting (PB-065) formaliza este caso, considerar mover el sentinel a `ports/in` o modelarlo como un tipo de error dedicado en el puerto.
- `internal/adapter/db/postgres/dbtest` es un paquete no-test que importa `testing` (patrón estándar para helpers de test compartidos entre paquetes de repositorio) — señalado por si el revisor prefiere un nombre `_test` interno en su lugar; no viola ninguna regla del profile.

## Result

pass
