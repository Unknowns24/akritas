# AKR-INVESTIGATION-LIFECYCLE

Crear y administrar `Investigation` sobre incidentes ya detectados, con ejecución
asíncrona vía `Operation` como infraestructura genérica reutilizable por H5.
Cubre PB-026; la ejecución real (QVAC) queda fuera de alcance.

## Alcance aprobado

- `GET /incidents/{incident_id}/investigations` — listado paginado por incidente.
- `POST /incidents/{incident_id}/investigations` — crea Investigation `pending`,
  responde 202 + Operation, ejecuta de forma asíncrona con `Idempotency-Key`.
- `GET /investigations/{investigation_id}` — detalle completo.
- `GET /operations/{operation_id}` — polling de estado de Operation.

## Frontera con H2 (obligatoria)

H2 (Detection + Incidents) solo existe como dominio (`incident.go`), sin
repositorio ni usecase. Se define `out.IncidentReader` (mínimo, solo `Exists`)
para validar `incident_id`; en producción se conecta un stub que siempre
deniega (`Exists=false`) hasta que H2 mergee un reader real. Tests usan un fake
controlable.

## Decisiones de producto aprobadas

- `started_at` (required en el contrato) se expone en el DTO como `StartedAt`
  si está seteado, o `CreatedAt` como fallback mientras el estado es `pending`
  — el dominio no cambia.
- `Operation` es dominio + repositorio + polling genéricos, parametrizados por
  `type`/`resource_type`/`resource_id`; no exclusivos de Investigation.
- La ejecución real de la Investigation se resuelve con `out.InvestigationRunner`,
  implementado hoy por un stub que falla la Investigation con un mensaje
  explícito de "pendiente de PB-028+"; el pipeline asíncrono completo (crear →
  encolar → poll) es demostrable end-to-end igual.
