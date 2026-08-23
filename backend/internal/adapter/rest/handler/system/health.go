package system

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	systemdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/system"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[systemdto.HealthDTO]{Data: systemdto.HealthDTO{Status: "ok"}})
}

func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[systemdto.ReadinessDTO]{Data: systemdto.ReadinessDTO{Status: "ready"}})
}
