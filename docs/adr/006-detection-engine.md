# ADR-006 — Detection Engine determinístico basado en reglas y fingerprints estables

## Estado

Accepted

## Contexto

ADR-003 establece que la detección de logs debe ocurrir antes de la investigación con IA y que QVAC no debe procesar individualmente el stream completo de logs.

Para implementar esa separación es necesario definir cómo Akritas transforma logs provenientes de Dokploy en `LogEvent` e `Incident` sin depender de inferencia probabilística.

Los logs contienen ruido, mensajes multilinea, identificadores variables, timestamps y valores específicos de cada request. Comparar líneas crudas produciría falsos duplicados o múltiples incidentes para manifestaciones equivalentes del mismo problema.

Al mismo tiempo, el Detection Engine del MVP debe ser explicable, reproducible y suficientemente simple para funcionar de manera confiable durante la demo.

## Decisión

El Detection Engine será un pipeline determinístico basado en reglas explícitas.

QVAC no participa en ninguna etapa de detección, normalización, fingerprinting o agrupación.

El pipeline conceptual será:

```text
Dokploy logs
    ↓
Cursor / deduplication
    ↓
Event assembly
    ↓
Ignored rules
    ↓
Detection rules
    ↓
Normalization
    ↓
Fingerprint
    ↓
Incident grouping
    ↓
Context capture
    ↓
Incident ready for investigation
```

### 1. Cursor y deduplicación

Cada Project mantendrá un cursor temporal o mecanismo equivalente que permita solicitar únicamente logs posteriores al último punto procesado.

Si la fuente puede devolver líneas repetidas entre consultas, el engine debe evitar procesar dos veces la misma ocurrencia.

El cursor pertenece al estado operativo del monitoring y no modifica el modelo conceptual del incidente.

### 2. Event assembly

Antes de aplicar reglas, Akritas debe reconstruir eventos lógicos cuando un error ocupa múltiples líneas.

Esto permite tratar como una única unidad señales como:

- stack traces;
- panic traces;
- excepciones con frames sucesivos;
- mensajes de error continuados en varias líneas.

El resultado de esta etapa es un evento lógico candidato, no necesariamente un `Incident`.

### 3. Ignored rules

Las reglas de exclusión se evalúan antes de crear un `LogEvent`.

`MonitoringConfiguration.ignored_patterns` permite descartar mensajes conocidos o irrelevantes para un Project.

Un evento ignorado no crea `LogEvent`, no genera fingerprint y no invoca QVAC.

### 4. Detection rules

Un evento se considera relevante cuando cumple al menos una regla determinística habilitada.

El MVP incluirá reglas built-in para señales reconocibles como:

- niveles `ERROR`, `FATAL` o equivalentes;
- `PANIC` o process crash;
- stack traces o excepciones no controladas;
- respuestas HTTP 5xx cuando puedan identificarse de forma confiable en el log;
- señales de restart/crash emitidas por la fuente de logs;
- patrones adicionales definidos en `MonitoringConfiguration.error_patterns`.

Cada detección debe registrar qué regla o reglas la activaron para mantener trazabilidad.

No se utilizará un modelo probabilístico ni un LLM para decidir si un evento es relevante.

### 5. Normalización

Antes de construir el fingerprint, Akritas generará una representación estable del error eliminando valores que normalmente cambian entre ocurrencias pero no representan una causa diferente.

La normalización puede remover o reemplazar elementos como:

- timestamps;
- request IDs y correlation IDs;
- UUIDs;
- identificadores numéricos claramente variables;
- direcciones de memoria;
- otros valores dinámicos reconocidos por reglas explícitas.

La normalización debe ser conservadora: ante la duda se preserva información para evitar agrupar errores distintos de manera incorrecta.

El contenido original permanece disponible como evidencia; la normalización se utiliza para fingerprinting y agrupación.

### 6. Fingerprint

El fingerprint será determinístico y estable para eventos equivalentes dentro de un Project.

Se derivará de una firma compuesta por información estable disponible, por ejemplo:

```text
project
+ detection rule / error type
+ normalized primary message
+ stable stack-trace location when available
```

La firma se convertirá en un identificador mediante una función hash estable.

Dos ocurrencias con el mismo problema pero valores dinámicos distintos deben tender a producir el mismo fingerprint.

Errores provenientes de Projects distintos nunca deben compartir un Incident aunque su mensaje sea idéntico.

### 7. Incident grouping

Una vez calculado el fingerprint:

- si existe un Incident compatible y abierto dentro de la política de agrupación, se agrega la nueva ocurrencia;
- si no existe, se crea un nuevo Incident;
- `occurrence_count` se incrementa por cada ocurrencia aceptada;
- `first_seen_at` conserva la primera ocurrencia;
- `last_seen_at` se actualiza con la más reciente.

`MonitoringConfiguration.grouping_window` controla la ventana temporal utilizada por la política de agrupación.

Un `LogEvent` representa la ocurrencia relevante detectada. El `Incident` representa el problema agrupado.

### 8. Context capture

Para cada ocurrencia relevante, Akritas conservará una cantidad acotada de logs anteriores y posteriores según:

- `MonitoringConfiguration.context_before`;
- `MonitoringConfiguration.context_after`.

Este contexto forma parte de la evidencia disponible para la investigación, pero no participa necesariamente del fingerprint.

## Invariantes

- la misma entrada y la misma configuración deben producir la misma decisión de detección;
- QVAC nunca decide si una línea de log debe convertirse en `LogEvent`;
- QVAC nunca genera fingerprints;
- los patrones ignorados se aplican antes de crear incidentes;
- el fingerprint se construye a partir de datos normalizados, no del contexto completo;
- un evento descartado no debe disparar una Investigation;
- el Detection Engine debe poder explicar qué regla produjo cada detección.

## Consecuencias

### Positivas

- comportamiento reproducible y fácil de demostrar;
- menor cantidad de ruido enviado a QVAC;
- menor uso de CPU y contexto del modelo local;
- agrupación consistente de ocurrencias equivalentes;
- detecciones auditables y explicables;
- posibilidad de testear el engine con fixtures sin ejecutar QVAC ni acceder a GitHub.

### Negativas

- las reglas built-in deben mantenerse explícitamente;
- una normalización demasiado agresiva puede agrupar errores distintos;
- una normalización demasiado conservadora puede fragmentar un mismo incidente;
- errores sin señales conocidas pueden no detectarse;
- detectar anomalías puramente estadísticas queda limitado en el MVP.

## Fuera de alcance del MVP

- detección basada en embeddings;
- clasificación de logs mediante LLM;
- anomaly detection estadístico generalizado;
- correlación entre múltiples fuentes de observabilidad;
- fingerprints semánticos aprendidos;
- reglas auto-generadas por QVAC.
