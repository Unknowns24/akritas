# FE-H5-06-VISUAL-GATE-H5 — Gate Visual H5: Golden Flow de Remediación

## Estado

pending

## Tipo de tarea

frontend-feature

## Modo de proyecto

existing_project

## Contexto

El Hito 5 culmina con el **Gate Visual de Remediación**, donde se debe demostrar y verificar visualmente la secuencia completa:
`Issue → fixable → branch → regression test → fix → tests pass → PR`

Este flujo debe proporcionar **contexto suficiente para la revisión humana**, incluyendo:
- Identificación de la Issue de origen con enlace a GitHub.
- Clasificación de causa raíz y justificación de por qué es `fixable`.
- Rama dedicada de fix.
- Test de regresión generado para reproducir el incidente.
- Parche de código propuesto con visor de diffs.
- Verificación de validaciones superadas con evidencias y trazas de ejecución.
- Pull Request generada con su commit referenciado y frontera de seguridad inmutable.

## Objetivo

Implementar el visualizador integral del Gate Visual de Hito 5 (`RemediationHumanReviewPacket` / `Milestone5VisualGateView`):
1. Panel interactivo de revisión humana que resume de manera estructurada los 7 hitos del Golden Flow.
2. Exhibición del contexto técnico completo:
   - Root Cause & Severity.
   - Regression Test (`remediation.regression_test`).
   - Code Diff (`remediation.changes`).
   - Validation Logs (`ValidationResultsViewer`).
   - Commit & Pull Request link.
3. Integración en la vista de detalle de incidente.

## Criterios de aceptación

1. El visualizador resume de manera fluida y comprensible la secuencia completa `Issue → fixable → branch → regression test → fix → tests pass → PR`.
2. Proporciona todo el contexto necesario para que un ingeniero humano tome una decisión informada sobre la Pull Request sin salir del contexto.
3. Pruebas unitarias de los componentes del Golden Flow.
4. Verificación de compilación de TypeScript y ESLint limpia.

## Restricciones técnicas

- Stack: Next.js App Router, React 19, TypeScript, CSS Modules.
- Apego a los contratos OpenAPI y ADRs de seguridad (ADR-002, ADR-004, ADR-006).

