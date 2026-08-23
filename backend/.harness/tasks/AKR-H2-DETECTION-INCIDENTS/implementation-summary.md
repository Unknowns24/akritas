# Implementation summary — AKR-H2-DETECTION-INCIDENTS

## Capacidades implementadas

- `LogSource` sobre el cliente Dokploy existente para `application.readLogs`,
  con `tail=10000`, `since` incremental y backfill explícito `all`;
- parser Docker RFC3339Nano, ordinal por timestamp y hash de contenido;
- checkpoint durable por epoch Project/source, overlap de un segundo,
  deduplicación y discontinuidad saturada fail-closed;
- multilinea determinística, precedencia de ignored patterns, siete built-ins,
  regex custom, redacción, normalización conservadora y SHA-256 versionado;
- contexto before/after acotado y ocurrencias pendientes durables con timeout
  de 30 segundos;
- creación idempotente de LogEvent y grouping transaccional de Incident contra
  `last_seen_at`, incluida promoción de severidad;
- endpoints REST de lista/detalle de Incident y LogEvents publicados;
- runner background cancelable, polling configurable y concurrencia máxima 4;
- `initial_log_ingestion` one-shot en Create Project y OpenAPI 1.5.0.

## Componentes principales

- `internal/service/detection`
- `internal/service/monitoring`
- `internal/adapter/external/dokploy/read_logs.go`
- repositorios PostgreSQL `monitoring` e `incident`
- usecase, DTOs, mappers, handlers y rutas de Incident
- bootstrap runtime compartido y shutdown ordenado en `cmd/main.go`

## Persistencia

- `20260823_01_add_monitoring_checkpoints`
- `20260823_02_add_incidents`
- `20260823_03_add_log_events`

Cursor, efectos de Incident/LogEvent y estado pendiente se actualizan dentro de
la misma transacción. El lock de Project define un orden global antes del lock
de checkpoint/Incident; la occurrence key única hace idempotentes los retries.
