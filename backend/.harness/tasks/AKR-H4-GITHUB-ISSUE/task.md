# AKR-H4-GITHUB-ISSUE — Hito 4 GitHub Issue

## Estado

in_progress — plan TDD aprobado explícitamente el 2026-08-23

## Tipo de tarea

backend-api-feature

## Modo de proyecto

existing_project

## Contexto

H1–H3 están integrados y el pipeline productivo termina actualmente con una
Investigation completada y el Incident aún en `investigating`. El dominio y
OpenAPI anticipan publicación de GitHub Issues, pero no existen todavía el
publisher, la persistencia normalizada ni la proyección REST completa.

## Objetivo

Implementar AKR-45, AKR-46, AKR-47 y el soporte backend de AKR-48:
`Investigation completed → GitHub Issue → IssueReference → REST`.

## Criterios de aceptación

- Toda Investigation completada publica una Issue para `fixable` y
  `requires_human`.
- La Issue usa el Project/repository configurado y credenciales exclusivamente
  desde el Credential Store.
- IssueReference queda ligada a Incident e Investigation con idempotencia
  durable por Investigation.
- El contenido separa Evidence observada de conclusiones QVAC y no filtra
  secretos.
- Incident detail y timeline cumplen OpenAPI 1.6.0.
- H1–H3 y todos los gates continúan verdes.

## Restricciones técnicas

- Profile `backend_api`; workflow `backend_api_feature`.
- Arquitectura hexagonal, SRP, PostgreSQL/GORM/gormigrate y OpenAPI canónico.
- Sin transacciones abiertas durante GitHub/QVAC.
- Sin Remediation, branches, commits, Pull Requests, merge o deploy.

## Zonas afectadas

Dominio/ports, Investigation orchestration, servicio de contenido, adapters
GitHub/PostgreSQL, REST Incident/timeline, migraciones, documentación y harness.

## Fuera de alcance

AKR-49+, H5, reconciliación robusta de side effects remotos de H6 y frontend.

## Instrucción para el harness

`implementation-brief.md` y `tdd-test-plan.md` fueron aprobados por el usuario
al pedir explícitamente implementar el plan presentado.
