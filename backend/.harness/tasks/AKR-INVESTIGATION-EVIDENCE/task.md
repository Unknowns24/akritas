# AKR-INVESTIGATION-EVIDENCE

Construir Evidence real para las Investigation ya implementadas en PB-026:
ensamblado durante la ejecución asíncrona, persistencia y listado paginado.
Cubre PB-027.

## Alcance aprobado

- `GET /investigations/{investigation_id}/evidence` — listado paginado por
  investigación, filtro opcional `type_in`.
- Persistencia real de Evidence (tabla `evidence`, con FK a `investigations`
  — a diferencia de `investigations.incident_id`, esta FK sí existe porque
  `investigations` ya está en este repo).
- `out.EvidenceAssembler` invocado dentro de `RunUseCase.Execute`, después de
  `Start()` y antes de `InvestigationRunner.Run`: cada Evidence ensamblada se
  persiste y `Investigation.EvidenceCount` se actualiza con la cantidad real.
  La evidencia ensamblada se conserva aunque la investigación termine en
  `failed`.

## Qué se ensambla de verdad hoy (nada inventado)

- `deployment_metadata`: SÍ, a partir de metadata no secreta de `Project`
  (nombre, estado de monitoreo/salud, snapshot de `GitHubRepository` y
  `DokployApplication`) — mismo patrón de snapshot que H1, resuelto vía
  `Investigation.IncidentID` → `Incident.ProjectID` → `Project`.
- `log_excerpt`, `stack_trace`: pendientes de H2 (LogEvent/Incident real).
- `code_location`, `commit`, `diff`: pendientes de PB-030/PB-031 (lectura de
  repositorio/commits).
- `validation_result`: pertenece a Remediation (H5), fuera de alcance.

No se generan placeholders para estos tipos: el assembler simplemente no
produce Evidence de esos tipos todavía.

## Frontera con H2 (extensión de un port existente, no reimplementación)

`out.IncidentReader` (definido en PB-026, hoy solo `Exists`) gana un método
`Get(ctx, incidentID) (*domain.Incident, error)` — es lo mínimo que
`EvidenceAssembler` necesita para llegar a `Incident.ProjectID`. El stub de
producción (`stub.DenyAllIncidentReader`) responde `ErrIncidentNotFound` en
`Get`, igual que `Exists=false`; el assembler trata ese caso (y un Project ya
borrado) como "sin evidencia disponible todavía" — slice vacío, no error —
en vez de fallar la investigación completa por un dato de contexto opcional.
