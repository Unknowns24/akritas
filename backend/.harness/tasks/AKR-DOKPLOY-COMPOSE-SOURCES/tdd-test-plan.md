# TDD test plan — AKR-DOKPLOY-COMPOSE-SOURCES

Estado: **aprobado por el usuario el 2026-08-23**.

La implementación queda habilitada conforme al workflow
`backend-api-feature`.

## 1. Dominio y contratos internos

- `DokploySource` acepta application sin service/runtime y rechaza service name
  o runtime Compose en esa variante.
- `DokploySource` acepta compose_service únicamente con composeId, service name,
  instance identifier y runtime `docker-compose` o `stack`.
- rechaza UUID nil, strings vacíos, tipo desconocido, status inválido y
  combinaciones cruzadas del discriminante.
- `IdentityKey` es estable, normaliza whitespace y diferencia dos servicios del
  mismo Compose.
- Project valida, reemplaza y refresca snapshots de ambas variantes; cambiar la
  identidad con monitoring activo sigue devolviendo conflicto.
- checkpoint y LogEvent copian type, resource, service e instance correctos.

## 2. Adapter Dokploy

### Discovery de Composes

- `compose.search` recibe limit, offset y `q` derivado de `name_like`.
- mapea `{items,total}`, nombres alternativos, environmentId y status conocido;
  un status futuro se normaliza a unknown.
- calcula next/prev boundary sin filtrar o duplicar páginas.
- rechaza cursor provider inválido, payload vacío/malformado y total incoherente.
- normaliza 401/403 como credencial rechazada y 404/5xx/timeouts como error de
  integración según el catálogo vigente.

### Servicios Compose

- `refresh=false` envía `type=cache`; `refresh=true` envía `type=fetch`.
- mapea un array de nombres, recorta whitespace, elimina vacíos/duplicados y
  ordena el resultado.
- diferencia Compose inexistente, servicio inexistente y respuesta inválida.
- resolver un selector Compose llama `compose.one`, valida runtime y confirma el
  service name con `compose.loadServices` sin refresh forzado.

### Resolución de contenedor y logs

- application continúa llamando `application.readLogs` con los mismos tail y
  since.
- compose_service consulta `docker.getContainersByAppNameMatch` con appName,
  appType y serverId remoto cuando exista.
- docker-compose filtra exactamente las labels project/service y descarta
  contenedores stopped o de otro servicio.
- stack filtra namespace/service labels y descarta réplicas ajenas.
- con varias réplicas running ordena por container ID y usa la primera.
- un fetch posterior puede seleccionar un container ID nuevo tras un redeploy.
- sin candidato running devuelve el error estable de contenedor no disponible.
- `compose.readLogs` recibe composeId, containerId, tail y since correctamente
  escapados.
- el parser conserva timestamps, ordinales, multiline y hashes para ambos tipos
  de fuente.

## 3. Casos de uso Project e integración

- ListComposes carga el servidor, delega al gateway, actualiza compose_count y
  last_synced_at; un fallo no persiste un conteo parcial.
- ListComposeServices valida server/compose/refresh y no modifica conteos.
- Create Project resuelve y persiste application y compose_service.
- un servicio declarado sin contenedor activo puede crear Project.
- selector inexistente no persiste Project ni checkpoint.
- la misma application o el mismo compose_service producen conflicto.
- dos servicios diferentes del mismo Compose se aceptan.
- Update requiere el selector completo y monitoring disabled cuando cambia.
- activar monitoring vuelve a resolver el snapshot y valida que el servicio aún
  existe.
- cambios concurrentes conservan los controles de versión actuales.

## 4. PostgreSQL

- la migración transforma una fila Project application existente sin perder
  metadata ni relaciones.
- checkpoints y LogEvents existentes quedan identificados como application y
  conservan cursores, ocurrencias e incidentes.
- persiste y recupera snapshots de las dos variantes.
- el índice rechaza la misma identidad exacta y permite servicios distintos del
  mismo Compose.
- los checks rechazan service/runtime inválidos para cada tipo.
- compose_count inicia en cero y se actualiza sin modificar application_count.
- rollback restaura el schema anterior cuando no hay fuentes Compose y falla de
  forma explícita/no destructiva cuando sí las hay.

## 5. REST y OpenAPI

- inventario de rutas conserva applications y agrega composes/services bajo el
  prefijo autenticado correcto.
- lista Compose devuelve envelope/paging y cursores Uker con scope separado.
- services devuelve lista ordenada y valida `refresh` estrictamente como boolean.
- UUID, composeId, filtros y cursores inválidos producen Problem Details estable.
- Create/Update aceptan cada rama válida de `dokploy_source`.
- rechazan el contrato anterior y combinaciones inválidas del discriminador.
- Project/ProjectSummary devuelven solamente `dokploy_source`, sin
  `dokploy_application` ni campos selectores superiores.
- DokployServer devuelve application_count y compose_count.
- OpenAPI incluye ejemplos, required/additionalProperties y responses de error;
  el validador confirma que implementación y contrato coinciden.

## 6. Monitoreo, evidencia y seguridad

- el worker construye LogFetchRequest con la fuente exacta para application y
  compose_service.
- rotación de contenedor no invalida el checkpoint lógico de la fuente.
- ausencia de réplica activa degrada el ciclo sin perder el cursor anterior.
- persist_occurrence guarda identidad genérica completa en LogEvent.
- evidencia nueva de log/deployment incluye tipo, resource, service e instance;
  no altera evidencia histórica.
- fixtures con API keys, tokens, env vars y secretos verifican redacción y
  ausencia en logs/respuestas/evidence.
- permisos Dokploy insuficientes no exponen body, URL con secretos ni credential.

## 7. Gates de validación

Después de implementar, ejecutar:

```text
go test ./...
go fmt ./...
.harness/kernel/scripts/check-backend-architecture.sh
.harness/kernel/scripts/check-openapi.sh
.harness/kernel/scripts/check-security.sh
```

Además:

- generar `implementation-summary.md` con archivos y validaciones;
- completar `architecture-review.md` y resolver findings bloqueantes;
- completar `security-review.md` y resolver findings bloqueantes;
- cerrar con `final-summary.md`, indicando cualquier validación no ejecutable.

## Criterio de aprobación

La implementación se considera aceptada cuando todos los tests anteriores
pasan, OpenAPI valida, los tres scripts del harness finalizan correctamente y no
quedan findings bloqueantes de arquitectura o seguridad.
