# ADR-007 — Contrato y semántica de `MonitoringConfiguration`

## Estado

Accepted

## Contexto

Akritas permite configurar el monitoreo de manera independiente para cada `Project`.

El modelo de dominio ya define una entidad/value object `MonitoringConfiguration`, utilizada por el Detection Engine descrito en ADR-006.

Hasta ahora se habían identificado los siguientes atributos:

```text
enabled
error_patterns
ignored_patterns
grouping_window
context_before
context_after
```

Sin embargo, no estaban definidos de manera precisa:

- sus tipos;
- la semántica de los patrones;
- la precedencia entre reglas;
- el comportamiento exacto de `grouping_window`;
- la unidad de `context_before` y `context_after`;
- los valores por defecto;
- qué parámetros operativos pertenecen o no a esta configuración.

Esta ambigüedad dejaría decisiones importantes en manos de la implementación.

## Decisión

Cada `Project` tendrá exactamente una `MonitoringConfiguration`.

Conceptualmente:

```text
Project
    ↓
MonitoringConfiguration
```

La configuración tendrá el siguiente contrato:

```text
MonitoringConfiguration

enabled: bool

error_patterns: []string

ignored_patterns: []string

grouping_window: duration

context_before: int

context_after: int
```

---

## `enabled`

Determina si el monitoring del `Project` se encuentra activo.

```text
enabled = false
```

implica que Akritas no debe consumir ni procesar nuevos logs para ese proyecto.

Desactivar monitoring:

- no elimina configuración;
- no elimina `LogEvent`;
- no elimina `Incident`;
- no elimina `Investigation`;
- no modifica el estado de incidentes existentes.

El valor por defecto será:

```text
false
```

El monitoring debe ser habilitado explícitamente.

---

## `error_patterns`

Lista de patrones adicionales definidos para el `Project` que permiten considerar un evento como relevante.

Los patrones serán **expresiones regulares**.

Ejemplo conceptual:

```text
connection refused

database .* unavailable

payment provider timeout
```

Los patrones complementan las reglas built-in del Detection Engine definidas por ADR-006.

Por lo tanto:

```text
builtin rule matches
        OR
error_pattern matches
        ↓
candidate detection
```

Un patrón inválido debe ser rechazado al guardar la configuración.

La ausencia de patrones personalizados no deshabilita las reglas built-in.

Valor por defecto:

```text
[]
```

---

## `ignored_patterns`

Lista de expresiones regulares utilizadas para descartar eventos conocidos o irrelevantes.

Ejemplos:

```text
healthcheck failed

expected test error

connection reset by monitoring probe
```

Los `ignored_patterns` tienen precedencia sobre cualquier regla de detección.

Por lo tanto:

```text
ignored_pattern matches
        ↓
      discard
```

incluso cuando el mismo evento coincida también con:

- una regla built-in;
- un `error_pattern`.

Formalmente:

```text
if ignored(event):
    discard

else if detected(event):
    process
```

Un evento ignorado:

- no crea `LogEvent`;
- no genera fingerprint;
- no crea ni actualiza `Incident`;
- no dispara `Investigation`;
- no llega a QVAC.

Valor por defecto:

```text
[]
```

---

## `grouping_window`

Define durante cuánto tiempo una nueva ocurrencia puede incorporarse a un `Incident` existente con el mismo fingerprint.

Valor por defecto:

```text
30m
```

Una nueva ocurrencia será agrupada en un `Incident` existente cuando se cumplan simultáneamente:

```text
same Project

AND

same fingerprint

AND

Incident is eligible for grouping

AND

occurrence.timestamp - Incident.last_seen_at <= grouping_window
```

Si cualquiera de estas condiciones no se cumple, se crea un nuevo `Incident`.

La ventana se calcula respecto de:

```text
Incident.last_seen_at
```

y no respecto de `first_seen_at`.

Esto permite que un problema que continúa ocurriendo mantenga activo el mismo incidente mientras las nuevas ocurrencias permanezcan dentro de la ventana configurada.

Ejemplo:

```text
grouping_window = 30m

10:00 error X
10:10 error X
10:35 error X
10:50 error X
```

Todas las ocurrencias pertenecen al mismo incidente porque ninguna ocurre más de 30 minutos después de la anterior.

En cambio:

```text
10:00 error X
10:10 error X
11:00 error X
```

la ocurrencia de `11:00` inicia un nuevo incidente.

---

## `context_before`

Cantidad máxima de **log records anteriores** que Akritas conservará como contexto de una ocurrencia detectada.

Valor por defecto:

```text
20
```

Estos registros:

- forman parte de la evidencia;
- pueden ser utilizados posteriormente por QVAC;
- no participan del fingerprint salvo que explícitamente formen parte del evento detectado.

---

## `context_after`

Cantidad máxima de **log records posteriores** que Akritas conservará como contexto de una ocurrencia detectada.

Valor por defecto:

```text
20
```

Estos registros tienen la misma finalidad que `context_before`.

---

## Unidad de contexto

`context_before` y `context_after` representan **log records recibidos desde la fuente**, no `LogEvent`.

Esto es necesario porque el contexto es capturado antes de determinar qué registros individuales representan eventos relevantes del dominio.

Conceptualmente:

```text
raw log records

20 anteriores
      ↓
detected event
      ↓
20 posteriores
```

El `LogEvent` representa posteriormente la ocurrencia relevante y puede conservar esta evidencia asociada.

---

## Valores por defecto

Una configuración creada sin personalización utilizará:

```text
enabled = false

error_patterns = []

ignored_patterns = []

grouping_window = 30m

context_before = 20

context_after = 20
```

Estos valores forman parte del comportamiento definido de Akritas y no deben depender de decisiones particulares de cada adapter.

---

## Validación

Una `MonitoringConfiguration` válida debe cumplir:

```text
grouping_window > 0

context_before >= 0

context_after >= 0

todos los error_patterns son regex válidas

todos los ignored_patterns son regex válidas
```

La configuración inválida debe rechazarse antes de activar monitoring.

---

## Polling de logs

La frecuencia con la que Akritas consulta Dokploy **no forma parte de `MonitoringConfiguration`**.

Por ejemplo:

```text
polling_interval
```

no será un atributo del dominio.

El polling es una preocupación operativa de infraestructura y pertenece a la configuración del runtime de Akritas.

Conceptualmente:

```text
Runtime Configuration
    ↓
polling interval

Project
    ↓
MonitoringConfiguration
    ↓
qué detectar y cómo agrupar
```

Esto evita que detalles sobre cómo Akritas obtiene los logs contaminen la configuración conceptual de monitoring de cada proyecto.

En el futuro diferentes adapters podrían incluso utilizar mecanismos distintos:

```text
polling
streaming
websocket
event subscription
```

sin modificar `MonitoringConfiguration`.

---

## Relación con el Detection Engine

ADR-006 define el pipeline:

```text
logs
 ↓
event assembly
 ↓
ignored rules
 ↓
detection rules
 ↓
normalization
 ↓
fingerprint
 ↓
grouping
 ↓
context
```

`MonitoringConfiguration` parametriza ese pipeline:

```text
ignored_patterns
        ↓
Ignored Rules

error_patterns
        ↓
Detection Rules

grouping_window
        ↓
Incident Grouping

context_before
context_after
        ↓
Context Capture
```

`enabled` determina si el pipeline debe ejecutarse para el `Project`.

---

## Invariantes

- cada `Project` posee una única `MonitoringConfiguration`;
- monitoring está deshabilitado por defecto;
- `ignored_patterns` siempre tiene precedencia sobre detección;
- `error_patterns` complementa, pero no reemplaza, las reglas built-in;
- los patrones deben ser válidos antes de activar monitoring;
- `grouping_window` se calcula respecto de `last_seen_at`;
- el contexto se expresa en cantidad de log records;
- la frecuencia de adquisición de logs no pertenece al dominio;
- modificar `MonitoringConfiguration` no altera retroactivamente incidentes ya procesados.

## Consecuencias

### Positivas

- el Detection Engine tiene un contrato explícito;
- la implementación no necesita inventar defaults;
- el comportamiento de agrupación queda definido;
- los patrones personalizados tienen semántica predecible;
- las reglas de exclusión tienen precedencia inequívoca;
- infraestructura y dominio permanecen separados;
- los tests pueden construirse directamente a partir de estas reglas.

### Negativas

- expresiones regulares incorrectamente diseñadas por el usuario pueden generar detecciones demasiado amplias;
- los valores por defecto pueden necesitar ajustes después de observar workloads reales;
- la captura basada en cantidad de registros puede representar ventanas temporales muy distintas dependiendo de la frecuencia de logs.

## Fuera de alcance del MVP

- configuración individual de reglas built-in;
- prioridades entre múltiples `error_patterns`;
- diferentes `grouping_window` según tipo de error;
- límites de contexto dinámicos;
- detección basada en tasa de errores;
- anomaly detection;
- configuración de polling por Project;
- reglas generadas automáticamente por QVAC;
- actualización retroactiva de incidentes al modificar la configuración.
