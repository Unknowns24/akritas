# Implementation summary — AKR-INVESTIGATION-LIFECYCLE

## Resultado

Se implementaron los 4 endpoints exactamente como los define
`docs/openapi.yaml`, reutilizando el dominio `Investigation` existente sin
tocar su lógica, y se agregó `Operation` como infraestructura genérica de
comando asíncrono (dominio + puertos + persistencia + polling), pensada para
ser reutilizada por remediation/pull_request más adelante.

## Capacidad incorporada

- `GET /incidents/{incident_id}/investigations` — listado paginado (Uker) y
  scoped por incidente, 404 si el incidente no existe.
- `POST /incidents/{incident_id}/investigations` — crea Investigation
  `pending` + Operation `queued`, responde 202 con la Operation, exige y
  valida `Idempotency-Key`, rechaza duplicados activos (409) y resuelve
  replays de la misma key devolviendo la Operation ya existente.
- `GET /investigations/{investigation_id}` — detalle completo, `started_at`
  siempre presente (fallback a `created_at` en `pending`).
- `GET /operations/{operation_id}` — polling de estado.
- Pipeline asíncrono real: `InvestigationDispatcher` (goroutine con timeout
  propio, desacoplada del contexto de la request) invoca `RunInvestigation`,
  que transiciona Investigation y Operation de forma discreta
  (pending/queued → running → completed/failed → succeeded/failed),
  persistiendo cada paso.

## Frontera con H2 y PB-028+

- `out.IncidentReader` (solo `Exists`) es el único punto de contacto con H2.
  En producción se conecta `stub.DenyAllIncidentReader`, que siempre
  responde `false`: las 4 rutas están vivas y auditables, pero crear o listar
  investigaciones da 404 hasta que H2 aporte un reader real. Tests usan un
  fake controlable.
- `out.InvestigationRunner` es hoy `qvac.StubRunner`, que siempre falla con
  un mensaje explícito ("QVAC integration is not implemented yet; see
  PB-028+"). El pipeline completo (crear → encolar → poll) es demostrable
  end-to-end aunque el resultado real de la investigación quede pendiente de
  PB-028+.

## Contrato y persistencia

No se modificó `docs/openapi.yaml`: los 4 paths y schemas ya estaban
publicados y se implementaron tal cual. Migraciones
`20260822_10_add_investigations` y `20260822_11_add_operations` (SQL
explícito, JSONB con default `[]`, checks de enum, índice único parcial de
`idempotency_key`). `investigations.incident_id` no tiene FK todavía —gap
documentado, pendiente de una migración de seguimiento cuando H2 mergee.
