# Security review — AKR-INVESTIGATION-EVIDENCE

## Veredicto

Sin hallazgos abiertos.

## Controles verificados

- La ruta nueva requiere sesión de Administrator, en el mismo grupo que el
  resto (verificado manualmente: 401 sin sesión).
- `deployment_metadata` solo incluye campos de `Project` ya expuestos por
  `ProjectSummaryDTO` (nombre, monitoring/health status, repo owner/name/
  branch, aplicación Dokploy/entorno/estado); nunca tokens, PATs ni
  credenciales de integraciones.
- `Evidence.Redacted` es siempre `true` (invariante ya impuesta por
  `domain.Evidence.Validate`, no reimplementada); el mapper REST no la
  pisa.
- El `IncidentReader` de producción sigue siendo deny-by-default: `Get`
  también responde `ErrIncidentNotFound`, así que `EvidenceAssembler` nunca
  filtra datos de un incidente/proyecto real hasta que H2 aporte un reader
  verdadero (verificado manualmente contra un server real: la investigación
  insertada a mano en la base sin pasar por el pipeline nunca contiene
  evidencia hasta que se inserta explícitamente).
- Un fallo de `EvidenceAssembler` que no sea "no encontrado" se propaga tal
  cual, sin exponer detalles internos en la respuesta (el pipeline async no
  produce respuesta HTTP directamente; el error solo queda en el estado de
  Operation/Investigation la próxima vez que se pueda persistir).
- `type_in` se valida contra el enum del dominio antes de llegar a
  Uker/GORM (verificado manualmente: 400 con un valor fuera del enum).
- El listado queda scoped por `investigation_id` a nivel de query SQL, no
  solo a nivel de filtro de aplicación.

El gate `.harness/kernel/scripts/check-security.sh` finalizó correctamente.
