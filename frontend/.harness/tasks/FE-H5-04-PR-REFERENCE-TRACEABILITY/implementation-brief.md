# Implementation Brief — FE-H5-04 Pull Request Reference y Traceability

## Task

FE-H5-04-PR-REFERENCE-TRACEABILITY: Implementar el visualizador interactivo de trazabilidad completa `Incident → Investigation → Issue → Remediation → Branch → Commit → Pull Request` con enlaces, referencias y estados contextuales.

## Current project context

- `IncidentDetailView` contiene cards individuales para Root Cause, Stack Trace, Context y Remediation.
- Se requiere una vista unificada de trazabilidad horizontal/vertical que conecte todos los artefactos de la cadena desde el incidente original hasta la PR creada.

## Proposed approach

1. **Modelo de Trazabilidad y Utilidades (`traceability.utils.ts`)**:
   - Función extractora `buildIncidentTraceabilityChain(incident: Incident): TraceabilityStep[]`.
   - Modela los 7 pasos con su ID, título, subtítulo, estado (`completed`, `active`, `halted`, `pending`), link externo si aplica, y metadata.

2. **Componente Visual `TraceabilityChainView.tsx` y `TraceabilityChainView.module.css`**:
   - Renderiza un stepper/flow gráfico con conectores visuales entre cada etapa.
   - Cada nodo muestra su ícono específico:
     - `Incident`: `AlertCircle`
     - `Investigation`: `Search`
     - `Issue`: `Bookmark` / `GitPullRequestDraft`
     - `Remediation`: `Wrench`
     - `Branch`: `GitBranch`
     - `Commit`: `GitCommit`
     - `Pull Request`: `GitPullRequest`
   - Si el flujo se detiene en un paso (ej. `requires_human` en Issue o `failed` en Remediation), el conector y los pasos siguientes se marcan como inactivos o bloqueados con un badge explicativo.

3. **Integración en `IncidentDetailView.tsx`**:
   - Se agrega el componente `TraceabilityChainView` en la vista de detalle del incidente.

## Architecture impact

- Feature `incidents` modular con CSS Modules.

## Test strategy

- Tests unitarios para `buildIncidentTraceabilityChain`:
  - Cadena completa exitosa con PR.
  - Cadena detenida en `requires_human`.
  - Cadena detenida en `failed` de validación.

