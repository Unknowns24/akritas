# AKR-OPENAPI-MVP - Contrato OpenAPI v1 y autenticacion TOTP

## Estado

complete

## Tipo de tarea

backend-api-feature

## Modo de proyecto

new_project

## Contexto

Akritas todavia no posee implementacion backend ni un contrato HTTP. La documentacion de producto y arquitectura vive principalmente en `../docs/`, y el ZIP de Stitch adjunto define las pantallas previstas para overview, proyectos, incidentes e integraciones.

El contrato OpenAPI debe permitir que el frontend comience a desarrollarse sin inventar endpoints o campos. La documentacion aceptada y los ADR tienen precedencia sobre conceptos exclusivos de la maqueta.

## Objetivo

Crear el contrato publico OpenAPI 3.1.0 version `1.0.0` para todos los hitos P0 del MVP y las proyecciones compatibles requeridas por la interfaz, incorporando autenticacion de administrador unico mediante bootstrap, password, TOTP y sesion segura.

## Requerimiento funcional

- Definir autenticacion, sistema/diagnosticos, overview, actividad y operaciones asincronas.
- Definir integraciones GitHub por PAT y GitHub App Manifest, Dokploy y QVAC.
- Definir Projects, MonitoringConfiguration y AutomationPolicy.
- Definir Incidents, LogEvents, Investigations, Evidence, Remediation, ValidationResults y referencias GitHub.
- Definir paginacion por cursor, errores seguros, idempotencia, security schemes y ejemplos.
- Actualizar la documentacion de producto, dominio, arquitectura, seguridad, configuracion y UX afectada.

## Criterios de aceptacion

- Existe un unico contrato canonico en `docs/openapi.yaml` y declara OpenAPI 3.1.0 / API 1.0.0.
- Todas las operaciones tienen `operationId` unico, schemas reutilizables y respuestas documentadas.
- Solo health/readiness, setup/login/recovery y callbacks GitHub son publicos.
- Ningun schema de respuesta expone passwords, sesiones, bootstrap tokens, secretos TOTP ni credenciales de GitHub, Dokploy o QVAC.
- Los endpoints y schemas cubren PB-001 a PB-055 y las pantallas compatibles del ZIP.
- MonitoringConfiguration respeta ADR-007 y AutomationPolicy mantiene sus invariantes.
- El lifecycle termina en PR o `requires_human` sin merge, deploy ni falso estado de resolucion productiva.
- La documentacion describe el bootstrap TOTP, las sesiones y ambos metodos de conexion GitHub.
- Los gates de OpenAPI y seguridad finalizan correctamente.

## Restricciones tecnicas

- Profile `backend_api` y workflow `backend-api-feature`.
- Arquitectura hexagonal y OpenAPI como fuente de verdad.
- JSON publico en `snake_case`; timestamps RFC 3339 UTC; UUID; duraciones ISO 8601.
- Paginacion firmada por cursor para colecciones operativas.
- Errores compatibles con el codigo de dominio `DxAAABBBT`.
- Inferencia QVAC exclusivamente local/privada y sin secretos.
- La creacion de GitHub Issue es obligatoria despues de una Investigation completada.
- No implementar handlers, persistencia, workers ni adapters en esta tarea.

## Archivos o zonas probablemente afectadas

- `docs/openapi.yaml` y el gate `.harness/kernel/scripts/check-openapi.sh`.
- Documentacion de producto/arquitectura en `../docs/`.
- `.harness/tasks/AKR-OPENAPI-MVP/`.

## Fuera de alcance

- Implementacion Go del contrato.
- Frontend Next.js.
- Team, invitaciones, RBAC, multi-tenancy o passkeys WebAuthn.
- Incidentes creados manualmente.
- GitHub Enterprise, merge, deploy o rollback automaticos.

## Instruccion para el harness

Primero generar `implementation-brief.md` y `tdd-test-plan.md`. No implementar el contrato ni modificar documentacion de producto hasta aprobacion humana explicita de este plan TDD.
