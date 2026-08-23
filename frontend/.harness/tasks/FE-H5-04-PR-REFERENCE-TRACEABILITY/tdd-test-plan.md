# TDD Test Plan — FE-H5-04 Pull Request Reference y Traceability

## Scope

Verificar el cálculo y renderizado de la cadena de trazabilidad de 7 etapas (`Incident → Investigation → Issue → Remediation → Branch → Commit → Pull Request`).

## Tests to add/update

Ubicación: `akritas/frontend/src/features/incidents/utils/traceability.utils.test.ts`.

1. **Test 1: Cadena completa exitosa con Pull Request**
   - Debe contener 7 pasos en estado `completed`.
   - El paso de Issue debe contener el enlace a la GitHub Issue.
   - El paso de PR debe contener el enlace a la GitHub Pull Request y el número `#<num>`.

2. **Test 2: Cadena detenida por `resolution_status = requires_human`**
   - Los pasos `Incident`, `Investigation` e `Issue` deben estar en `completed`.
   - El paso `Remediation` debe marcarse como `halted` (Manual Intervention Required).
   - Los pasos `Branch`, `Commit` y `PR` deben marcarse como `not_applicable` o `pending`.

3. **Test 3: Cadena detenida por fallo en validación (`remediation.status = failed`)**
   - Los pasos hasta `Remediation` deben reflejar la ejecución.
   - `Remediation` debe marcarse como `failed`.
   - Los pasos `Commit` y `PR` deben marcarse como `halted` / `blocked` (No PR created).

## Expected failing tests before implementation

- Actualmente no existe el módulo `traceability.utils.ts`.

## Acceptance criteria covered

- Todos los criterios de aceptación de `FE-H5-04` cubiertos.

