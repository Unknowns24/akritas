package github

import (
	"net/http"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

func (h *Handler) StartManifest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-store")
	var body dto.GitHubManifestRegistrationRequestDTO
	if request.DecodeJSON(w, r, &body) != nil {
		response.Invalid(w, r)
		return
	}
	if len(strings.TrimSpace(body.DisplayName)) < 1 || len(body.DisplayName) > 100 || len(body.Organization) > 100 {
		response.Invalid(w, r)
		return
	}
	result, err := h.apps.StartRegistration(r.Context(), mapper.GitHubManifestRegistrationToCommand(body))
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusCreated, dto.DataResponseDTO[dto.GitHubManifestRegistrationDTO]{Data: mapper.GitHubManifestRegistrationResultToDTO(result)})
}
