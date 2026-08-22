# Implementation Brief

## Task

`AKR-BACKEND-FOUNDATION`: crear la fundación compilable del backend y el dominio completo del MVP.

## Current project context

El directorio backend contiene `docs/openapi.yaml`, el harness y documentación compartida en `../docs/`, pero no contiene `go.mod`, `cmd/` ni `internal/`. El contrato vigente usa UUID, timestamps RFC 3339, duraciones ISO 8601 y estados explícitos. La arquitectura aceptada es hexagonal, con OpenAPI como contrato HTTP y GORM/gormigrate reservados para adapters futuros.

Fuentes autoritativas revisadas:

- `../docs/backend-architecture.md`, `../docs/domain.md`, `../docs/spec.md`, `../docs/mvp.md` y `../docs/incident-lifecycle.md`.
- ADR-001 a ADR-009, especialmente detección, monitoring, auth y seguridad de credenciales.
- `docs/openapi.yaml` y `.harness/memory/{project-summary,decisions}.md`.

## Proposed approach

1. Crear el módulo `github.com/Unknowns24/akritas/backend` con Go 1.26 y `github.com/google/uuid`.
2. Agregar un `cmd/main.go` inerte que solo garantice compilación.
3. Crear placeholders versionados para las capas hexagonales futuras, sin declaraciones Go artificiales.
4. Implementar un package plano `internal/core/domain` con un archivo por entidad/value object y tests espejados.
5. Usar constructores y `Validate`; devolver punteros para entidades mutables y valores para value objects.
6. Proteger slices mediante copias y no agregar tags JSON/GORM.
7. Implementar solamente invariantes y transiciones definidas por documentación/OpenAPI.
8. Incorporar errores enriquecidos con catálogo estable y mensajes públicos seguros en español.

## Architecture impact

Se crea la estructura inicial:

```text
cmd/main.go
config/
internal/core/domain/
internal/core/ports/{in,out,paging}/
internal/usecase/
internal/service/
internal/adapter/db/
internal/adapter/external/{github,dokploy,qvac,git}/
internal/adapter/rest/{dto,handler,middleware,router}/
```

El dominio dependerá únicamente de stdlib y UUID. No se creará bootstrap, wiring, puerto ni adapter concreto hasta que una feature lo necesite.

## API/OpenAPI impact

No hay cambio de contrato. `docs/openapi.yaml` no se modifica. Los nombres de estados e invariantes del dominio se mantienen compatibles con sus schemas, pero las entidades no funcionan como DTOs.

## Data/persistence impact

No se elige motor de base de datos, no se crean modelos GORM ni migraciones. El dominio no contiene tags ni decisiones de schema.

## Error handling impact

Se introduce `domain.Error` con `Error`, `Unwrap`, `Is` y `Wrap`. Se reservan componentes de capa dominio:

- `0x401`: autenticación.
- `0x402`: integraciones.
- `0x403`: Project, monitoring y automation.
- `0x404`: detección e Incident.
- `0x405`: Investigation y Evidence.
- `0x406`: Remediation y Validation.

Los sentinels usarán causas incrementales y tipos `V`, `U` o `C`; cada código quedará registrado en `docs/errors/aaa-map.md`.

## Test strategy

Tests unitarios en el mismo package cubrirán enums, constructores, validaciones, copias defensivas, defaults, sesiones, grouping, ciclos de estado, precondiciones de Issue/PR y semántica de errores. Luego se ejecutarán test normal/race, vet, formato y todos los gates del profile.

## Risks

- Introducir reglas de workflow no documentadas bloquearía features posteriores; los métodos se limitarán a transiciones aprobadas.
- Tratar entidades como DTOs o modelos de persistencia rompería los límites hexagonales.
- Un package plano puede crecer; se mitiga con un archivo y test por concepto, tal como el usuario solicitó.
- La estrategia de reaparición de incidentes completados está fuera del MVP: los incidentes terminales no agruparán nuevas ocurrencias.

## Files likely to change

- `go.mod`, `go.sum`, `cmd/main.go`.
- `internal/core/domain/*.go` y `*_test.go`.
- Placeholders de la estructura hexagonal.
- `docs/errors/aaa-map.md`.
- Artefactos de `.harness/tasks/AKR-BACKEND-FOUNDATION/`.
