# AKR-BACKEND-FOUNDATION - Fundación del backend y dominio MVP

## Estado

complete

## Tipo de tarea

backend-api-feature

## Modo de proyecto

new_project

## Contexto

El backend contiene el contrato OpenAPI y la documentación autoritativa del producto, pero todavía no posee un módulo Go ni código de aplicación. El equipo necesita una base común y estable que permita desarrollar features en ramas paralelas sin inventar estructuras, estados ni reglas de dominio incompatibles.

## Objetivo

Crear un módulo Go 1.26 compilable, establecer la estructura hexagonal definida por el profile `backend_api` y modelar en un único paquete `internal/core/domain` los conceptos e invariantes documentados para el MVP.

## Requerimiento funcional

- Crear `github.com/Unknowns24/akritas/backend` y un entrypoint deliberadamente inerte.
- Versionar el esqueleto de las capas futuras sin inventar puertos, DTOs, repositorios o adapters concretos.
- Modelar auth, integraciones, proyectos, monitoring, detección, incidentes, investigaciones, evidencia, remediación, validaciones y referencias externas.
- Usar UUID para identidades internas, tipos nativos de tiempo y colecciones defensivas.
- Implementar constructores, validación e invariantes/transiciones expresamente documentadas.
- Crear el contrato enriquecido de errores de dominio y su catálogo.

## Criterios de aceptación

- `go test ./...`, `go test -race ./...` y `go vet ./...` finalizan correctamente.
- Los defaults y validaciones de `MonitoringConfiguration` coinciden con ADR-007 y OpenAPI.
- Los ciclos de Incident, Investigation, Remediation y ValidationResult rechazan transiciones inválidas.
- Un Incident terminal no agrupa nuevas ocurrencias y la ventana usa `last_seen_at`.
- La publicación de Issue es obligatoria antes de cerrar o remediar; una PR requiere remediación validada.
- Ninguna entidad contiene tags de transporte/persistencia ni secretos de autenticación o integración.
- Los errores cumplen `DxAAABBBT`, son estables y están registrados en `docs/errors/aaa-map.md`.
- Los checks de arquitectura, OpenAPI y seguridad del harness pasan sin modificar el contrato HTTP.

## Restricciones técnicas

- Profile: `.harness/kernel/profiles/go-hexagonal-api.yaml`.
- Workflow: `.harness/kernel/workflows/backend-api-feature.yaml`.
- Dependencias permitidas en dominio: biblioteca estándar y `github.com/google/uuid`.
- No importar HTTP, Chi, GORM, adapters, SDKs externos, filesystem ni ejecución de procesos desde `internal/core`.
- Mantener un solo package Go `domain`, con un archivo por concepto o responsabilidad cohesiva.
- No implementar código antes de la aprobación humana de `tdd-test-plan.md`.

## Archivos o zonas probablemente afectadas

- `go.mod`, `go.sum`, `cmd/main.go`.
- `internal/core/domain/`.
- Directorios base bajo `config/`, `internal/core/ports/`, `internal/usecase/`, `internal/service/` e `internal/adapter/`.
- `docs/errors/aaa-map.md`.
- `.harness/tasks/AKR-BACKEND-FOUNDATION/` y `.harness/tasks/index.md`.

## Fuera de alcance

- Servidor HTTP, handlers, middleware, routers y DTOs.
- Puertos de entrada/salida, casos de uso y services concretos.
- Persistencia, tecnología de base de datos y migraciones.
- Clientes GitHub, Dokploy, QVAC, Git o filesystem.
- OpenAPI, overview, activity, timeline, system status, paginación y operaciones asíncronas.

## Instrucción para el harness

Primero generar `implementation-brief.md` y `tdd-test-plan.md`. No implementar código hasta aprobación humana.
