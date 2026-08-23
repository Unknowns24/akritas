package domain

import "errors"

// Error is the stable error contract shared by domain and application boundaries.
// The wrapped cause is intentionally excluded from Error() to avoid leaking details.
type Error struct {
	Code        string
	Message     string
	UserMessage string
	cause       error
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.cause
}

func (e *Error) Is(target error) bool {
	var other *Error
	return e != nil && errors.As(target, &other) && e.Code == other.Code
}

func (e *Error) Wrap(cause error) error {
	if cause == nil {
		return e
	}
	wrapped := *e
	wrapped.cause = cause
	return &wrapped
}

func newDomainError(code, message, userMessage string) *Error {
	return &Error{Code: code, Message: message, UserMessage: userMessage}
}

var (
	ErrInvalidAdministrator              = newDomainError("0x401001V", "invalid administrator", "Los datos del administrador no son válidos.")
	ErrInvalidAdministratorSession       = newDomainError("0x401002V", "invalid administrator session", "La sesión no es válida.")
	ErrInactiveAdministratorSession      = newDomainError("0x401003U", "inactive administrator session", "La sesión no está activa.")
	ErrAdministratorSessionTransition    = newDomainError("0x401004C", "invalid administrator session transition", "No se pudo actualizar la sesión.")
	ErrInvalidBootstrapToken             = newDomainError("0x401005V", "invalid bootstrap token", "El token de instalación no es válido.")
	ErrInvalidPendingEnrollment          = newDomainError("0x401006V", "invalid pending enrollment", "Los datos de inscripción no son válidos.")
	ErrAdministratorAlreadyExists        = newDomainError("0x401007C", "administrator already exists", "El registro inicial ya no está disponible.")
	ErrInvalidTotpEnrollmentVerification = newDomainError("0x401008V", "invalid totp enrollment verification", "El código o la inscripción no son válidos.")
	ErrInvalidCredentials                = newDomainError("0x401009U", "invalid credentials", "Las credenciales no son válidas.")
	ErrAuthenticationRateLimited         = newDomainError("0x501001R", "authentication rate limited", "Alcanzaste el límite de intentos. Probá nuevamente más tarde.")
	ErrInvalidIntegrationStatus          = newDomainError("0x402001V", "invalid integration status", "El estado de la integración no es válido.")
	ErrInvalidGitHubAccount              = newDomainError("0x402002V", "invalid GitHub account", "La cuenta de GitHub no es válida.")
	ErrInvalidGitHubRepository           = newDomainError("0x402003V", "invalid GitHub repository", "El repositorio de GitHub no es válido.")
	ErrInvalidDokployServer              = newDomainError("0x402004V", "invalid Dokploy server", "El servidor Dokploy no es válido.")
	ErrInvalidDokployApplication         = newDomainError("0x402005V", "invalid Dokploy application", "La aplicación Dokploy no es válida.")
	ErrInvalidConnectionTestStatus       = newDomainError("0x402006V", "invalid connection test status", "El resultado de conexión no es válido.")
	ErrInvalidDokploySource              = newDomainError("0x402007V", "invalid Dokploy source", "La fuente Dokploy no es válida.")
	ErrInvalidDokployCompose             = newDomainError("0x402008V", "invalid Dokploy compose", "El Compose Dokploy no es válido.")
	ErrInvalidDokployComposeService      = newDomainError("0x402009V", "invalid Dokploy compose service", "El servicio Compose Dokploy no es válido.")
	ErrInvalidMonitoringStatus           = newDomainError("0x403001V", "invalid monitoring status", "El estado de monitoreo no es válido.")
	ErrInvalidProjectHealthStatus        = newDomainError("0x403002V", "invalid project health status", "El estado de salud del proyecto no es válido.")
	ErrInvalidProject                    = newDomainError("0x403003V", "invalid project", "El proyecto no es válido.")
	ErrInvalidMonitoringConfiguration    = newDomainError("0x403004V", "invalid monitoring configuration", "La configuración de monitoreo no es válida.")
	ErrInvalidAutomationPolicy           = newDomainError("0x403005V", "invalid automation policy", "La configuración de automatización no es válida.")
	ErrInvalidDetectionRule              = newDomainError("0x403006V", "invalid detection rule", "La regla de detección no es válida.")
	ErrInvalidSeverity                   = newDomainError("0x404001V", "invalid severity", "La severidad no es válida.")
	ErrInvalidIncidentPhase              = newDomainError("0x404002V", "invalid incident phase", "La fase del incidente no es válida.")
	ErrInvalidTerminalOutcome            = newDomainError("0x404003V", "invalid terminal outcome", "El resultado del incidente no es válido.")
	ErrInvalidSanitizedLogRecord         = newDomainError("0x404004V", "invalid sanitized log record", "El registro de log no es válido.")
	ErrInvalidLogEvent                   = newDomainError("0x404005V", "invalid log event", "El evento de log no es válido.")
	ErrInvalidIncident                   = newDomainError("0x404006V", "invalid incident", "El incidente no es válido.")
	ErrIncidentTransition                = newDomainError("0x404007C", "invalid incident transition", "El incidente no puede cambiar a ese estado.")
	ErrIncidentNotGroupable              = newDomainError("0x404008C", "incident occurrence is not groupable", "La ocurrencia no puede agruparse en este incidente.")
	ErrInvalidGitHubIssueReference       = newDomainError("0x404009V", "invalid GitHub issue reference", "La referencia a la Issue no es válida.")
	ErrInvalidInvestigationStatus        = newDomainError("0x405001V", "invalid investigation status", "El estado de la investigación no es válido.")
	ErrInvalidRootCauseStatus            = newDomainError("0x405002V", "invalid root cause status", "El estado de causa raíz no es válido.")
	ErrInvalidResolutionStatus           = newDomainError("0x405003V", "invalid resolution status", "El estado de resolución no es válido.")
	ErrInvalidEvidenceType               = newDomainError("0x405004V", "invalid evidence type", "El tipo de evidencia no es válido.")
	ErrInvalidInvestigation              = newDomainError("0x405005V", "invalid investigation", "La investigación no es válida.")
	ErrInvestigationTransition           = newDomainError("0x405006C", "invalid investigation transition", "La investigación no puede cambiar a ese estado.")
	ErrInvalidEvidence                   = newDomainError("0x405007V", "invalid evidence", "La evidencia no es válida.")
	ErrInvalidRemediationStatus          = newDomainError("0x406001V", "invalid remediation status", "El estado de remediación no es válido.")
	ErrInvalidValidationStatus           = newDomainError("0x406002V", "invalid validation status", "El estado de validación no es válido.")
	ErrInvalidValidationType             = newDomainError("0x406003V", "invalid validation type", "El tipo de validación no es válido.")
	ErrInvalidCodeChangeType             = newDomainError("0x406004V", "invalid code change type", "El tipo de cambio no es válido.")
	ErrInvalidRemediation                = newDomainError("0x406005V", "invalid remediation", "La remediación no es válida.")
	ErrRemediationTransition             = newDomainError("0x406006C", "invalid remediation transition", "La remediación no puede cambiar a ese estado.")
	ErrInvalidValidationResult           = newDomainError("0x406007V", "invalid validation result", "El resultado de validación no es válido.")
	ErrValidationTransition              = newDomainError("0x406008C", "invalid validation transition", "La validación no puede cambiar a ese estado.")
	ErrInvalidCodeChange                 = newDomainError("0x406009V", "invalid code change", "El cambio de código no es válido.")
	ErrInvalidPullRequestReference       = newDomainError("0x40600AV", "invalid pull request reference", "La referencia a la Pull Request no es válida.")
	ErrIntegrationNotFound               = newDomainError("0x502001N", "integration not found", "La integración solicitada no existe.")
	ErrIntegrationConflict               = newDomainError("0x502002C", "integration conflict", "La integración entra en conflicto con una configuración existente.")
	ErrIntegrationInUse                  = newDomainError("0x502003C", "integration in use", "La integración está asociada a un Project y no puede eliminarse.")
	ErrGitHubCredentialRejected          = newDomainError("0x502004V", "GitHub credential rejected", "GitHub rechazó la credencial o la cuenta configurada.")
	ErrDokployCredentialRejected         = newDomainError("0x502005V", "Dokploy credential rejected", "Dokploy rechazó la credencial configurada.")
	ErrManifestStateInvalid              = newDomainError("0x502006V", "invalid GitHub manifest state", "El intento de conexión con GitHub no es válido.")
	ErrManifestStateConflict             = newDomainError("0x502007C", "GitHub manifest state conflict", "El intento de conexión con GitHub ya fue utilizado o expiró.")
	ErrIntegrationUnavailable            = newDomainError("0x502008I", "integration unavailable", "No se pudo contactar la integración.")
	ErrDokployContainerUnavailable       = newDomainError("0x502009I", "Dokploy container unavailable", "El servicio no tiene un contenedor activo disponible.")
	ErrProjectNotFound                   = newDomainError("0x503001N", "project not found", "El proyecto no existe.")
	ErrProjectRepositoryNotFound         = newDomainError("0x503002N", "project repository not found", "El repositorio seleccionado no existe o no pertenece a la integración.")
	ErrProjectApplicationNotFound        = newDomainError("0x503003N", "project application not found", "La aplicación seleccionada no existe.")
	ErrProjectNameConflict               = newDomainError("0x503004C", "project name conflict", "Ya existe un proyecto con ese nombre.")
	ErrProjectApplicationConflict        = newDomainError("0x503005C", "project application conflict", "La aplicación Dokploy ya está asociada a otro proyecto.")
	ErrProjectMustBeDisabled             = newDomainError("0x503006C", "project must be disabled", "Desactivá el monitoreo antes de realizar esta operación.")
	ErrProjectConcurrentModification     = newDomainError("0x503007C", "project changed concurrently", "El proyecto cambió; volvé a intentar la operación.")
	ErrProjectDefaultBranchMismatch      = newDomainError("0x503008V", "project default branch mismatch", "La rama predeterminada no coincide con GitHub.")
	ErrProjectHasDependencies            = newDomainError("0x503009C", "project has dependencies", "El proyecto tiene registros asociados y no puede eliminarse.")

	ErrInvalidOperationType         = newDomainError("0x407001V", "invalid operation type", "El tipo de operación no es válido.")
	ErrInvalidOperationStatus       = newDomainError("0x407002V", "invalid operation status", "El estado de la operación no es válido.")
	ErrInvalidOperationResourceType = newDomainError("0x407003V", "invalid operation resource type", "El tipo de recurso de la operación no es válido.")
	ErrInvalidOperation             = newDomainError("0x407004V", "invalid operation", "La operación no es válida.")
	ErrOperationTransition          = newDomainError("0x407005C", "invalid operation transition", "La operación no puede cambiar a ese estado.")

	ErrInvalidInitialLogIngestion   = newDomainError("0x50300AV", "invalid initial log ingestion", "La opción de ingesta inicial no es válida.")
	ErrProjectDokploySourceNotFound = newDomainError("0x50300BN", "project Dokploy source not found", "La fuente Dokploy seleccionada no existe.")
	ErrProjectDokploySourceConflict = newDomainError("0x50300CC", "project Dokploy source conflict", "La fuente Dokploy ya está asociada a otro proyecto.")

	ErrIncidentNotFound            = newDomainError("0x504001N", "incident not found", "El incidente no existe.")
	ErrInvestigationNotFound       = newDomainError("0x504002N", "investigation not found", "La investigación no existe.")
	ErrInvestigationAlreadyActive  = newDomainError("0x504003C", "investigation already active", "Ya hay una investigación en curso para este incidente.")
	ErrGitHubIssueAlreadyPublished = newDomainError("0x504004C", "GitHub issue already published", "La investigación ya tiene una Issue publicada.")

	ErrOperationNotFound                = newDomainError("0x505001N", "operation not found", "La operación no existe.")
	ErrMonitoringContinuityLost         = newDomainError("0x505001I", "monitoring log continuity lost", "No se pudo verificar la continuidad de los logs.")
	ErrMonitoringConcurrentModification = newDomainError("0x505002C", "monitoring checkpoint changed concurrently", "El monitoreo cambió; se reintentará el procesamiento.")
)

// DomainErrors returns the complete stable catalog keyed by sentinel name.
func DomainErrors() map[string]*Error {
	return map[string]*Error{
		"ErrInvalidAdministrator":              ErrInvalidAdministrator,
		"ErrInvalidAdministratorSession":       ErrInvalidAdministratorSession,
		"ErrInactiveAdministratorSession":      ErrInactiveAdministratorSession,
		"ErrAdministratorSessionTransition":    ErrAdministratorSessionTransition,
		"ErrInvalidBootstrapToken":             ErrInvalidBootstrapToken,
		"ErrInvalidPendingEnrollment":          ErrInvalidPendingEnrollment,
		"ErrAdministratorAlreadyExists":        ErrAdministratorAlreadyExists,
		"ErrInvalidTotpEnrollmentVerification": ErrInvalidTotpEnrollmentVerification,
		"ErrInvalidCredentials":                ErrInvalidCredentials,
		"ErrInvalidIntegrationStatus":          ErrInvalidIntegrationStatus,
		"ErrInvalidGitHubAccount":              ErrInvalidGitHubAccount,
		"ErrInvalidGitHubRepository":           ErrInvalidGitHubRepository,
		"ErrInvalidDokployServer":              ErrInvalidDokployServer,
		"ErrInvalidDokployApplication":         ErrInvalidDokployApplication,
		"ErrInvalidDokploySource":              ErrInvalidDokploySource,
		"ErrInvalidDokployCompose":             ErrInvalidDokployCompose,
		"ErrInvalidDokployComposeService":      ErrInvalidDokployComposeService,
		"ErrInvalidConnectionTestStatus":       ErrInvalidConnectionTestStatus,
		"ErrInvalidMonitoringStatus":           ErrInvalidMonitoringStatus,
		"ErrInvalidProjectHealthStatus":        ErrInvalidProjectHealthStatus,
		"ErrInvalidProject":                    ErrInvalidProject,
		"ErrInvalidMonitoringConfiguration":    ErrInvalidMonitoringConfiguration,
		"ErrInvalidAutomationPolicy":           ErrInvalidAutomationPolicy,
		"ErrInvalidDetectionRule":              ErrInvalidDetectionRule,
		"ErrInvalidSeverity":                   ErrInvalidSeverity,
		"ErrInvalidIncidentPhase":              ErrInvalidIncidentPhase,
		"ErrInvalidTerminalOutcome":            ErrInvalidTerminalOutcome,
		"ErrInvalidSanitizedLogRecord":         ErrInvalidSanitizedLogRecord,
		"ErrInvalidLogEvent":                   ErrInvalidLogEvent,
		"ErrInvalidIncident":                   ErrInvalidIncident,
		"ErrIncidentTransition":                ErrIncidentTransition,
		"ErrIncidentNotGroupable":              ErrIncidentNotGroupable,
		"ErrInvalidGitHubIssueReference":       ErrInvalidGitHubIssueReference,
		"ErrInvalidInvestigationStatus":        ErrInvalidInvestigationStatus,
		"ErrInvalidRootCauseStatus":            ErrInvalidRootCauseStatus,
		"ErrInvalidResolutionStatus":           ErrInvalidResolutionStatus,
		"ErrInvalidEvidenceType":               ErrInvalidEvidenceType,
		"ErrInvalidInvestigation":              ErrInvalidInvestigation,
		"ErrInvestigationTransition":           ErrInvestigationTransition,
		"ErrInvalidEvidence":                   ErrInvalidEvidence,
		"ErrInvalidRemediationStatus":          ErrInvalidRemediationStatus,
		"ErrInvalidValidationStatus":           ErrInvalidValidationStatus,
		"ErrInvalidValidationType":             ErrInvalidValidationType,
		"ErrInvalidCodeChangeType":             ErrInvalidCodeChangeType,
		"ErrInvalidRemediation":                ErrInvalidRemediation,
		"ErrRemediationTransition":             ErrRemediationTransition,
		"ErrInvalidValidationResult":           ErrInvalidValidationResult,
		"ErrValidationTransition":              ErrValidationTransition,
		"ErrInvalidCodeChange":                 ErrInvalidCodeChange,
		"ErrInvalidPullRequestReference":       ErrInvalidPullRequestReference,
		"ErrInvalidOperationType":              ErrInvalidOperationType,
		"ErrInvalidOperationStatus":            ErrInvalidOperationStatus,
		"ErrInvalidOperationResourceType":      ErrInvalidOperationResourceType,
		"ErrInvalidOperation":                  ErrInvalidOperation,
		"ErrOperationTransition":               ErrOperationTransition,
	}
}

// IntegrationErrors returns stable errors introduced by the integration application boundary.
func IntegrationErrors() map[string]*Error {
	return map[string]*Error{
		"ErrIntegrationNotFound":         ErrIntegrationNotFound,
		"ErrIntegrationConflict":         ErrIntegrationConflict,
		"ErrIntegrationInUse":            ErrIntegrationInUse,
		"ErrGitHubCredentialRejected":    ErrGitHubCredentialRejected,
		"ErrDokployCredentialRejected":   ErrDokployCredentialRejected,
		"ErrManifestStateInvalid":        ErrManifestStateInvalid,
		"ErrManifestStateConflict":       ErrManifestStateConflict,
		"ErrIntegrationUnavailable":      ErrIntegrationUnavailable,
		"ErrDokployContainerUnavailable": ErrDokployContainerUnavailable,
	}
}

func AuthenticationErrors() map[string]*Error {
	return map[string]*Error{
		"ErrAuthenticationRateLimited": ErrAuthenticationRateLimited,
	}
}

func ProjectErrors() map[string]*Error {
	return map[string]*Error{
		"ErrProjectNotFound":               ErrProjectNotFound,
		"ErrProjectRepositoryNotFound":     ErrProjectRepositoryNotFound,
		"ErrProjectApplicationNotFound":    ErrProjectApplicationNotFound,
		"ErrProjectNameConflict":           ErrProjectNameConflict,
		"ErrProjectApplicationConflict":    ErrProjectApplicationConflict,
		"ErrProjectMustBeDisabled":         ErrProjectMustBeDisabled,
		"ErrProjectConcurrentModification": ErrProjectConcurrentModification,
		"ErrProjectDefaultBranchMismatch":  ErrProjectDefaultBranchMismatch,
		"ErrProjectHasDependencies":        ErrProjectHasDependencies,
		"ErrInvalidInitialLogIngestion":    ErrInvalidInitialLogIngestion,
		"ErrProjectDokploySourceNotFound":  ErrProjectDokploySourceNotFound,
		"ErrProjectDokploySourceConflict":  ErrProjectDokploySourceConflict,
	}
}

func IncidentErrors() map[string]*Error {
	return map[string]*Error{"ErrIncidentNotFound": ErrIncidentNotFound}
}

func MonitoringErrors() map[string]*Error {
	return map[string]*Error{
		"ErrMonitoringContinuityLost":         ErrMonitoringContinuityLost,
		"ErrMonitoringConcurrentModification": ErrMonitoringConcurrentModification,
	}
}

// InvestigationErrors returns stable errors introduced by the investigation application boundary.
func InvestigationErrors() map[string]*Error {
	return map[string]*Error{
		"ErrInvestigationNotFound":       ErrInvestigationNotFound,
		"ErrInvestigationAlreadyActive":  ErrInvestigationAlreadyActive,
		"ErrGitHubIssueAlreadyPublished": ErrGitHubIssueAlreadyPublished,
	}
}

// OperationErrors returns stable errors introduced by the operation application boundary.
func OperationErrors() map[string]*Error {
	return map[string]*Error{
		"ErrOperationNotFound": ErrOperationNotFound,
	}
}
