# Akritas domain error catalog

Los errores siguen `DxAAABBBT`. En esta fundación, `D=0` representa la instalación/plataforma y el primer nibble de `AAA=4` identifica la capa de dominio.

| Sentinel | Código | Componente | Tipo | Significado público |
| --- | --- | --- | --- | --- |
| `ErrInvalidAdministrator` | `0x401001V` | Auth | Validation | Los datos del administrador no son válidos. |
| `ErrInvalidAdministratorSession` | `0x401002V` | Auth | Validation | La sesión no es válida. |
| `ErrInactiveAdministratorSession` | `0x401003U` | Auth | Unauthorized | La sesión no está activa. |
| `ErrAdministratorSessionTransition` | `0x401004C` | Auth | Conflict | No se pudo actualizar la sesión. |
| `ErrInvalidBootstrapToken` | `0x401005V` | Auth | Validation | El token de instalación no es válido. |
| `ErrInvalidPendingEnrollment` | `0x401006V` | Auth | Validation | Los datos de inscripción no son válidos. |
| `ErrAdministratorAlreadyExists` | `0x401007C` | Auth | Conflict | El registro inicial ya no está disponible. |
| `ErrInvalidTotpEnrollmentVerification` | `0x401008V` | Auth | Validation | El código o la inscripción no son válidos. |
| `ErrInvalidIntegrationStatus` | `0x402001V` | Integrations | Validation | El estado de la integración no es válido. |
| `ErrInvalidGitHubAccount` | `0x402002V` | Integrations | Validation | La cuenta de GitHub no es válida. |
| `ErrInvalidGitHubRepository` | `0x402003V` | Integrations | Validation | El repositorio de GitHub no es válido. |
| `ErrInvalidDokployServer` | `0x402004V` | Integrations | Validation | El servidor Dokploy no es válido. |
| `ErrInvalidDokployApplication` | `0x402005V` | Integrations | Validation | La aplicación Dokploy no es válida. |
| `ErrInvalidConnectionTestStatus` | `0x402006V` | Integrations | Validation | El resultado de conexión no es válido. |
| `ErrInvalidMonitoringStatus` | `0x403001V` | Project | Validation | El estado de monitoreo no es válido. |
| `ErrInvalidProjectHealthStatus` | `0x403002V` | Project | Validation | El estado de salud del proyecto no es válido. |
| `ErrInvalidProject` | `0x403003V` | Project | Validation | El proyecto no es válido. |
| `ErrInvalidMonitoringConfiguration` | `0x403004V` | Project | Validation | La configuración de monitoreo no es válida. |
| `ErrInvalidAutomationPolicy` | `0x403005V` | Project | Validation | La configuración de automatización no es válida. |
| `ErrInvalidDetectionRule` | `0x403006V` | Project | Validation | La regla de detección no es válida. |
| `ErrInvalidSeverity` | `0x404001V` | Incident | Validation | La severidad no es válida. |
| `ErrInvalidIncidentPhase` | `0x404002V` | Incident | Validation | La fase del incidente no es válida. |
| `ErrInvalidTerminalOutcome` | `0x404003V` | Incident | Validation | El resultado del incidente no es válido. |
| `ErrInvalidSanitizedLogRecord` | `0x404004V` | Incident | Validation | El registro de log no es válido. |
| `ErrInvalidLogEvent` | `0x404005V` | Incident | Validation | El evento de log no es válido. |
| `ErrInvalidIncident` | `0x404006V` | Incident | Validation | El incidente no es válido. |
| `ErrIncidentTransition` | `0x404007C` | Incident | Conflict | El incidente no puede cambiar a ese estado. |
| `ErrIncidentNotGroupable` | `0x404008C` | Incident | Conflict | La ocurrencia no puede agruparse en este incidente. |
| `ErrInvalidGitHubIssueReference` | `0x404009V` | Incident | Validation | La referencia a la Issue no es válida. |
| `ErrInvalidInvestigationStatus` | `0x405001V` | Investigation | Validation | El estado de la investigación no es válido. |
| `ErrInvalidRootCauseStatus` | `0x405002V` | Investigation | Validation | El estado de causa raíz no es válido. |
| `ErrInvalidResolutionStatus` | `0x405003V` | Investigation | Validation | El estado de resolución no es válido. |
| `ErrInvalidEvidenceType` | `0x405004V` | Investigation | Validation | El tipo de evidencia no es válido. |
| `ErrInvalidInvestigation` | `0x405005V` | Investigation | Validation | La investigación no es válida. |
| `ErrInvestigationTransition` | `0x405006C` | Investigation | Conflict | La investigación no puede cambiar a ese estado. |
| `ErrInvalidEvidence` | `0x405007V` | Investigation | Validation | La evidencia no es válida. |
| `ErrInvalidRemediationStatus` | `0x406001V` | Remediation | Validation | El estado de remediación no es válido. |
| `ErrInvalidValidationStatus` | `0x406002V` | Remediation | Validation | El estado de validación no es válido. |
| `ErrInvalidValidationType` | `0x406003V` | Remediation | Validation | El tipo de validación no es válido. |
| `ErrInvalidCodeChangeType` | `0x406004V` | Remediation | Validation | El tipo de cambio no es válido. |
| `ErrInvalidRemediation` | `0x406005V` | Remediation | Validation | La remediación no es válida. |
| `ErrRemediationTransition` | `0x406006C` | Remediation | Conflict | La remediación no puede cambiar a ese estado. |
| `ErrInvalidValidationResult` | `0x406007V` | Remediation | Validation | El resultado de validación no es válido. |
| `ErrValidationTransition` | `0x406008C` | Remediation | Conflict | La validación no puede cambiar a ese estado. |
| `ErrInvalidCodeChange` | `0x406009V` | Remediation | Validation | El cambio de código no es válido. |
| `ErrInvalidPullRequestReference` | `0x40600AV` | Remediation | Validation | La referencia a la Pull Request no es válida. |

Los adapters deben mapear el tipo final del código a HTTP sin exponer la causa envuelta: `V` a 400, `U` a 401 y `C` a 409.

## Códigos de capa REST (`0x1...`)

Fallos de transporte que no representan una regla de negocio (body malformado,
error interno inesperado) usan el componente de capa `1` (REST) en vez de `4`
(dominio), y no forman parte del catálogo `internal/core/domain` ni de su test
de contrato — viven en `internal/adapter/rest/response`.

| Código | Componente | Tipo | Significado público |
| --- | --- | --- | --- |
| `0x100001I` | REST | Internal | Ocurrió un error inesperado. Intentá nuevamente más tarde. |
| `0x100002V` | REST | Validation | La solicitud contiene datos inválidos (falla de shape/formato, no de regla de negocio). |
| `0x100003C` | REST | Rate limit | Alcanzaste el límite de intentos. Probá nuevamente más tarde. Sin letra propia en el pattern de `ErrorCode` (no cubre 429); el status 429 lo fija el handler explícitamente, no este código. |
