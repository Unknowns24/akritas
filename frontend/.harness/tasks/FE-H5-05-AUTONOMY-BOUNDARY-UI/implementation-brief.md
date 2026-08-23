# Implementation Brief — FE-H5-05 Autonomy Boundary UI

## Task

FE-H5-05-AUTONOMY-BOUNDARY-UI: Mostrar el flujo como terminado tras la creación de la Pull Request, sin ofrecer acciones automáticas de merge, deploy, rollback ni promoción productiva, reforzando visualmente el límite estricto de autonomía (ADR-004).

## Current project context

- `RemediationCard` muestra el botón de enlace externo a GitHub PR (`View Pull Request #X`).
- Se necesita un componente explícito y estilizado (`AutonomyBoundaryBanner`) que certifique que el flujo autónomo concluyó satisfactoriamente y que las acciones posteriores (revisión, merge, deploy) corresponden exclusivamente al equipo humano.

## Proposed approach

1. **Componente `AutonomyBoundaryBanner`**:
   - Iconografía de seguridad (`ShieldCheck`, `Lock`).
   - Título: *"Autonomous Remediation Completed"*.
   - Explicación de gobernanza:
     - *"Akritas has completed all autonomous remediation activities for this incident."*
     - *"Per safety design (ADR-004), merge, deployment, and rollback actions are never performed automatically. Please review and merge the Pull Request on GitHub according to your team's code review policies."*
   - Badges visuales de seguridad: `No Auto-Merge`, `No Auto-Deploy`, `Human Review Required`.

2. **Integración en `RemediationCard`**:
   - Renderizar `AutonomyBoundaryBanner` cuando `status === 'pull_request_created'`.
   - Mantener el footer de seguridad en todos los estados.

3. **Auditoría de Componentes**:
   - Verificar que ningún componente de `IncidentDetailView` exponga botones de mutación productiva.

## Architecture impact

- Feature `incidents` modular en `src/features/incidents/views/IncidentDetailView/components/`.

## Test strategy

- Tests unitarios para verificar que ante `pull_request_created`, el banner de frontera de autonomía se renderiza y no existen botones de merge ni deploy.

