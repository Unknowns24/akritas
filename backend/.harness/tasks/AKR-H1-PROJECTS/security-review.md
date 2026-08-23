# Security review — AKR-H1-PROJECTS

## Veredicto

Sin hallazgos abiertos.

## Controles verificados

- Project y sus DTOs sólo contienen referencias y snapshots no secretos.
- PATs, installation tokens y API keys se obtienen del Credential Store y se
  limpian después de usarse; nunca se persisten ni retornan con Project.
- Errores de proveedor se normalizan sin exponer bodies o credenciales.
- Resolución Dokploy conserva la validación SSRF/DNS del gateway vigente.
- Todas las rutas Project requieren sesión de Administrator.
- POST, PATCH, PUT y DELETE exigen un Origin exacto allowlisted para CSRF.
- Filtros y sorts están allowlisted antes de llegar a Uker/GORM.
- Updates usan optimistic concurrency y los fallos externos ocurren antes de
  cualquier escritura.
- FKs restrict y consultas de uso protegen integraciones referenciadas.

El gate `.harness/kernel/scripts/check-security.sh` finalizó correctamente.
