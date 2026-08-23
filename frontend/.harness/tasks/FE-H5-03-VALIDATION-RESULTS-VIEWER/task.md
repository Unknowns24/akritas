# FE-H5-03-VALIDATION-RESULTS-VIEWER — Validation Results Viewer

## Estado

pending

## Tipo de tarea

frontend-feature

## Modo de proyecto

existing_project

## Contexto

En el Hito 5 de Akritas, la remediación automática no crea una Pull Request directamente tras generar código; ejecuta validaciones determinísticas (`test`, `build`, `static_analysis`) según el contrato OpenAPI (`ValidationResult`).

Si alguna validación falla (`ValidationStatus = failed`):
1. El estado de la remediación se consolida como `failed`.
2. Se debe mostrar al usuario qué validaciones se ejecutaron, cuáles pasaron, cuáles fallaron y la evidencia técnica (`output_excerpt`).
3. Debe quedar **completamente explícito que NO se creó Pull Request** (ADR-004: frontera de seguridad).

Si todas las validaciones pasan (`ValidationStatus = passed`):
1. El estado avanza a `validated` o `pull_request_created`.
2. Se muestra el reporte de checks superados con sus evidencias.

## Objetivo

Implementar el visualizador detallado de resultados de validación (`ValidationResultsViewer`):
- Servicio API cliente para obtener `GET /remediations/{remediation_id}/validation-results`.
- Componente expandible/acordeón con lista de validaciones ejecutadas (`test`, `build`, `static_analysis`).
- Indicador de estado por check (`passed`, `failed`, `running`, `pending`).
- Visor de evidencia / output sanitizado (`output_excerpt`) con estilo monospaciado.
- Banner explícito de fallo de remediación cuando alguna validación falla, confirmando que la PR fue bloqueada.

## Criterios de aceptación

1. Permite listar y visualizar los `ValidationResult` de una remediación activa o completada.
2. Cada resultado muestra:
   - Tipo de validación (`test`, `build`, `static_analysis`) y nombre del comando/check (`name`).
   - Estado visual con badge e ícono (`passed` en verde, `failed` en rojo, `running` con spinner, `pending` neutral).
   - Duración / timestamps (`started_at`, `finished_at`).
   - Resumen del resultado (`summary`).
   - Evidencia técnica / traza (`output_excerpt`) formateada.
3. Si existe al menos un check fallido:
   - Se muestra banner de alerta destacando el fallo.
   - Se incluye el texto mandatorio: **"Remediation Failed — No Pull Request was created"**.
4. Pruebas unitarias de componentes y validación limpia de TypeScript y ESLint.

## Restricciones técnicas

- Stack: Next.js App Router, React 19, TypeScript, CSS Modules.
- Uso estricto de tipos de `@/core/libs/api-client`.

