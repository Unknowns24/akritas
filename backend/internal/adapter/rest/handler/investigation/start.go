package investigation

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	operationdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/operation"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	"github.com/google/uuid"
)

func (h *Handler) Start(w http.ResponseWriter, r *http.Request) {
	incidentID, err := uuid.Parse(r.PathValue("incident_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	idempotencyKey, err := request.IdempotencyKey(r)
	if err != nil {
		response.Invalid(w, r)
		return
	}
	operation, err := h.investigations.StartIncidentInvestigation(r.Context(), portsin.StartIncidentInvestigationCommand{
		IncidentID: incidentID, IdempotencyKey: idempotencyKey,
	})
	if err != nil {
		response.Error(w, r, err)
		return
	}
	w.Header().Set("Retry-After", "5")
	response.JSON(w, http.StatusAccepted, commondto.DataResponseDTO[operationdto.OperationDTO]{Data: mapper.OperationToDTO(*operation)})
}
