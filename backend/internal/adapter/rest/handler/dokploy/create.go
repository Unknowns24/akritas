package dokploy

import (
	"net/http"
	"strings"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	dokploydto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/dokploy"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var body dokploydto.CreateDokployServerRequestDTO
	if request.DecodeJSON(w, r, &body) != nil || len(strings.TrimSpace(body.Name)) < 1 || len(body.Name) > 100 || len(body.BaseURL) < 1 || len(body.BaseURL) > 2048 || len(body.APICredential) < 16 || len(body.APICredential) > 1024 {
		response.Invalid(w, r)
		return
	}
	server, err := h.servers.Create(r.Context(), mapper.CreateDokployServerToCommand(body))
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusCreated, commondto.DataResponseDTO[dokploydto.DokployServerDTO]{Data: mapper.DokployServerToDTO(*server)})
}
