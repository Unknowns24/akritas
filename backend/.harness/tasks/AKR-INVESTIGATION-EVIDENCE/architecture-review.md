# Architecture review — AKR-INVESTIGATION-EVIDENCE

## Veredicto

Conforme con el profile `backend_api` y el workflow `backend-api-feature`.

## Revisión

- `internal/core/domain/evidence.go` solo ganó tags GORM; `NewEvidence` y
  `Validate` quedaron intactos, igual que el criterio ya usado con
  `Investigation`/`Operation` en PB-026.
- `out.IncidentReader` se extendió con `Get` en vez de crear un port
  paralelo: es el mismo boundary hacia H2 que PB-026 ya definió, ahora con
  lo mínimo adicional que `EvidenceAssembler` necesita
  (`Incident.ProjectID`). `stub.DenyAllIncidentReader` y el fake de tests se
  extendieron en consecuencia.
- `internal/service/evidenceassembly` es hermano de
  `investigationdispatch`: mismo criterio de "servicio = orquestación
  reutilizable sobre ports existentes", sin SDK de terceros de por medio.
  No inventa datos: "no encontrado" en cualquiera de sus dos pasos devuelve
  slice vacío, nunca una Evidence ficticia.
- `RunUseCase` (investigación) creció dos dependencias (`EvidenceStore`,
  `EvidenceAssembler`) sin romper la separación ya establecida en PB-026
  entre `UseCase` (REST-facing) y `RunUseCase` (solo dispatcher) — ninguna
  de las dos nuevas dependencias reintroduce el ciclo de construcción que
  esa separación evita.
- `internal/usecase/evidence` es un paquete propio de una sola
  responsabilidad, mismo criterio que `usecase/operation`, en vez de
  agregar un método más a `InvestigationUseCase`.
- `EvidenceStore` no tiene `Update`: Evidence es write-once, reflejado
  directamente en la forma del port en vez de dejarlo como una convención
  implícita.
- REST sigue el patrón ya establecido: `EvidenceDTO` con sufijo `DTO`,
  mapper de una responsabilidad, handler con `r.PathValue`/
  `restpagination.Parse`/`response.Invalid/Error/JSON`, envelope genérico
  `commondto.ListResponseDTO` reutilizado sin crear
  `EvidenceListResponseDTO`. La ruta se agregó al mismo grupo `Admin`+
  `Origin` que el resto, sin excepciones.
- `internal/bootstrap/integrations/module.go` reutiliza la misma instancia
  de `stub.DenyAllIncidentReader` tanto para `InvestigationUseCase` como
  para `EvidenceAssembler`, evitando dos stubs deny-all independientes para
  el mismo boundary.

No se detectaron dependencias invertidas ni lógica de negocio filtrada a
adapters o handlers.
