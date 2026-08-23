package dashboard

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	dashboarddto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/dashboard"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

func (h *Handler) Overview(w http.ResponseWriter, r *http.Request) {
	value, err := h.dashboard.GetOverview(r.Context())
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[dashboarddto.OverviewDTO]{Data: mapper.OverviewToDTO(value)})
}
