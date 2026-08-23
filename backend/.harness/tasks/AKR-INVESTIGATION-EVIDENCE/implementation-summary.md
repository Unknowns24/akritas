# Implementation summary — AKR-INVESTIGATION-EVIDENCE

## Resultado

Se implementó `GET /investigations/{investigation_id}/evidence` exactamente
como lo define `docs/openapi.yaml`, con persistencia real de Evidence y
ensamblado real de `deployment_metadata` dentro del pipeline asíncrono de
PB-026, sin tocar la lógica ya existente de `Investigation`/`Operation`.

## Capacidad incorporada

- `GET /investigations/{investigation_id}/evidence` — listado paginado
  (Uker), filtro opcional `type_in`, 404 si la investigación no existe.
- Persistencia real de Evidence (tabla `evidence`, FK real a
  `investigations`, write-once).
- `out.EvidenceAssembler` invocado dentro de `RunUseCase.Execute`, justo
  después de persistir `Start()` y antes de `InvestigationRunner.Run`: cada
  Evidence ensamblada se persiste y `Investigation.EvidenceCount` se
  actualiza con la cantidad real antes de invocar al runner.
- La evidencia ensamblada y el `EvidenceCount` sobreviven aunque la
  investigación termine en `failed` (verificado con test dedicado).
- `out.IncidentReader` extendido con `Get` (además del `Exists` de PB-026),
  reutilizado por el nuevo `EvidenceAssembler`.

## Qué evidencia quedó real vs. pendiente

- `deployment_metadata`: real, ensamblada a partir de `Project` (nombre,
  monitoring/health status, snapshot de GitHubRepository/DokployApplication)
  vía `Investigation.IncidentID` → `Incident.ProjectID` → `Project`. Nunca
  credenciales.
- `log_excerpt`, `stack_trace`: sin implementar — dependen de H2
  (LogEvent/Incident real).
- `code_location`, `commit`, `diff`: sin implementar — dependen de
  PB-030/PB-031 (lectura de repositorio/commits).
- `validation_result`: sin implementar — pertenece a Remediation (H5).

No se generaron placeholders para ninguno de estos cuatro tipos.

## Contrato y persistencia

No se modificó `docs/openapi.yaml`: el path y los schemas ya estaban
publicados. Migración `20260822_12_add_evidence` (SQL explícito, FK real a
`investigations` con `ON DELETE CASCADE`, checks de `type`/`redacted`/
coherencia de líneas). Error de persistencia nuevo:
`ErrEvidencePersistence` (`0x206001I`).
