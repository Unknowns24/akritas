# Implementation brief — AKR-DOKPLOY-COMPOSE-SOURCES

## Estado inicial

Akritas descubre únicamente `/api/application.search`, representa cada Project
con un `DokployApplication` embebido y obtiene logs desde
`/api/application.readLogs`. Por eso una instalación Dokploy que contiene solo
Composes aparece sin aplicaciones y no puede convertirse en fuente de
monitoreo.

La identidad actual está repetida como `application_identifier` e
`instance_identifier` en Projects, checkpoints, eventos y evidencia. El cambio
debe generalizar esa identidad sin perder filas existentes ni introducir IDs de
contenedor inestables.

## Modelo de dominio

Se incorporarán `DokployCompose`, `DokployComposeService` y `DokploySource`.
`DokploySource` será el snapshot embebido en Project y contendrá:

- `type`: `application` o `compose_service`;
- `dokploy_server_id`;
- `resource_identifier`: applicationId o composeId;
- `service_name`: obligatorio solo para `compose_service`;
- `instance_identifier`: appName de Dokploy;
- `display_name`, `environment` y estado normalizado;
- `runtime_type`: `docker-compose` o `stack`, solo para Compose;
- identificador opcional del servidor remoto administrado por Dokploy, solo en
  persistencia interna.

La identidad exclusiva será
`dokploy_server_id + type + resource_identifier + service_name normalizado`.
Para application, `service_name` debe estar vacío. Un Compose podrá asociarse a
varios Projects únicamente cuando cada uno seleccione un servicio diferente.

## API REST y OpenAPI

Se mantendrá:

```text
GET /api/v1/integrations/dokploy/servers/{server_id}/applications
```

Se agregarán:

```text
GET /api/v1/integrations/dokploy/servers/{server_id}/composes
GET /api/v1/integrations/dokploy/servers/{server_id}/composes/{compose_id}/services?refresh=false
```

La lista de Composes usará los mismos límites y cursores Uker firmados que la
lista de aplicaciones, con scope propio y filtro allowlisted `name_like`. La
lista de servicios no será paginada: deduplicará y ordenará alfabéticamente la
respuesta del proveedor. `refresh=false` llamará `compose.loadServices` con
`type=cache`; `refresh=true`, con `type=fetch`.

Create/Update Project reemplazarán los campos superiores
`dokploy_server_id + application_identifier` por un objeto completo
`dokploy_source`. Project y ProjectSummary reemplazarán `dokploy_application`
por el snapshot `dokploy_source`. En PATCH, si cambia la fuente, el selector
completo será obligatorio y continuará aplicando la regla de monitoring
desactivado.

`DokployServer` conservará `application_count` y agregará `compose_count`. Cada
listado actualizará solamente su conteo y `last_synced_at`. El contrato OpenAPI
pasará a `2.0.0` y declarará los dos selectors mediante `oneOf` y discriminador.

## Integración Dokploy y logs

El gateway incorporará:

- `compose.search` para discovery y total paginado;
- `compose.one` para resolver y snapshotear un Compose exacto;
- `compose.loadServices` para validar/listar servicios;
- `docker.getContainersByAppNameMatch` para localizar réplicas actuales;
- `compose.readLogs` para obtener logs del contenedor elegido.

Para cada fetch Compose, el adapter consultará los contenedores actuales usando
`instance_identifier`, `runtime_type` y el serverId remoto opcional. Filtrará
solo contenedores running con las etiquetas del proyecto y servicio:

- docker-compose: `com.docker.compose.project` y
  `com.docker.compose.service`;
- stack: `com.docker.stack.namespace` y
  `com.docker.swarm.service.name`.

Los candidatos se ordenarán por ID ascendente y se elegirá el primero. Después
se invocará `compose.readLogs` con composeId, containerId, tail y since. Resolver
el contenedor en cada ciclo permite continuar tras un redeploy sin guardar el
ID. Si el servicio está declarado pero no tiene un contenedor activo, Project
puede existir, pero el ciclo de monitoreo falla de forma normalizada y pasa a
degraded según el comportamiento vigente del worker.

## Persistencia y compatibilidad de datos

Una migración reversible posterior a `20260823_05` realizará, dentro de una
transacción:

- renombre de columnas de aplicación en `projects` a columnas `source_*`;
- alta y backfill de `source_type = 'application'`;
- alta de service name, runtime type y serverId remoto opcionales;
- reemplazo de `uq_projects_dokploy_application` por un índice único que use
  `COALESCE(source_service_name, '')`;
- checks coherentes con el discriminante;
- generalización de las columnas de fuente en `monitoring_checkpoints` y
  `log_events`, con backfill application;
- alta de `compose_count` en `dokploy_servers`.

El rollback restaurará nombres, constraints e índices anteriores y solo podrá
rechazarse explícitamente si existen fuentes Compose que no pueden representarse
en el schema antiguo. No se modificarán blobs Evidence ya persistidos; la
evidencia nueva usará nombres genéricos de fuente.

## Errores y seguridad

Se agregarán errores estables para selector inválido, Compose inexistente,
servicio inexistente, fuente ya asociada y ausencia de contenedor activo. Los
401/403/404/5xx del proveedor continuarán pasando por la normalización vigente.

No se loguearán credenciales, respuestas crudas sensibles, archivos Compose ni
variables de entorno. El refresh remoto será siempre explícito. Los IDs y
etiquetas Docker se tratarán como metadata operativa y se validarán antes de
construir queries al proveedor.
