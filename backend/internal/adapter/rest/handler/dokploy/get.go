package dokploy

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	dokploydto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/dokploy"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/google/uuid"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("server_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	server, err := h.servers.Get(r.Context(), id)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[dokploydto.DokployServerDTO]{Data: mapper.DokployServerToDTO(*server)})
}
