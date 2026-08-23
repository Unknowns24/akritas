# Implementation summary — AKR-H1-PROJECTS

## Resultado

La historia de `origin/feat/project-handling` se integró mediante un merge con
estrategia `ours`. PB-009 a PB-013 se reconstruyeron sobre la arquitectura de
`feat/backend-milestone-1`; no se importaron el bootstrap PostgreSQL, router,
sesión, paginación, envelopes, repositorios de integraciones ni configuración de
la rama antigua.

## Capacidad incorporada

- CRUD completo de Project, con borrado físico 204 y conflictos para monitoring
  activo o dependencias referenciales.
- snapshots no secretos de GitHubRepository y DokployApplication obtenidos de
  gateways reales antes de create, cambios de asociación y activación.
- compatibilidad GitHub por identificador opaco o `owner/name`, validación del
  owner configurado y coincidencia exacta de `default_branch`.
- resolución exacta del `application_identifier` de Dokploy.
- nombre único case-insensitive, aplicación Dokploy exclusiva y optimistic
  concurrency por `updated_at`.
- MonitoringConfiguration completa; toda configuración habilitada revalida
  proveedores y transiciona a `starting`.
- listado Uker con cursores firmados, filtros `name_like` y
  `monitoring_status_in`, y sorts allowlisted.
- consultas reales de uso de Projects para impedir eliminar las integraciones
  GitHub/Dokploy referenciadas.

## Contrato y persistencia

OpenAPI se elevó a 1.4.0 y suma `DELETE /projects/{project_id}`. La migración
`20260822_09_add_projects` declara SQL PostgreSQL explícito, FKs restrict,
snapshots, JSONB, checks e índices únicos. Los errores estables y las decisiones
de lifecycle quedaron documentados.
