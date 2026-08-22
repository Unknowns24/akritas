package dokploy

import (
	"net/http"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/google/uuid"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("server_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	var body dto.UpdateDokployServerRequestDTO
	if request.DecodeJSON(w, r, &body) != nil || (body.Name == nil && body.BaseURL == nil && body.APICredential == nil) || (body.Name != nil && (len(strings.TrimSpace(*body.Name)) < 1 || len(*body.Name) > 100)) || (body.BaseURL != nil && (len(*body.BaseURL) < 1 || len(*body.BaseURL) > 2048)) || (body.APICredential != nil && (len(*body.APICredential) < 16 || len(*body.APICredential) > 1024)) {
		response.Invalid(w, r)
		return
	}
	server, err := h.servers.Update(r.Context(), id, mapper.UpdateDokployServerToCommand(body))
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.DataResponseDTO[dto.DokployServerDTO]{Data: mapper.DokployServerToDTO(*server)})
}
