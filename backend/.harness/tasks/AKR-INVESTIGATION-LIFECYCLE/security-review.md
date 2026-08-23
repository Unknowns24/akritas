# Security review — AKR-INVESTIGATION-LIFECYCLE

## Veredicto

Sin hallazgos abiertos.

## Controles verificados

- Las 4 rutas requieren sesión de Administrator; `startIncidentInvestigation`
  además exige Origin exacto allowlisted (verificado manualmente: 401 sin
  sesión, 403 sin Origin).
- El `IncidentReader` de producción (`stub.DenyAllIncidentReader`) es
  deny-by-default: nunca autoriza de más mientras H2 no exista, aunque las
  rutas queden expuestas y auditables (verificado manualmente: 404 en los 3
  endpoints que dependen de la existencia del incidente).
- `Idempotency-Key` se valida como UUID antes de tocar cualquier dependencia
  (400 si falta o es inválida); el replay de una key ya usada devuelve la
  misma Operation sin crear estado nuevo ni volver a encolar trabajo.
- El pipeline asíncrono corre con su propio `context.WithTimeout` desacoplado
  de la request HTTP (no hereda cancelación ni puede quedar colgado
  indefinidamente si la request original ya terminó).
- Ninguna llamada externa ocurre dentro de una transacción de base de datos;
  cada transición de estado se persiste de forma discreta (`background-
  processes.md`).
- `InvestigationDTO`/`OperationDTO` solo exponen campos ya públicos en el
  contrato; `Operation.IdempotencyKey` nunca se serializa (no está en el
  schema `Operation` del OpenAPI y el DTO no tiene ese campo).
- Los mensajes de fallo (`FailureUserMessage`/`Operation.UserMessage`) se
  originan en el propio proceso de investigación (dominio interno o el error
  del runner), nunca en input del cliente ni en detalles de infraestructura
  no saneados.

El gate `.harness/kernel/scripts/check-security.sh` finalizó correctamente.
