package system

import (
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/envelope"
	systemdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/system"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/utils"
)

func Health(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, envelope.DataEnvelopeDTO[systemdto.HealthDTO]{Data: systemdto.HealthDTO{Status: "ok"}})
}

func Readiness(w http.ResponseWriter, _ *http.Request) {
	utils.WriteJSON(w, http.StatusOK, envelope.DataEnvelopeDTO[systemdto.ReadinessDTO]{Data: systemdto.ReadinessDTO{Status: "ready"}})
}
