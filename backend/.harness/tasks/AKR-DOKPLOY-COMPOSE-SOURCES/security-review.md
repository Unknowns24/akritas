# Security Review

## Summary

La expansión Compose reutiliza el Credential Store y el cliente Dokploy
existentes, no agrega secretos al dominio público y normaliza fallos remotos sin
devolver payloads sensibles.

## Auth / permissions

- Las dos rutas nuevas se registran bajo el grupo autenticado de integraciones
  Dokploy existente.
- No se relajaron middleware, sesión ni permisos.
- Permisos insuficientes del API key siguen el error normalizado de integración.

## Input validation

- Server UUID, compose ID, cursor, límite y `name_like` se validan antes de
  delegar al caso de uso.
- `refresh` acepta exclusivamente un único valor `true` o `false`; vacío,
  duplicado o parámetros desconocidos se rechazan.
- El discriminante prohíbe service/runtime en applications y exige ambos en
  compose_service.
- Los filtros de contenedor requieren estado running y coincidencia exacta de
  labels de proyecto y servicio.

## Data exposure

- API keys se obtienen como bytes, se borran tras el request y no forman parte de
  DTOs, errores ni evidencia.
- El ID remoto de servidor se persiste sólo para resolver contenedores y no se
  expone en Project responses.
- Container IDs no se persisten; se resuelven nuevamente en cada lectura.

## Error leakage

- 401/403, 404 y errores del proveedor se convierten al catálogo estable.
- Bodies malformados se reportan como integración no disponible sin incluir la
  respuesta cruda.
- La ausencia de contenedor activo usa un error estable y no revela inventario
  Docker.

## Findings

No se encontraron secretos hardcodeados ni hallazgos bloqueantes. `refresh=true`
puede provocar acceso remoto al repositorio configurado en Dokploy, pero sólo se
activa por solicitud explícita y nunca es el valor por defecto.

## Result

pass
