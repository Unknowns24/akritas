# Implementation brief — AKR-H2-DETECTION-INCIDENTS

## Línea base

H1 ya proporciona Projects, MonitoringConfiguration, asociaciones verificadas,
Credential Store, adapter Dokploy con transporte SSRF-safe, PostgreSQL/GORM,
transacciones de aplicación, REST/Chi, Uker y el contrato OpenAPI 1.4.0.

## Estrategia

El incremento extiende esas capacidades sin crear integraciones paralelas:

```text
Project habilitado -> LogSource Dokploy -> checkpoint/overlap
  -> multilinea -> ignored -> built-in/custom -> normalización/fingerprint
  -> contexto durable -> LogEvent -> Incident -> REST
```

El fetch y el cálculo de drafts son determinísticos. La persistencia de
LogEvents, agrupación de Incidents, estado pendiente y avance de cursor ocurre
en una sola transacción con lock de Project y checkpoint. Una clave de
ocurrencia estable hace idempotentes los reintentos.

## Contrato

OpenAPI pasa a 1.5.0. `CreateProjectRequest` agrega la opción one-shot
`initial_log_ingestion` (`from_now | last_10000`, default `from_now`). Se
implementan los endpoints publicados de Incidents: lista, detalle y LogEvents.
Los campos pertenecientes a milestones posteriores permanecen ausentes.

## Persistencia

Se agregan migraciones inmutables para `monitoring_checkpoints`, `incidents` y
`log_events`. Checkpoints operativos eliminan en cascade; la historia de
Incidents/LogEvents restringe el borrado de Project. No se usa AutoMigrate.

## Seguridad y determinismo

Las credenciales siguen resolviéndose por el Credential Store. Los secretos se
redactan antes de persistir evidencia o estado. QVAC/LLMs no intervienen en la
detección. Regex, normalización, prioridad de reglas, hash y agrupación tienen
orden explícito y testeable.
