# TDD Test Plan — FE-H5-01 Remediation Lifecycle UI

## Scope

Verificar el comportamiento y renderizado de la UI de remediación para incidentes en todas las variantes de `resolution_status` y `remediation.status` según OpenAPI.

## Tests to add/update

Ubicación: `akritas/frontend/src/features/incidents/views/IncidentDetailView/__tests__/RemediationCard.test.tsx` (o archivo de pruebas unitarias correspondiente).

1. **Test 1: Renderizado condicional por `resolution_status`**
   - Cuando `resolution_status = "requires_human"`:
     - No debe renderizar el card de remediación activa ni el botón de PR.
     - Debe renderizar el mensaje explicativo de que el incidente requiere intervención humana.
   - Cuando `resolution_status = "fixable"`:
     - Debe renderizar el componente de `Remediation`.
   - Cuando `resolution_status` no está presente (ej. investigación pendiente):
     - Debe indicar que la evaluación de remediación aguarda el resultado de la investigación.

2. **Test 2: Estado `planned`**
   - Debe renderizar el badge de estado "Planned".
   - Debe mostrar la rama planificada (ej. `branch_name` o fallback derivado del incidente).
   - Debe mostrar mensaje de remediación en cola.

3. **Test 3: Estado `in_progress`**
   - Debe renderizar el badge de estado "In Progress" con indicador de actividad.
   - Debe indicar que los cambios de código o las validaciones se están ejecutando.

4. **Test 4: Estado `validated`**
   - Debe renderizar el badge de estado "Validated".
   - Debe mostrar el resumen de validaciones pasadas (ej. `Validation Passed: X/X checks`).
   - Debe mostrar el visor de diffs con los archivos modificados.

5. **Test 5: Estado `failed`**
   - Debe renderizar el badge de estado "Failed".
   - Debe mostrar el mensaje de error `failure_user_message`.
   - Debe mostrar el resumen de validaciones con conteo de fallos.
   - Debe mostrar el aviso de que **no se ha creado Pull Request**.

6. **Test 6: Estado `pull_request_created`**
   - Debe renderizar el badge "Pull Request Created".
   - Debe mostrar el enlace externo a la Pull Request (`pull_request_reference.url`) con el número `#<number>`.
   - Debe mostrar la rama `pull_request_reference.branch`.
   - Debe mostrar el disclaimer: *"Akritas never merges changes automatically"*.

7. **Test 7: Visor de Diffs (`CodeChangesDiffViewer`)**
   - Debe renderizar la lista de archivos (`file_path` y `change_type`).
   - Debe parsear y colorear adecuadamente las líneas con `+`, `-`, y `@@`.

## Expected failing tests before implementation

- Componentes actuales no contemplan `resolution_status !== 'fixable'` de forma diferenciada.
- `RemediationCard` actual no tiene soporte para `planned`, `in_progress`, `failed`, `validated`, ni muestra `failure_user_message` o badges de estado.

## Acceptance criteria covered

- Todos los criterios de aceptación de `FE-H5-01` cubiertos.

## Open questions / human approval notes

- Todo está alineado con OpenAPI 3.1.0 (`backend/docs/openapi.yaml`) y ADRs 002, 004 y 006.
