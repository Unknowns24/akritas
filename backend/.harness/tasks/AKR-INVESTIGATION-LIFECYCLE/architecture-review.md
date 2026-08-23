# Architecture review — AKR-INVESTIGATION-LIFECYCLE

## Veredicto

Conforme con el profile `backend_api` y el workflow `backend-api-feature`.

## Revisión

- `internal/core/domain` gana `Operation` (mismo estilo de máquina de estados
  que `Investigation`) y tags GORM en `Investigation`; sigue sin importar
  GORM, chi ni adapters (`check-backend-architecture.sh` en verde).
- `RunInvestigationUseCase` (`RunUseCase`) y `InvestigationUseCase`
  (`UseCase`) se separaron deliberadamente en dos structs dentro del mismo
  paquete `usecase/investigation`: `RunUseCase` no depende de
  `InvestigationDispatcher` y `UseCase` no depende de `InvestigationRunner`,
  lo que evita un ciclo de construcción (`InvestigationDispatcher` envuelve a
  `RunUseCase`; `UseCase` depende de `InvestigationDispatcher`).
- `Operation` está diseñado como infraestructura genérica: dominio, puertos
  (`OperationStore`) y persistencia no mencionan `investigation` en ningún
  tipo — solo `resource_type`/`resource_id` los conectan, exactamente como
  pide la tarea para que H5 lo reutilice.
- `internal/service/investigationdispatch` es la única pieza en `internal/
  service`; implementa `InvestigationDispatcher` con una goroutine simple,
  documentada como mínima a propósito (nada de worker pool/queue).
- `internal/adapter/external/qvac` aloja el stub del runner en el mismo lugar
  donde vivirá la implementación real de QVAC (ADR-001), y
  `internal/adapter/stub` es un paquete nuevo, explícitamente temporal, para
  el `IncidentReader` de producción — no se mezcló con `qvac` porque
  responde a una frontera distinta (H2, no QVAC).
- REST sigue el patrón de Project: DTOs con sufijo `DTO` uno por archivo,
  mappers de una responsabilidad por archivo, handlers con `r.PathValue`,
  `request.DecodeJSON`/`request.IdempotencyKey`, `response.Invalid/Error/
  JSON`, envelopes genéricos (`commondto.DataResponseDTO`/`ListResponseDTO`)
  reutilizados sin duplicar `InvestigationResponseDTO`/`OperationResponseDTO`.
- El router agrega las 3 rutas nuevas al mismo grupo `Admin`+`Origin` que ya
  usa Project, siguiendo el agrupamiento grueso ya establecido en el repo en
  lugar de inventar un split fino por verbo.
- `internal/bootstrap/integrations/module.go` extiende el mismo composition
  root que ya wirea Project (no se creó un segundo bootstrap), consistente
  con que Project tampoco es una "integración" en sentido estricto y ya vive
  ahí.

No se detectaron dependencias invertidas ni lógica de negocio filtrada a
adapters o handlers.
