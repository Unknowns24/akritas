# Security Review — FE-H5-01 Remediation Lifecycle UI

## Evaluación de Seguridad y Autonomía

- **Aislamiento de Secretos**:
  - Ningún secreto, token o credencial se renderiza ni se almacena en el cliente.
  - Los diffs provienen de `CodeChange.patch`, el cual llega saneado y redactado desde el backend (`redacted: true`).
- **Frontera de Autonomía (ADR-004)**:
  - Para `resolution_status = requires_human`, el sistema no expone acciones mutativas ni botones de PR.
  - Para `resolution_status = fixable`, la creación de PR es el límite absoluto de autonomía.
  - Mensaje explícito visible: *"Akritas never merges changes automatically"*.
- **Veredicto**: APROBADO (PASS).
