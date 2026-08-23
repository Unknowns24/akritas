# FE-H5-04-PR-REFERENCE-TRACEABILITY — Pull Request Reference y Traceability

## Estado

pending

## Tipo de tarea

frontend-feature

## Modo de proyecto

existing_project

## Contexto

En el Hito 5 de Akritas, la remediación genera trazabilidad de extremo a extremo desde el momento de la detección hasta la creación de la Pull Request.

Para auditoría, gobernanza y transparencia técnica, el usuario debe poder ver la cadena completa de trazabilidad:
`Incident → Investigation → Issue → Remediation → Branch → Commit → Pull Request`

Cada nodo de esta cadena debe exponer sus identificadores, estados y enlaces externos correspondientes (ej. Issue URL, PR URL, nombres de rama y hashes de commit).

## Objetivo

Implementar el visualizador de trazabilidad completa (`TraceabilityTimelineCard` / `TraceabilityChainView`):
1. Representación visual de la cadena lineal de 7 etapas:
   - **Incident**: Key (`INC-XXX`), severidad y timestamp.
   - **Investigation**: ID de investigación, categoría de causa raíz y estado.
   - **GitHub Issue**: Número `#<num>`, repositorio y enlace externo a GitHub.
   - **Remediation**: ID de remediación y estado (`planned`, `in_progress`, `validated`, `failed`, `pull_request_created`).
   - **Branch**: Nombre de la rama (`akritas/fix-...`).
   - **Commit**: Hash / SHA del commit de fix o correlación de deployment.
   - **Pull Request**: Número `#<num>`, enlace directo a GitHub, y disclaimer inmutable de seguridad.
2. Manejo de estados incompletos o bifurcaciones:
   - Si `resolution_status = requires_human`: la cadena se detiene en `Issue` con un badge de "Manual Intervention Required".
   - Si `remediation.status = failed`: la cadena se detiene en `Remediation` con badge de "Validation Failed — PR Blocked".
   - Si `pull_request_created`: la cadena muestra los 7 nodos conectados en verde/éxito.
3. Integración en `IncidentDetailView`.

## Criterios de aceptación

1. La cadena de trazabilidad renderiza los 7 eslabones con datos reales del `Incident`.
2. Los enlaces a GitHub (Issue y PR) abren en pestaña nueva con `target="_blank"` y `rel="noreferrer"`.
3. Se representan fielmente las interrupciones de seguridad (cuando no hay PR por fallo o por intervención humana).
4. Pruebas unitarias de cálculo de cadena de trazabilidad.
5. Verificaciones de TypeScript y ESLint limpias.

## Restricciones técnicas

- Stack: Next.js App Router, React 19, TypeScript, CSS Modules.
- Tipos de datos tomados exclusivamente del cliente OpenAPI.

