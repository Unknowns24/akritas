# FE-H5-05-AUTONOMY-BOUNDARY-UI — Autonomy Boundary UI

## Estado

pending

## Tipo de tarea

frontend-feature

## Modo de proyecto

existing_project

## Contexto

La arquitectura de seguridad de Akritas (ADR-004, ADR-006, ADR-008) establece que la autonomía del agente concluye de forma terminante tras la creación de la Pull Request en GitHub.

Bajo ninguna circunstancia la UI debe ofrecer ni sugerir botones o flujos de:
- Auto-merge de Pull Requests.
- Despliegue automático a producción.
- Rollback automatizado en infraestructura.
- Promoción de artefactos o triggers de release sin intervención humana.

Una vez que un incidente alcanza `pull_request_created` o un estado terminal, la interfaz debe mostrar el flujo como **completado**, documentar el límite de autonomía con un componente visual de seguridad y delegar el merge y release a los ingenieros humanos a través de las políticas de revisión de GitHub y CI/CD existente.

## Objetivo

Implementar el componente de frontera de autonomía (`AutonomyBoundaryBanner` / `AutonomyTerminalView`):
1. Tras la creación de la PR (`pull_request_created`), renderizar un banner / card de finalización autónoma con diseño de seguridad y auditoría.
2. Declarar explícitamente:
   - *"Autonomous workflow complete. Akritas operates up to Pull Request creation only."*
   - *"Code merge, deployment, and production promotion require human review and authorization through your standard CI/CD workflow."*
3. Auditar toda la interfaz de detalle de incidente para asegurar la **ausencia total** de controles interactivos de auto-merge, deploy o rollback automático.
4. Integrar en `RemediationCard.tsx` y `IncidentDetailView.tsx`.

## Criterios de aceptación

1. En estado `pull_request_created`, la UI presenta el flujo como terminado (badge "Completed", icono de candado de seguridad o escudo).
2. Se muestra el aviso explícito de frontera de autonomía (`AutonomyBoundaryBanner`).
3. Cero botones o APIs para merge, deploy o rollback automático en el frontend.
4. Pruebas unitarias que aseguren que el banner de frontera de autonomía se renderiza ante `pull_request_created`.
5. Verificaciones de TypeScript (`tsc --noEmit`) y ESLint limpias.

## Restricciones técnicas

- Stack: Next.js App Router, React 19, TypeScript, CSS Modules.
- Cumplimiento estricto con ADR-004.

