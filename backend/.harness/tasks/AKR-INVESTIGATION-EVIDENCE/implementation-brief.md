# Implementation brief — AKR-INVESTIGATION-EVIDENCE

## Estado inicial

`internal/core/domain/evidence.go` está completo (7 tipos, `NewEvidence`,
`Validate`) y sin tags gorm, igual que `Investigation` antes de PB-026. No
existe port, usecase, persistencia ni REST para Evidence.
`RunUseCase.Execute` (PB-026) ya orquesta pending→running→completed/failed
pero no toca evidencia; `Investigation.EvidenceCount` existe pero nadie lo
escribe.

## Estrategia

```text
REST/Chi -> ListInvestigationEvidence -> InvestigationStore.FindByID (404 si no existe)
                                       -> EvidenceStore.ListByInvestigation

RunUseCase.Execute (ya existe) -> ... Start() persistido ...
                                -> EvidenceAssembler.Assemble(investigation)
                                     -> IncidentReader.Get (extensión del port de PB-026)
                                     -> ProjectStore.Get (ya existe, H1)
                                -> EvidenceStore.Create (una vez por Evidence)
                                -> Investigation.EvidenceCount = len(...); Update()
                                -> ... Run()/Complete()/Fail() sin cambios ...
```

Un fallo de `Assemble` (error real de infraestructura, no un "no encontrado")
se propaga igual que ya propagan los `Update`/`Create` intermedios de
`Execute` hoy: `return err` directo, sin pasar por `failInvestigation` — es
el mismo criterio ya usado en PB-026 (fallas de infra se dejan ver tal cual;
solo el resultado del `InvestigationRunner` es un "outcome de negocio"
capturado en el dominio). Un "no encontrado" (incident/project) no es un
error: `Assemble` devuelve slice vacío.

## Dominio

- `Evidence`: se agregan tags `gorm` a los campos existentes (mismo criterio
  que `Investigation`/`Operation` en PB-026), sin tocar `NewEvidence` ni
  `Validate`.
- Sin códigos de error nuevos en el dominio: `ErrInvalidEvidence`/
  `ErrInvalidEvidenceType` (405) y `ErrIncidentNotFound`/
  `ErrInvestigationNotFound` (504) ya existen y se reutilizan tal cual.

## Ports

- `out.IncidentReader` (extensión): se agrega `Get(ctx, incidentID)
  (*domain.Incident, error)`, devuelve `domain.ErrIncidentNotFound` si no
  existe. `stub.DenyAllIncidentReader` implementa el nuevo método igual que
  `Exists` (siempre deniega). El fake de tests de `usecase/investigation`
  también se extiende.
- `out.EvidenceStore` (nuevo) — `Create`, `ListByInvestigation` (paginado,
  `paging.Params`/`paging.Slice`). Sin `Update`: Evidence es write-once.
- `out.EvidenceAssembler` (nuevo) — `Assemble(ctx, domain.Investigation)
  ([]domain.Evidence, error)`.

## Usecases

- `internal/usecase/investigation/run_investigation.go` (extensión de
  `RunUseCase`): nuevas dependencias `evidence out.EvidenceStore` y
  `assembler out.EvidenceAssembler` en el constructor; `Execute` ensambla,
  persiste cada Evidence y actualiza `EvidenceCount` entre el `Start()` ya
  persistido y la llamada a `InvestigationRunner.Run`.
- `internal/usecase/evidence/` (nuevo, paquete propio — mismo criterio que
  `usecase/operation`): `ListInvestigationEvidence` valida que la
  investigación exista (`InvestigationStore.FindByID`, reusa
  `ErrInvestigationNotFound` para el 404) y delega en
  `EvidenceStore.ListByInvestigation`.

## Servicio — evidenceassembly

`internal/service/evidenceassembly/` (nuevo, hermano de
`investigationdispatch`) implementa `out.EvidenceAssembler`: resuelve
`Investigation.IncidentID` → `Incident.ProjectID` → `Project` y arma un
único `Evidence` tipo `deployment_metadata` con nombre del proyecto, estado
de monitoreo/salud y snapshot de `GitHubRepository`/`DokployApplication`
(nunca credenciales — son los mismos campos que ya expone
`ProjectSummaryDTO`). "No encontrado" en cualquiera de los dos pasos
devuelve slice vacío sin error, no un placeholder.

## Persistencia

Migración `20260822_12_add_evidence` (mismo formato que las anteriores):
`evidence.investigation_id uuid NOT NULL REFERENCES investigations(id) ON
DELETE CASCADE` (FK real, a diferencia del gap documentado de
`investigations.incident_id`), check de `type` contra los 7 valores del
enum, check `redacted = true`, check de coherencia `line_start`/`line_end`.
Error de persistencia nuevo: `ErrEvidencePersistence` (`0x206001I`,
siguiente componente libre de la capa DB).

## REST

`EvidenceDTO` refleja el schema `Evidence` campo a campo (mismos opcionales
que el schema: `content`, `file_path`, `line_start`, `line_end`,
`commit_sha`, `patch`, `occurred_at` con `omitempty`). Handler nuevo
`internal/adapter/rest/handler/evidence/` con un único método `List`;
`type_in` se valida contra `domain.EvidenceType.Validate()` por cada valor
(mismo patrón que `monitoring_status_in` en Project). Ruta nueva
`GET /investigations/{investigation_id}/evidence`, registrada en el mismo
grupo `Admin`+`Origin` que el resto (GET no muta, pero es la convención ya
establecida del repo).

## Wiring de producción

`internal/bootstrap/integrations/module.go` agrega el repositorio de
Evidence, el `evidenceassembly.Assembler` (reusando la misma instancia de
`stub.DenyAllIncidentReader` que ya usa `InvestigationUseCase`), y pasa
ambos al `RunUseCase` extendido. Como el `IncidentReader` de producción
sigue denegando todo, este código es real pero, igual que `InvestigationRunner`,
no se ejerce en producción hasta que H2 mergee — sí se ejerce completo en
tests con fakes.
