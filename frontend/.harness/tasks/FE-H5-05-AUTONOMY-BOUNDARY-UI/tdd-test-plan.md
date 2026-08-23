# TDD Test Plan — FE-H5-05 Autonomy Boundary UI

## Scope

Verificar el renderizado y cumplimiento del límite de autonomía cuando la remediación alcanza el estado `pull_request_created` o termina su ejecución.

## Tests to add/update

Ubicación: `akritas/frontend/src/features/incidents/views/IncidentDetailView/__tests__/AutonomyBoundary.test.ts`.

1. **Test 1: Renderizado de AutonomyBoundaryBanner en `pull_request_created`**
   - Debe contener las etiquetas de seguridad: `Human Review Required`, `No Auto-Merge`, `No Auto-Deploy`.
   - Debe certificar que el flujo autónomo concluyó exitosamente.

2. **Test 2: Verificación de ausencia de acciones mutativas de producción**
   - Asegurar que no se expongan métodos o botones para merge, deploy o rollback.

## Expected failing tests before implementation

- Actualmente no existe el componente `AutonomyBoundaryBanner`.

## Acceptance criteria covered

- Todos los criterios de aceptación de `FE-H5-05` cubiertos.

