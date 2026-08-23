# AKR-QVAC-INFERENCE

Reemplazar `qvac.StubRunner` por un `InvestigationRunner` real que habla con
QVAC local vía HTTP OpenAI-compatible (`qvac serve openai`), valida el
resultado estructurado y permite completar Investigations de punta a punta.

Cubre PB-028, PB-032, PB-033, PB-034 y PB-035 (persistencia verificada con
httptest, no con el stub).

## Alcance

- Cliente HTTP encapsulado en `internal/adapter/external/qvac/`.
- Runner que produce `out.InvestigationRunResult` validado.
- Mapeo estricto de `root_cause_status` / `resolution_status` a enums de dominio.
- Errores de runtime/modelo/output inválido → error (Investigation failed).
- Wiring de producción reemplazando `NewStubRunner()`.
- Tests de ciclo runner→Complete→campos listos para GET.

## Fuera de alcance

- Tool calling (AKR-QVAC-TOOL-LOOP).
- Tools GitHub de contenido (AKR-GITHUB-REPO-TOOLS).
- CRUD REST de configuración QVAC.
- H4/H5 (Issues, Remediation, PRs).
