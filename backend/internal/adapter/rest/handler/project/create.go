package project

import (
	"net/http"
	"strings"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/google/uuid"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body projectdto.CreateProjectRequestDTO
	if request.DecodeJSON(w, r, &body) != nil || !validCreate(body) {
		response.Invalid(w, r)
		return
	}
	command, err := mapper.CreateProjectToCommand(body)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	created, err := h.projects.Create(r.Context(), command)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusCreated, commondto.DataResponseDTO[projectdto.ProjectDTO]{Data: mapper.ProjectResultToDTO(created)})
}

func validCreate(value projectdto.CreateProjectRequestDTO) bool {
	return value.GitHubAccountID != uuid.Nil && value.DokployServerID != uuid.Nil && len(strings.TrimSpace(value.Name)) >= 1 && len(value.Name) <= 150 && len(value.Description) <= 2000 && len(strings.TrimSpace(value.RepositoryIdentifier)) >= 1 && len(value.RepositoryIdentifier) <= 255 && len(strings.TrimSpace(value.DefaultBranch)) >= 1 && len(value.DefaultBranch) <= 255 && len(strings.TrimSpace(value.ApplicationIdentifier)) >= 1 && len(value.ApplicationIdentifier) <= 255
}
