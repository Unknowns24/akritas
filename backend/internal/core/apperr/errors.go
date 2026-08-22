package apperr

import "github.com/Unknowns24/akritas/backend/internal/core/domain"

var (
	ErrUnauthenticated            = domain.NewError("0x501001U", "unauthenticated", "La sesión no es válida o está ausente.")
	ErrProjectNotFound            = domain.NewError("0x503001N", "project not found", "El proyecto no existe.")
	ErrGitHubAccountNotFound      = domain.NewError("0x503002N", "GitHub account not found", "La cuenta de GitHub no existe.")
	ErrDokployServerNotFound      = domain.NewError("0x503003N", "Dokploy server not found", "El servidor Dokploy no existe.")
	ErrRepositoryNotResolvable    = domain.NewError("0x503004N", "repository snapshot not resolvable", "No se pudo resolver el repositorio de GitHub.")
	ErrApplicationNotResolvable   = domain.NewError("0x503005N", "application snapshot not resolvable", "No se pudo resolver la aplicación Dokploy.")
	ErrProjectNameConflict        = domain.NewError("0x503006C", "project name already exists", "Ya existe un proyecto con ese nombre.")
	ErrProjectApplicationConflict = domain.NewError("0x503007C", "Dokploy application already assigned", "La aplicación Dokploy ya está asignada a otro proyecto.")
	ErrInvalidProjectCommand      = domain.NewError("0x503008V", "invalid project command", "La solicitud del proyecto no es válida.")
)

func Catalog() map[string]*domain.Error {
	return map[string]*domain.Error{
		"ErrUnauthenticated":            ErrUnauthenticated,
		"ErrProjectNotFound":            ErrProjectNotFound,
		"ErrGitHubAccountNotFound":      ErrGitHubAccountNotFound,
		"ErrDokployServerNotFound":      ErrDokployServerNotFound,
		"ErrRepositoryNotResolvable":    ErrRepositoryNotResolvable,
		"ErrApplicationNotResolvable":   ErrApplicationNotResolvable,
		"ErrProjectNameConflict":        ErrProjectNameConflict,
		"ErrProjectApplicationConflict": ErrProjectApplicationConflict,
		"ErrInvalidProjectCommand":      ErrInvalidProjectCommand,
	}
}
