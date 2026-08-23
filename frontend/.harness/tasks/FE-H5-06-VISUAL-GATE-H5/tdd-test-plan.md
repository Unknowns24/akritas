# TDD Test Plan — FE-H5-06 Gate Visual H5

## Scope

Verificar que el componente `RemediationReviewPacket` exponga de forma completa y precisa el flujo Golden Flow `Issue → fixable → branch → regression test → fix → tests pass → PR` con todo el contexto requerido para la revisión humana.

## Tests to add/update

Ubicación: `akritas/frontend/src/features/incidents/utils/review-packet.utils.test.ts`.

1. **Test 1: Ensamblado del Paquete de Revisión Humana**
   - Debe extraer la Issue documentada con su link.
   - Debe certificar la clasificación `fixable`.
   - Debe listar la rama objetivo.
   - Debe listar los tests de regresión y cambios de código.
   - Debe resumir las validaciones pasadas.
   - Debe enlazar la Pull Request generada con su commit SHA.

2. **Test 2: Integridad del Contexto de Revisión**
   - Validar que no falten datos clave (causa raíz, archivos tocados, resumen de cambios).

## Acceptance criteria covered

- Todos los criterios de aceptación de `FE-H5-06` cubiertos.

