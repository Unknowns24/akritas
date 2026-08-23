# TDD Test Plan — AKR-H3-FINAL-INTEGRATION

## Scope

AKR-35..AKR-44 sobre el estado combinado H2/H3, incluyendo compilación, lifecycle, Evidence real, QVAC local, tools GitHub read-only, salida estructurada, persistencia, REST y migraciones.

## Tests to add/update

- Catálogo PostgreSQL único con errores H1/H2/H3; unicidad global y mapa AAA exacto.
- Composición REST con Incident/Investigation/Operation/Evidence, rutas H2+H3 y monitoring wiring preservado.
- Incident PostgreSQL real: start, not-found, idempotencia, una Investigation activa, Project y repo resueltos desde Incident.
- LogEvents generan deployment/log_excerpt/stack_trace sólo con datos reales; before/after/source/rules/redaction/bounds; Evidence persiste ante fallo QVAC y es recuperable.
- Runner recibe contenido real; prompt Evidence ≤24 KiB y se reduce por contexto; framing untrusted y secretos ausentes.
- Cinco tools read-only, límites 8/24/8KiB/16KiB, deduplicación, path safety, unknown fail-closed y aislamiento repo A/B.
- Structured output: combinaciones válidas, unknown completa, confidence 0/1, JSON/campos/enums/arrays inválidos y evidence IDs desconocidos.
- Lifecycle y recovery: transiciones atómicas, technical failure, requeue pending, fail running e Incident permanece investigating al completar H3.
- Migraciones/FKs RESTRICT/índice activo/rollback y escenario Testcontainers Project→Incident→LogEvent→Investigation→Evidence→QVAC/tools→resultado.
- Regresiones H1/H2 existentes.

## Expected failing tests before implementation

- El paquete no compila por `Catalog()` mal fusionado.
- Corregido temporalmente el catálogo, la composición falla porque `incidentrepo.Repository` no implementa `IncidentReader.Exists`.
- El catálogo agregado detecta `ErrIncidentNotFound` duplicado.
- Evidence/QVAC/tools/persistence no cumplen todavía los contratos anteriores.

## Acceptance criteria covered

AKR-35, AKR-36, AKR-37, AKR-38, AKR-39, AKR-40, AKR-41, AKR-42, AKR-43 y AKR-44.

## Human approval

Aprobado explícitamente por el usuario el 2026-08-23, incluyendo: corpus 128 KiB/25, Evidence inicial 24 KiB máximo adaptativo a `ctx_size=16384`, y prohibición de llamar `StartIssuePublication` desde H3.
