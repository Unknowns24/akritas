package git

import "github.com/Unknowns24/akritas/backend/internal/core/domain"

var ErrGitBinaryUnavailable = &domain.Error{
	Code: "0x303001I", Message: "git binary unavailable", UserMessage: "No se pudo ejecutar Git para preparar el workspace.",
}

var ErrInvalidWorkspace = &domain.Error{
	Code: "0x303002V", Message: "invalid git workspace", UserMessage: "El workspace de Git no es válido.",
}

var ErrBaseBranchNotFound = &domain.Error{
	Code: "0x303003N", Message: "base branch not found", UserMessage: "La rama base no existe en el repositorio.",
}

var ErrBranchAlreadyExists = &domain.Error{
	Code: "0x303004C", Message: "branch already exists", UserMessage: "La rama de remediación ya existe.",
}

var ErrProtectedBranchTarget = &domain.Error{
	Code: "0x303005F", Message: "protected branch target", UserMessage: "No se puede operar directamente sobre la rama base/protegida.",
}

var ErrGitCommandFailed = &domain.Error{
	Code: "0x303006I", Message: "git command failed", UserMessage: "El comando de Git falló inesperadamente.",
}

func Catalog() map[string]*domain.Error {
	return map[string]*domain.Error{
		"ErrGitBinaryUnavailable":  ErrGitBinaryUnavailable,
		"ErrInvalidWorkspace":      ErrInvalidWorkspace,
		"ErrBaseBranchNotFound":    ErrBaseBranchNotFound,
		"ErrBranchAlreadyExists":   ErrBranchAlreadyExists,
		"ErrProtectedBranchTarget": ErrProtectedBranchTarget,
		"ErrGitCommandFailed":      ErrGitCommandFailed,
	}
}
