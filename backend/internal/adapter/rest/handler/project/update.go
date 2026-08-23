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

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("project_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	var body projectdto.UpdateProjectRequestDTO
	if request.DecodeJSON(w, r, &body) != nil || !validUpdate(body) {
		response.Invalid(w, r)
		return
	}
	updated, err := h.projects.Update(r.Context(), mapper.UpdateProjectToCommand(id, body))
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[projectdto.ProjectDTO]{Data: mapper.ProjectResultToDTO(updated)})
}

func validUpdate(value projectdto.UpdateProjectRequestDTO) bool {
	if value.Name == nil && value.Description == nil && value.GitHubAccountID == nil && value.RepositoryIdentifier == nil && value.DefaultBranch == nil && value.DokploySource == nil {
		return false
	}
	if value.Name != nil && (len(strings.TrimSpace(*value.Name)) < 1 || len(*value.Name) > 150) {
		return false
	}
	if value.Description != nil && len(*value.Description) > 2000 {
		return false
	}
	if value.GitHubAccountID != nil && *value.GitHubAccountID == uuid.Nil {
		return false
	}
	if value.DokploySource != nil && !validSource(*value.DokploySource) {
		return false
	}
	for _, field := range []*string{value.RepositoryIdentifier, value.DefaultBranch} {
		if field != nil && (len(strings.TrimSpace(*field)) < 1 || len(*field) > 255) {
			return false
		}
	}
	return true
}
