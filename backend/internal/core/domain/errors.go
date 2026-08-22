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
	ErrInvalidAdministrator           = newDomainError("0x401001V", "invalid administrator", "Los datos del administrador no son válidos.")
	ErrInvalidAdministratorSession    = newDomainError("0x401002V", "invalid administrator session", "La sesión no es válida.")
	ErrInactiveAdministratorSession   = newDomainError("0x401003U", "inactive administrator session", "La sesión no está activa.")
	ErrAdministratorSessionTransition = newDomainError("0x401004C", "invalid administrator session transition", "No se pudo actualizar la sesión.")
	ErrInvalidBootstrapToken          = newDomainError("0x401005V", "invalid bootstrap token", "El token de instalación no es válido.")
	ErrInvalidPendingEnrollment       = newDomainError("0x401006V", "invalid pending enrollment", "Los datos de inscripción no son válidos.")
	ErrAdministratorAlreadyExists     = newDomainError("0x401007C", "administrator already exists", "El registro inicial ya no está disponible.")
	ErrInvalidIntegrationStatus       = newDomainError("0x402001V", "invalid integration status", "El estado de la integración no es válido.")
	ErrInvalidGitHubAccount           = newDomainError("0x402002V", "invalid GitHub account", "La cuenta de GitHub no es válida.")
	ErrInvalidGitHubRepository        = newDomainError("0x402003V", "invalid GitHub repository", "El repositorio de GitHub no es válido.")
	ErrInvalidDokployServer           = newDomainError("0x402004V", "invalid Dokploy server", "El servidor Dokploy no es válido.")
	ErrInvalidDokployApplication      = newDomainError("0x402005V", "invalid Dokploy application", "La aplicación Dokploy no es válida.")
	ErrInvalidConnectionTestStatus    = newDomainError("0x402006V", "invalid connection test status", "El resultado de conexión no es válido.")
	ErrInvalidMonitoringStatus        = newDomainError("0x403001V", "invalid monitoring status", "El estado de monitoreo no es válido.")
	ErrInvalidProjectHealthStatus     = newDomainError("0x403002V", "invalid project health status", "El estado de salud del proyecto no es válido.")
	ErrInvalidProject                 = newDomainError("0x403003V", "invalid project", "El proyecto no es válido.")
	ErrInvalidMonitoringConfiguration = newDomainError("0x403004V", "invalid monitoring configuration", "La configuración de monitoreo no es válida.")
	ErrInvalidAutomationPolicy        = newDomainError("0x403005V", "invalid automation policy", "La configuración de automatización no es válida.")
	ErrInvalidDetectionRule           = newDomainError("0x403006V", "invalid detection rule", "La regla de detección no es válida.")
	ErrInvalidSeverity                = newDomainError("0x404001V", "invalid severity", "La severidad no es válida.")
	ErrInvalidIncidentPhase           = newDomainError("0x404002V", "invalid incident phase", "La fase del incidente no es válida.")
	ErrInvalidTerminalOutcome         = newDomainError("0x404003V", "invalid terminal outcome", "El resultado del incidente no es válido.")
	ErrInvalidSanitizedLogRecord      = newDomainError("0x404004V", "invalid sanitized log record", "El registro de log no es válido.")
	ErrInvalidLogEvent                = newDomainError("0x404005V", "invalid log event", "El evento de log no es válido.")
	ErrInvalidIncident                = newDomainError("0x404006V", "invalid incident", "El incidente no es válido.")
	ErrIncidentTransition             = newDomainError("0x404007C", "invalid incident transition", "El incidente no puede cambiar a ese estado.")
	ErrIncidentNotGroupable           = newDomainError("0x404008C", "incident occurrence is not groupable", "La ocurrencia no puede agruparse en este incidente.")
	ErrInvalidGitHubIssueReference    = newDomainError("0x404009V", "invalid GitHub issue reference", "La referencia a la Issue no es válida.")
	ErrInvalidInvestigationStatus     = newDomainError("0x405001V", "invalid investigation status", "El estado de la investigación no es válido.")
	ErrInvalidRootCauseStatus         = newDomainError("0x405002V", "invalid root cause status", "El estado de causa raíz no es válido.")
	ErrInvalidResolutionStatus        = newDomainError("0x405003V", "invalid resolution status", "El estado de resolución no es válido.")
	ErrInvalidEvidenceType            = newDomainError("0x405004V", "invalid evidence type", "El tipo de evidencia no es válido.")
	ErrInvalidInvestigation           = newDomainError("0x405005V", "invalid investigation", "La investigación no es válida.")
	ErrInvestigationTransition        = newDomainError("0x405006C", "invalid investigation transition", "La investigación no puede cambiar a ese estado.")
	ErrInvalidEvidence                = newDomainError("0x405007V", "invalid evidence", "La evidencia no es válida.")
	ErrInvalidRemediationStatus       = newDomainError("0x406001V", "invalid remediation status", "El estado de remediación no es válido.")
	ErrInvalidValidationStatus        = newDomainError("0x406002V", "invalid validation status", "El estado de validación no es válido.")
	ErrInvalidValidationType          = newDomainError("0x406003V", "invalid validation type", "El tipo de validación no es válido.")
	ErrInvalidCodeChangeType          = newDomainError("0x406004V", "invalid code change type", "El tipo de cambio no es válido.")
	ErrInvalidRemediation             = newDomainError("0x406005V", "invalid remediation", "La remediación no es válida.")
	ErrRemediationTransition          = newDomainError("0x406006C", "invalid remediation transition", "La remediación no puede cambiar a ese estado.")
	ErrInvalidValidationResult        = newDomainError("0x406007V", "invalid validation result", "El resultado de validación no es válido.")
	ErrValidationTransition           = newDomainError("0x406008C", "invalid validation transition", "La validación no puede cambiar a ese estado.")
	ErrInvalidCodeChange              = newDomainError("0x406009V", "invalid code change", "El cambio de código no es válido.")
	ErrInvalidPullRequestReference    = newDomainError("0x40600AV", "invalid pull request reference", "La referencia a la Pull Request no es válida.")
)

// DomainErrors returns the complete stable catalog keyed by sentinel name.
func DomainErrors() map[string]*Error {
	return map[string]*Error{
		"ErrInvalidAdministrator":           ErrInvalidAdministrator,
		"ErrInvalidAdministratorSession":    ErrInvalidAdministratorSession,
		"ErrInactiveAdministratorSession":   ErrInactiveAdministratorSession,
		"ErrAdministratorSessionTransition": ErrAdministratorSessionTransition,
		"ErrInvalidBootstrapToken":          ErrInvalidBootstrapToken,
		"ErrInvalidPendingEnrollment":       ErrInvalidPendingEnrollment,
		"ErrAdministratorAlreadyExists":     ErrAdministratorAlreadyExists,
		"ErrInvalidIntegrationStatus":       ErrInvalidIntegrationStatus,
		"ErrInvalidGitHubAccount":           ErrInvalidGitHubAccount,
		"ErrInvalidGitHubRepository":        ErrInvalidGitHubRepository,
		"ErrInvalidDokployServer":           ErrInvalidDokployServer,
		"ErrInvalidDokployApplication":      ErrInvalidDokployApplication,
		"ErrInvalidConnectionTestStatus":    ErrInvalidConnectionTestStatus,
		"ErrInvalidMonitoringStatus":        ErrInvalidMonitoringStatus,
		"ErrInvalidProjectHealthStatus":     ErrInvalidProjectHealthStatus,
		"ErrInvalidProject":                 ErrInvalidProject,
		"ErrInvalidMonitoringConfiguration": ErrInvalidMonitoringConfiguration,
		"ErrInvalidAutomationPolicy":        ErrInvalidAutomationPolicy,
		"ErrInvalidDetectionRule":           ErrInvalidDetectionRule,
		"ErrInvalidSeverity":                ErrInvalidSeverity,
		"ErrInvalidIncidentPhase":           ErrInvalidIncidentPhase,
		"ErrInvalidTerminalOutcome":         ErrInvalidTerminalOutcome,
		"ErrInvalidSanitizedLogRecord":      ErrInvalidSanitizedLogRecord,
		"ErrInvalidLogEvent":                ErrInvalidLogEvent,
		"ErrInvalidIncident":                ErrInvalidIncident,
		"ErrIncidentTransition":             ErrIncidentTransition,
		"ErrIncidentNotGroupable":           ErrIncidentNotGroupable,
		"ErrInvalidGitHubIssueReference":    ErrInvalidGitHubIssueReference,
		"ErrInvalidInvestigationStatus":     ErrInvalidInvestigationStatus,
		"ErrInvalidRootCauseStatus":         ErrInvalidRootCauseStatus,
		"ErrInvalidResolutionStatus":        ErrInvalidResolutionStatus,
		"ErrInvalidEvidenceType":            ErrInvalidEvidenceType,
		"ErrInvalidInvestigation":           ErrInvalidInvestigation,
		"ErrInvestigationTransition":        ErrInvestigationTransition,
		"ErrInvalidEvidence":                ErrInvalidEvidence,
		"ErrInvalidRemediationStatus":       ErrInvalidRemediationStatus,
		"ErrInvalidValidationStatus":        ErrInvalidValidationStatus,
		"ErrInvalidValidationType":          ErrInvalidValidationType,
		"ErrInvalidCodeChangeType":          ErrInvalidCodeChangeType,
		"ErrInvalidRemediation":             ErrInvalidRemediation,
		"ErrRemediationTransition":          ErrRemediationTransition,
		"ErrInvalidValidationResult":        ErrInvalidValidationResult,
		"ErrValidationTransition":           ErrValidationTransition,
		"ErrInvalidCodeChange":              ErrInvalidCodeChange,
		"ErrInvalidPullRequestReference":    ErrInvalidPullRequestReference,
	}
}
