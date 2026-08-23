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
| `ErrInvalidCredentials` | `0x401009U` | Auth | Unauthorized | Las credenciales no son válidas. |
| `ErrInvalidIntegrationStatus` | `0x402001V` | Integrations | Validation | El estado de la integración no es válido. |
| `ErrInvalidGitHubAccount` | `0x402002V` | Integrations | Validation | La cuenta de GitHub no es válida. |
| `ErrInvalidGitHubRepository` | `0x402003V` | Integrations | Validation | El repositorio de GitHub no es válido. |
| `ErrInvalidDokployServer` | `0x402004V` | Integrations | Validation | El servidor Dokploy no es válido. |
| `ErrInvalidDokployApplication` | `0x402005V` | Integrations | Validation | La aplicación Dokploy no es válida. |
| `ErrInvalidConnectionTestStatus` | `0x402006V` | Integrations | Validation | El resultado de conexión no es válido. |
| `ErrInvalidDokploySource` | `0x402007V` | Integrations | Validation | La fuente Dokploy no es válida. |
| `ErrInvalidDokployCompose` | `0x402008V` | Integrations | Validation | El Compose Dokploy no es válido. |
| `ErrInvalidDokployComposeService` | `0x402009V` | Integrations | Validation | El servicio Compose Dokploy no es válido. |
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
| `ErrInvalidOperationType` | `0x407001V` | Operation | Validation | El tipo de operación no es válido. |
| `ErrInvalidOperationStatus` | `0x407002V` | Operation | Validation | El estado de la operación no es válido. |
| `ErrInvalidOperationResourceType` | `0x407003V` | Operation | Validation | El tipo de recurso de la operación no es válido. |
| `ErrInvalidOperation` | `0x407004V` | Operation | Validation | La operación no es válida. |
| `ErrOperationTransition` | `0x407005C` | Operation | Conflict | La operación no puede cambiar a ese estado. |
| `ErrIntegrationNotFound` | `0x502001N` | Integrations usecase | Not found | La integración solicitada no existe. |
| `ErrIntegrationConflict` | `0x502002C` | Integrations usecase | Conflict | La integración entra en conflicto con una configuración existente. |
| `ErrIntegrationInUse` | `0x502003C` | Integrations usecase | Conflict | La integración está asociada a un Project y no puede eliminarse. |
| `ErrGitHubCredentialRejected` | `0x502004V` | Integrations usecase | Validation | GitHub rechazó la credencial o la cuenta configurada. |
| `ErrDokployCredentialRejected` | `0x502005V` | Integrations usecase | Validation | Dokploy rechazó la credencial configurada. |
| `ErrManifestStateInvalid` | `0x502006V` | Integrations usecase | Validation | El intento de conexión con GitHub no es válido. |
| `ErrManifestStateConflict` | `0x502007C` | Integrations usecase | Conflict | El intento de conexión con GitHub ya fue utilizado o expiró. |
| `ErrIntegrationUnavailable` | `0x502008I` | Integrations usecase | Internal | No se pudo contactar la integración. |
| `ErrDokployContainerUnavailable` | `0x502009I` | Integrations usecase | Internal | El servicio no tiene un contenedor activo disponible. |
| `ErrProjectNotFound` | `0x503001N` | Project usecase | Not found | El proyecto no existe. |
| `ErrProjectRepositoryNotFound` | `0x503002N` | Project usecase | Not found | El repositorio seleccionado no existe o no pertenece a la integración. |
| `ErrProjectApplicationNotFound` | `0x503003N` | Project usecase | Not found | La aplicación seleccionada no existe. |
| `ErrProjectNameConflict` | `0x503004C` | Project usecase | Conflict | Ya existe un proyecto con ese nombre. |
| `ErrProjectApplicationConflict` | `0x503005C` | Project usecase | Conflict | La aplicación Dokploy ya está asociada a otro proyecto. |
| `ErrProjectMustBeDisabled` | `0x503006C` | Project usecase | Conflict | El monitoreo debe desactivarse antes de la operación. |
| `ErrProjectConcurrentModification` | `0x503007C` | Project usecase | Conflict | El proyecto cambió durante la operación. |
| `ErrProjectDefaultBranchMismatch` | `0x503008V` | Project usecase | Validation | La rama predeterminada no coincide con GitHub. |
| `ErrProjectHasDependencies` | `0x503009C` | Project usecase | Conflict | El proyecto tiene registros asociados y no puede eliminarse. |
| `ErrInvalidInitialLogIngestion` | `0x50300AV` | Project usecase | Validation | La opción de ingesta inicial no es válida. |
| `ErrProjectDokploySourceNotFound` | `0x50300BN` | Project usecase | Not found | La fuente Dokploy seleccionada no existe. |
| `ErrProjectDokploySourceConflict` | `0x50300CC` | Project usecase | Conflict | La fuente Dokploy ya está asociada a otro proyecto. |
| `ErrIncidentNotFound` | `0x504001N` | Incident usecase | Not found | El incidente no existe. |
| `ErrMonitoringContinuityLost` | `0x505001I` | Monitoring service | Internal | No se pudo verificar la continuidad de los logs. |
| `ErrMonitoringConcurrentModification` | `0x505002C` | Monitoring service | Conflict | El estado de monitoreo cambió durante el procesamiento. |
| `ErrAuthenticationRateLimited` | `0x501001R` | Auth usecase | Rate limited | Alcanzaste el límite de intentos. Probá nuevamente más tarde. |
| `ErrInvestigationNotFound` | `0x504002N` | Investigation usecase | Not found | La investigación no existe. |
| `ErrInvestigationAlreadyActive` | `0x504003C` | Investigation usecase | Conflict | Ya hay una investigación en curso para este incidente. |
| `ErrGitHubIssueAlreadyPublished` | `0x504004C` | Investigation usecase | Conflict | La investigación ya tiene una Issue publicada. |
| `ErrOperationNotFound` | `0x505001N` | Operation usecase | Not found | La operación no existe. |
| `ErrIntegrationPersistence` | `0x202001I` | Integrations database | Internal | No se pudo guardar la integración. |
| `ErrProjectPersistence` | `0x203001I` | Project database | Internal | No se pudo guardar el proyecto. |
| `ErrInvestigationPersistence` | `0x204001I` | Investigation database | Internal | No se pudo guardar la investigación. |
| `ErrOperationPersistence` | `0x205001I` | Operation database | Internal | No se pudo guardar la operación. |
| `ErrEvidencePersistence` | `0x206001I` | Evidence database | Internal | No se pudo guardar la evidencia. |
| `ErrIncidentPersistence` | `0x207001I` | Incident database | Internal | No se pudo consultar el incidente. |
| `ErrMonitoringPersistence` | `0x208001I` | Monitoring database | Internal | No se pudo guardar el estado de monitoreo. |
| `ErrGitHubIssueReferencePersistence` | `0x209001I` | GitHub Issue database | Internal | No se pudo guardar la referencia a la Issue. |
| `ErrInvalidRequest` | `0x102001V` | REST request | Validation | La solicitud contiene datos inválidos. |
| `ErrRequestFailed` | `0x102002I` | REST request | Internal | No se pudo completar la solicitud. |
| `ErrRateLimited` | `0x102003R` | REST request | Rate limited | Alcanzaste el límite de solicitudes. Probá nuevamente más tarde. |
| `ErrOriginForbidden` | `0x102004F` | REST request | Forbidden | El origen de la solicitud no está permitido. |
| `ErrInvalidGitHubAppPrivateKey` | `0x302001I` | GitHub adapter | Internal | La clave privada de la GitHub App no pudo utilizarse. |

Los adapters deben mapear el tipo final del código a HTTP sin exponer la causa envuelta: `V` a 400, `U` a 401, `F` a 403, `N` a 404, `C` a 409, `R` a 429 e `I` a 500.

Cada sentinel se declara en la capa que representa: errores REST bajo
`internal/adapter/rest/errors`, errores PostgreSQL bajo
`internal/adapter/db/postgres/errors` y errores internos de proveedor dentro del
adapter correspondiente. `internal/core/domain` conserva sólo errores de dominio
y casos de uso mediante el contrato común `domain.Error`.
