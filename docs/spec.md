# Akritas — Product Specification

## 1. Problema

Cuando una aplicación falla en producción, un ingeniero normalmente debe:

1. detectar el error;
2. revisar logs;
3. identificar el servicio y flujo afectados;
4. inspeccionar stack traces;
5. revisar el código relevante;
6. correlacionar el error con cambios recientes;
7. formular una causa probable;
8. registrar el incidente;
9. implementar y validar una corrección si corresponde.

Este proceso es costoso, repetitivo y exige acceso a información sensible como logs de producción, código fuente, arquitectura interna y potencialmente secretos o datos de clientes.

## 2. Propuesta

Akritas automatiza ese proceso utilizando un agente local basado en QVAC.

Cada proyecto monitoreado vincula al menos:

- un repositorio GitHub;
- una aplicación/contenedor desplegado mediante Dokploy.

Akritas observa la ejecución de la aplicación, detecta eventos que puedan constituir un incidente, investiga el contexto y utiliza el repositorio para entender la posible causa.

## 3. Flujo principal

```text
Dokploy Logs
    ↓
Detection Engine
    ↓
Incident Candidate
    ↓
Deduplication / Grouping
    ↓
Incident
    ↓
QVAC Investigation
    ↓
Repository Inspection
    ↓
Root Cause Assessment
    ↓
GitHub Issue (always)
    ↓
Fixable?
 ┌───────┴────────┐
 No              Yes
 │                ↓
End          Generate Fix
                  ↓
             Validate Fix
                  ↓
             Pull Request
```

## 4. Configuración por proyecto

Cada Project debe permitir configurar:

### Identidad
- nombre;
- descripción opcional;
- estado de monitoreo.

### GitHub
- repositorio;
- branch por defecto;
- credencial/autorización necesaria para lectura y escritura;
- permisos para crear Issues y Pull Requests.

### Dokploy
- instancia/servidor;
- aplicación asociada;
- credencial de acceso;
- configuración mínima para obtener logs.

### Monitoring
- activado/desactivado;
- patrones adicionales de error opcionales;
- reglas de exclusión opcionales;
- ventana de agrupación de eventos;
- cantidad de contexto previa/posterior al error.

## 5. Detección de incidentes

QVAC no debe ser responsable de decidir sobre cada línea individual de logs.

Debe existir un Detection Engine convencional que identifique eventos potencialmente importantes utilizando señales como:

- `ERROR`;
- `FATAL`;
- `PANIC`;
- stack traces;
- HTTP 5xx;
- container restarts;
- process crashes;
- patrones configurados por proyecto;
- repetición anormal del mismo error.

Los eventos detectados se convierten en candidatos y posteriormente se agrupan para evitar crear múltiples incidentes por el mismo problema.

## 6. Investigación

Una investigación debe poder usar como contexto:

- logs que originaron el incidente;
- líneas anteriores y posteriores;
- frecuencia del error;
- timestamps;
- stack trace;
- aplicación afectada;
- repositorio asociado;
- archivos relevantes;
- búsquedas de código;
- commits recientes;
- diffs recientes;
- branch por defecto.

El objetivo de la investigación es producir:

- descripción del incidente;
- impacto observado;
- evidencia;
- causa raíz o hipótesis;
- nivel de confianza;
- archivos o componentes relevantes;
- evaluación de si el incidente es solucionable automáticamente.

## 7. GitHub Issue

Toda investigación debe terminar con la creación de una Issue, sin importar si Akritas puede solucionar el problema.

La Issue debe incluir al menos:

- título claro;
- aplicación/proyecto afectado;
- descripción del incidente;
- evidencia relevante;
- cantidad de ocurrencias cuando esté disponible;
- momento de primera y última detección;
- root cause status;
- posible causa raíz;
- confidence score;
- archivos o componentes involucrados;
- estado de remediación automática.

## 8. Remediación automática

Si Akritas determina que el problema es solucionable de manera suficientemente segura:

1. crea una branch específica;
2. realiza los cambios necesarios;
3. agrega o modifica tests cuando corresponda;
4. ejecuta validaciones disponibles;
5. crea una Pull Request;
6. referencia la Issue correspondiente.

La Pull Request debe explicar:

- qué incidente resuelve;
- cuál fue la causa identificada;
- qué cambios realizó;
- qué validaciones ejecutó;
- cualquier riesgo o limitación conocida.

## 9. Clasificaciones importantes

### Root Cause Status

- `identified`: Akritas considera identificada la causa raíz.
- `suspected`: existe una hipótesis probable pero no concluyente.
- `unknown`: no existe evidencia suficiente.

### Resolution Status

- `fixable`: Akritas puede proponer una corrección automática.
- `requires_human`: requiere intervención humana.

Estos estados son independientes.

Ejemplo:

```text
root_cause_status = identified
resolution_status = requires_human
```

Un proveedor externo caído puede tener una causa perfectamente identificada, pero no ser solucionable modificando el repositorio.

## 10. Fuera de alcance del MVP

- auto-merge de Pull Requests;
- auto-deploy;
- rollback automático;
- Slack/Discord/Teams;
- Grafana/Loki;
- Kubernetes;
- traces distribuidos;
- análisis avanzado de métricas;
- múltiples agentes especializados;
- remediación directa de infraestructura;
- integración con proveedores cloud adicionales.
