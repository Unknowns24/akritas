# Implementation Brief

## Task

`AKR-OPENAPI-MVP`: contrato OpenAPI v1 completo para backend/frontend, autenticacion single-admin con TOTP y actualizacion de la documentacion autoritativa.

## Current project context

- El backend es un proyecto nuevo: no existen `go.mod`, handlers, modelos ni OpenAPI previo.
- El profile activo es `backend_api`; la ubicacion canonica preferida es `docs/openapi.yaml`.
- `../docs/` contiene spec, MVP, backlog, dominio, arquitecturas y ADR aceptados.
- ADR-001 a ADR-007 fijan QVAC local, deteccion deterministica, Issue obligatoria, frontera humana en PR, Credential Store y MonitoringConfiguration.
- El ZIP de Stitch aporta nueve pantallas. Sus conceptos no documentados o contradictorios no se transforman automaticamente en API.

## Proposed approach

1. Fortalecer el gate OpenAPI con validaciones estructurales deterministas y sin dependencia de red.
2. Crear un contrato OpenAPI 3.1.0 modularizado mediante `components`, pero conservado en un unico YAML canonico para facilitar consumo y validacion.
3. Organizar operaciones por tags: Auth, System, Dashboard, GitHub, Dokploy, QVAC, Projects, Automation, Incidents, Investigations, Remediations, PullRequests y Operations.
4. Usar envelopes consistentes, paginacion por cursor y operaciones asincronas consultables.
5. Actualizar documentacion y agregar ADR aceptados para autenticacion TOTP y autenticacion GitHub PAT/App Manifest.
6. Ejecutar gates, revisar arquitectura/seguridad y cerrar los artefactos del harness.

## Architecture impact

El contrato define los futuros input ports y boundaries REST, pero no implementa capas. La autenticacion se modela como capacidad transversal del backend con middleware de sesion, sin introducir usuarios empresariales ni RBAC. Los secretos siguen perteneciendo a infraestructura y nunca a DTOs de salida.

Los procesos de investigacion, Issue y remediacion se representan como operaciones asincronas, compatibles con workers retry-safe y sin mantener requests HTTP abiertas.

## API/OpenAPI impact

- OpenAPI `3.1.0`, `info.version: 1.0.0`, server `/api/v1`.
- Cookie de sesion opaca como `cookieAuth`; callbacks GitHub protegidos por `state` de un solo uso.
- Envelopes `data`, colecciones con `paging`, error estable y `Operation` para `202`.
- Recursos publicos para auth, integraciones, control plane, monitoring, incident lifecycle y remediation.
- `Idempotency-Key` requerido en comandos manuales con efectos externos.
- Campos sensibles presentes solo en request schemas `writeOnly` o mantenidos completamente fuera del contrato.

## Data/persistence impact

No se crean migraciones ni modelos persistentes. El contrato anticipa identidades UUID, estados y timestamps que la futura persistencia debera soportar. La documentacion incorporara Administrator, enrollment, session y separacion entre GitHub account type y authentication method.

## Error handling impact

Todas las operaciones reutilizaran un envelope con `code`, `message`, `user_message`, `request_id` y `details` seguros. Se documentaran 400, 401, 403, 404, 409, 429 y 500 donde correspondan. Los tests de conexion devolveran resultados normalizados sin filtrar payloads de proveedores.

## Test strategy

- Validacion estatica YAML/OpenAPI y resolucion de referencias.
- Unicidad de operaciones, consistencia de parametros y cobertura de seguridad.
- Analisis de schemas para prevenir secretos en respuestas.
- Matriz de endpoints/schemas contra PB-001 a PB-055 y pantallas compatibles.
- Validacion de enums, defaults, invariantes, paginacion, idempotencia y operaciones asincronas.
- Ejecucion de gates OpenAPI/security y reviews del harness.

## Risks

- El contrato es amplio y fundacional; inconsistencias ahora se propagarian al backend y frontend.
- GitHub App Manifest combina redirects, callbacks y secretos generados; debe distinguir registro de instalacion.
- Cookies `SameSite=Lax` requieren validacion estricta de Origin/state en futuras implementaciones.
- Exponer evidence, logs o diffs sin redaccion podria filtrar secretos; el contrato debe expresar contenido sanitizado.
- La documentacion principal esta en `../docs/`, fuera del root de escritura actualmente declarado; puede requerir autorizacion adicional al aplicar esos cambios.

## Files likely to change

- `docs/openapi.yaml`.
- `.harness/kernel/scripts/check-openapi.sh`.
- `../docs/` para ADR, producto, dominio, arquitectura, integraciones, demo, diseño y configuracion.
- `.harness/tasks/AKR-OPENAPI-MVP/` para artefactos y reviews.
