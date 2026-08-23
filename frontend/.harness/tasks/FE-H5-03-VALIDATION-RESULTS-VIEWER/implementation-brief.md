# Implementation Brief — FE-H5-03 Validation Results Viewer

## Task

FE-H5-03-VALIDATION-RESULTS-VIEWER: Visualizar los resultados detallados de validación (`tests`, `build`, `static checks`), sus estados individuales, resúmenes y extractos de salida (`output_excerpt`), mostrando de forma inequívoca el bloqueo de Pull Request ante cualquier validación fallida.

## Current project context

- `FE-H5-01` introdujo el resumen de validaciones numérico (`ValidationSummaryView`) en `RemediationCard`.
- `FE-H5-03` expande esa capacidad para mostrar los resultados individuales y la evidencia técnica completa de cada check (`ValidationResult`).
- El endpoint OpenAPI `GET /remediations/{remediation_id}/validation-results` ya está documentado y expone `ValidationResult[]`.

## Proposed approach

1. **Servicio API en `src/features/incidents/services/`**:
   - `get-remediation-validation-results.service.ts`: Llama a `GET /remediations/{remediation_id}/validation-results` con paginación cursor.

2. **Subcomponentes en `src/features/incidents/views/IncidentDetailView/components/`**:
   - `ValidationResultItem.tsx`: Renders a single check (name, type icon, status badge, duration, summary, collapsible output excerpt).
   - `ValidationResultsViewer.tsx`: Contenedor de la lista de checks con switch/acordeón para desplegar la evidencia de ejecución.
   - `ValidationFailureBanner.tsx`: Alerta prominente cuando `hasFailed = true` confirmando que la remediación falló y la PR no fue creada.

3. **Integración en `RemediationCard`**:
   - Conectar `ValidationResultsViewer` dentro de `RemediationCard` para mostrar la traza de validaciones cuando la remediación está en estado `validated`, `failed`, `in_progress` o `pull_request_created`.

## Architecture impact

- Feature-based en `src/features/incidents/`.
- Tipos de datos estrictos desde OpenAPI `@/core/libs/api-client`.

## Test strategy

- Tests unitarios para `ValidationResultsViewer` y `ValidationResultItem` evaluando estados `passed`, `failed`, `running`, `pending`, y renderizado de `output_excerpt`.

