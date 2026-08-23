# Implementation Brief — FE-H5-01 Remediation Lifecycle UI

## Task

FE-H5-01-REMEDIATION-LIFECYCLE: Representar `Remediation` únicamente cuando `resolution_status = fixable` y cubrir sus estados (`planned`, `in_progress`, `validated`, `failed`, `pull_request_created`) según el contrato OpenAPI.

## Current project context

- `akritas/frontend` contiene la vista `IncidentDetailView` (`src/features/incidents/views/IncidentDetailView/`).
- Actualmente, `RemediationCard.tsx` asume de forma básica que hay cambios diffs y un botón estático de PR, sin validar si `resolution_status` es `fixable` ni reflejar los estados de ciclo de vida (`planned`, `in_progress`, `validated`, `failed`, `pull_request_created`).
- El contrato `backend/docs/openapi.yaml` y los tipos generados en `@/core/libs/api-client` ya incluyen la especificación completa de `Remediation`, `RemediationStatus`, `CodeChange`, `ValidationSummary` y `PullRequestReference`.

## Proposed approach

1. **Condición de entrada en `IncidentDetailView`**:
   - Evaluar `incident.resolution_status`.
   - Si `resolution_status === 'requires_human'`, mostrar `RequiresHumanCard` (explicación de que el incidente requiere intervención humana tras la Issue y no posee remediación automática).
   - Si `resolution_status === 'fixable'`, renderizar `RemediationCard` con el lifecycle completo.
   - Si `resolution_status` es indefinido (ej. investigación no concluida), mostrar estado de espera / investigación en progreso.

2. **Diseño de Subcomponentes Modulares en `src/features/incidents/views/IncidentDetailView/components/`**:
   - `RemediationCard.tsx`: Orquestador principal del componente de remediación.
   - `RemediationStatusBadge.tsx`: Badge con estilo e icono según el estado (`planned`, `in_progress`, `validated`, `failed`, `pull_request_created`).
   - `ValidationSummaryView.tsx`: Indicadores visuales de validación (checks pasados vs fallidos).
   - `CodeChangesDiffViewer.tsx`: Visor estructurado de parches (`patch`) para múltiples archivos con diferenciación de líneas agregadas/eliminadas.
   - `RequiresHumanCard.tsx`: Visualización para casos no solucionables automáticamente.

3. **Mapeo de Estados según OpenAPI**:
   - `planned`: Muestra rama objetivo/planificada, indicador de cola, descripción de cambios planeados.
   - `in_progress`: Spinner/pulso activo, rama en preparación, validaciones en curso.
   - `validated`: Validaciones superadas (badge verde con checks totales y pasados), visor de diffs listo.
   - `failed`: Banner de error, mensaje `failure_user_message`, desglose de tests fallidos y aviso explícito: *No Pull Request Created*.
   - `pull_request_created`: Badge de PR creada, botón con enlace externo a GitHub (`pull_request_reference.url`), rama (`pull_request_reference.branch`) y disclaimer de seguridad *Akritas never merges changes automatically*.

## Architecture impact

- Respeta la arquitectura hexagonal y de features (`src/features/incidents/`).
- Sin dependencias de kits externos no autorizados.
- Uso de CSS Modules con variables globales de diseño (`IncidentDetailView.module.css` / `RemediationCard.module.css`).

## API/OpenAPI impact

- Consume directamente `components["schemas"]["Incident"]`, `components["schemas"]["Remediation"]`, `components["schemas"]["RemediationStatus"]`.
- Cero desvíos de contrato.

## Data/persistence impact

- Ninguno (exclusivo de UI).

## Error handling impact

- Maneja de forma robusta propiedades opcionales (`changes`, `validation_summary`, `failure_user_message`, `pull_request_reference`).

## Test strategy

- Pruebas unitarias en componentes para validar:
  1. Renderizado condicional por `resolution_status` (`fixable` vs `requires_human` vs pendiente).
  2. Renderizado de cada uno de los 5 estados de `RemediationStatus`.
  3. Formato y resaltado del visor de diffs.
- Verificación mediante `npm run typecheck` y `npm run lint`.

## Risks

- Propiedades `changes` o `validation_summary` nulas en estados tempranos (`planned`). Solución: defensas y valores por defecto seguros en componentes.

## Files likely to change

- `akritas/frontend/src/features/incidents/views/IncidentDetailView/components/RemediationCard.tsx`
- `akritas/frontend/src/features/incidents/views/IncidentDetailView/components/RequiresHumanCard.tsx` (Nuevo)
- `akritas/frontend/src/features/incidents/views/IncidentDetailView/components/RemediationStatusBadge.tsx` (Nuevo)
- `akritas/frontend/src/features/incidents/views/IncidentDetailView/components/ValidationSummaryView.tsx` (Nuevo)
- `akritas/frontend/src/features/incidents/views/IncidentDetailView/components/CodeChangesDiffViewer.tsx` (Nuevo)
- `akritas/frontend/src/features/incidents/views/IncidentDetailView/IncidentDetailView.module.css`
- `akritas/frontend/src/features/incidents/views/IncidentDetailView/IncidentDetailView.tsx`
- `akritas/frontend/src/features/incidents/views/IncidentDetailView/__tests__/RemediationCard.test.tsx` (Nuevo)

