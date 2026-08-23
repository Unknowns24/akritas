# FE-H5-01-REMEDIATION-LIFECYCLE — Remediation Lifecycle UI

## Estado

pending

## Tipo de tarea

frontend-feature

## Modo de proyecto

existing_project

## Contexto

El MVP de Akritas implementa el ciclo autónomo de remediación en el Hito 5. Según la especificación, el ciclo de remediación sólo debe iniciarse e interactuarse cuando el incidente fue investigado y clasificado con `resolution_status = fixable`.

Para cualquier otro caso (`resolution_status = requires_human`), el flujo automático finaliza tras la publicación de la GitHub Issue (ADR-002, ADR-004), por lo que la UI debe reflejar de forma explícita este límite sin presentar acciones o estados de remediación erróneos.

El contrato OpenAPI (`backend/docs/openapi.yaml`) define `Remediation` y `RemediationStatus` con los estados:
- `planned`
- `in_progress` (running)
- `validated`
- `failed`
- `pull_request_created` (completed)

## Objetivo

Implementar la representación visual completa y tipada del ciclo de vida de `Remediation` en `IncidentDetailView` respetando:
1. Condición estricta: `Remediation` se representa **únicamente para `resolution_status = fixable`**.
2. Representación de los 5 estados según OpenAPI (`planned`, `in_progress`, `validated`, `failed`, `pull_request_created`).
3. Si `resolution_status = requires_human`, mostrar el estado informativo correspondiente sin flujo de remediación activa.
4. Mostrar diffs de cambios (`changes: CodeChange[]`), resumen de validación (`validation_summary`), mensajes de fallo (`failure_user_message`), branch y enlace a Pull Request cuando aplique.
5. Mantener visible el principio de seguridad: *Akritas never merges changes automatically*.

## Criterios de aceptación

1. Si `incident.resolution_status === 'requires_human'`, no se muestra un flujo de remediación activa; se renderiza un card/indicador que explica que el incidente requiere resolución manual y no admite remediación automática.
2. Si `incident.resolution_status === 'fixable'`:
   - Estado `planned`: Badge informativo, branch asignada/planeada, estado pendiente de inicio.
   - Estado `in_progress`: Indicador de ejecución activa, validaciones en curso.
   - Estado `validated`: Badge de validación superada, visor de diffs formateado, resumen de checks pasados.
   - Estado `failed`: Badge de error, mensaje de fallo al usuario (`failure_user_message`), resumen de checks fallidos, aviso explícito de que **no se creó Pull Request**.
   - Estado `pull_request_created`: Badge de éxito, enlace directo a la Pull Request en GitHub (`pull_request_reference.url`), nombre de la rama y disclaimer de seguridad.
3. El visor de diffs maneja múltiples archivos y colorea líneas añadidas (`+`), eliminadas (`-`) y cabeceras (`@@`).
4. Typecheck (`npm run typecheck` / `tsc --noEmit`) y lint pasan sin advertencias ni errores.

## Restricciones técnicas

- Stack: Next.js 16 App Router, React 19, TypeScript, CSS Modules.
- Sin librerías UI externas adicionales (Tailwind, Radix wrappers, shadcn).
- Consumo tipado del OpenAPI client (`@/core/libs/api-client`).
- Código modular en `src/features/incidents/`.

## Instrucción para el harness

Generar `implementation-brief.md` y `tdd-test-plan.md`. Organizar la ejecución en actividades atómicas con commits individuales.

