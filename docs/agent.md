# Akritas — QVAC Agent

## Rol

El agente QVAC es responsable de transformar un Incident y su evidencia en una investigación estructurada y, cuando corresponda, una propuesta de remediación.

QVAC es el único motor de inferencia de IA del MVP.

## Lo que el agente sí hace

- interpreta contexto de logs;
- entiende stack traces;
- formula hipótesis;
- decide qué información adicional necesita;
- solicita herramientas de lectura;
- inspecciona código;
- inspecciona commits/diffs;
- determina una causa raíz o hipótesis;
- clasifica el incidente;
- redacta el contenido de la Issue;
- propone cambios de código;
- genera una explicación para la Pull Request.

## Lo que el agente no hace directamente

- leer continuamente todos los logs;
- decidir mediante LLM si cada línea constituye un error;
- almacenar secretos;
- mergear Pull Requests;
- desplegar código;
- modificar producción;
- ejecutar comandos arbitrarios sin una tool explícitamente permitida.

## Tool Use

El agente interactúa con el sistema exclusivamente mediante herramientas explícitas.

Conjunto conceptual inicial:

### Observability

- `get_incident_logs`
- `get_surrounding_logs`

### Repository

- `search_code`
- `read_file`
- `list_recent_commits`
- `read_commit`
- `read_diff`

### GitHub Actions

- `create_issue`
- `create_branch`
- `commit_changes`
- `create_pull_request`

### Validation

- `run_tests`
- `run_static_checks`

El nombre exacto de las tools es una decisión de implementación.

## Structured Output

Las decisiones importantes no deben depender de texto libre.

El agente debe producir estructuras equivalentes a:

```json
{
  "root_cause_status": "identified | suspected | unknown",
  "resolution_status": "fixable | requires_human",
  "confidence": 0.0,
  "summary": "...",
  "root_cause": "...",
  "evidence": [],
  "relevant_files": [],
  "recommended_actions": []
}
```

## Reliability

Un modelo pequeño local debe trabajar con contexto reducido y herramientas específicas.

Akritas debe priorizar:

- outputs estructurados;
- contexto incremental;
- herramientas pequeñas;
- límites de iteración;
- validación determinística;
- separación entre análisis y acciones mutativas.

## Seguridad

Toda acción de escritura debe ser explícita y auditable.

Una conclusión de `fixable` no implica que el cambio sea válido: la remediación debe superar las validaciones configuradas antes de crear una Pull Request.
