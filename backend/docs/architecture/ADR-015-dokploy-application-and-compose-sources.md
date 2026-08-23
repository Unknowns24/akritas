# ADR-015 — Fuentes Dokploy application y compose_service

## Estado

Accepted

## Contexto

Dokploy expone Applications y Composes como recursos distintos. Akritas
modelaba Project exclusivamente con un snapshot `DokployApplication` y usaba
`application.readLogs`; una instalación que ejecutaba sus workloads mediante
Compose aparecía vacía y no podía ser monitoreada.

Los contenedores de un Compose tampoco son identidades durables: cambian durante
redeploys y un servicio Swarm puede tener múltiples réplicas.

## Decisión

Project embebe un `DokploySource` discriminado con dos variantes:

- `application`, identificada por Dokploy server y applicationId;
- `compose_service`, identificada por Dokploy server, composeId y service name.

La misma identidad de fuente sólo puede pertenecer a un Project. Servicios
distintos de un mismo Compose son identidades diferentes.

El snapshot guarda metadata no secreta necesaria para operar: appName,
displayName, environment, status y, para Compose, runtime type y serverId remoto
opcional. Nunca guarda el containerId.

Para cada lectura de logs Compose, el adapter consulta los contenedores actuales,
filtra las etiquetas oficiales de docker-compose o stack, selecciona solamente
réplicas running, ordena por containerId y elige la primera. Luego invoca
`compose.readLogs`. Esta resolución dinámica permite sobrevivir redeploys sin
alterar la identidad lógica ni el checkpoint.

La lista de servicios usa caché por defecto. Un refresh de la fuente remota sólo
ocurre mediante el parámetro público explícito `refresh=true`.

## Consecuencias

- el contrato Project cambia de `dokploy_application` a `dokploy_source`;
- Projects, checkpoints y LogEvents existentes se migran como `application`;
- la evidencia nueva usa identidad de fuente genérica y la histórica permanece
  inmutable;
- un servicio puede asociarse antes de tener contenedor activo, pero el ciclo de
  monitoreo queda degraded hasta que exista una réplica running;
- esta versión no agrega streams ni cursores independientes por réplica.

## Alternativas descartadas

- Guardar containerId: se vuelve obsoleto en cada redeploy.
- Monitorear el Compose completo: mezcla logs de servicios con responsabilidades
  diferentes y rompe la asociación uno-a-uno con Project.
- Agregar todos los logs de todas las réplicas: requiere cursores por réplica y
  deduplicación adicional fuera del alcance actual.
