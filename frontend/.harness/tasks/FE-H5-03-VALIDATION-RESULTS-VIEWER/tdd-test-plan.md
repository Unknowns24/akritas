# TDD Test Plan — FE-H5-03 Validation Results Viewer

## Scope

Verificar el servicio, tipos y renderizado de la lista detallada de resultados de validación (`ValidationResult[]`), su desglose por tipo (`test`, `build`, `static_analysis`), y el banner explícito de bloqueo de PR en caso de fallo.

## Tests to add/update

Ubicación: `akritas/frontend/src/features/incidents/views/IncidentDetailView/__tests__/ValidationResultsViewer.test.tsx`.

1. **Test 1: Renderizado de lista de checks exitosos**
   - Debe listar cada check con su ícono (`test`, `build`, `static_analysis`), nombre (`name`), y badge `passed`.
   - Debe permitir desplegar el `output_excerpt` con la traza de ejecución.

2. **Test 2: Renderizado de check fallido y banner de bloqueo de PR**
   - Ante un check con `status = "failed"`, debe renderizar el badge de error en rojo.
   - Debe renderizar el banner explícito: *"Remediation Failed — No Pull Request was created"*.
   - Debe mostrar la salida de error (`output_excerpt`) con estilo de consola.

3. **Test 3: Checks en ejecución (`running`) o en espera (`pending`)**
   - Debe mostrar spinner de progreso en checks activos.

## Expected failing tests before implementation

- Actualmente no existe el componente `ValidationResultsViewer` ni el servicio `get-remediation-validation-results.service.ts`.

## Acceptance criteria covered

- Todos los criterios de aceptación de `FE-H5-03` cubiertos.

