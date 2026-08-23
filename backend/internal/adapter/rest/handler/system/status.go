package system

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	systemdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/system"
	operationdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/operation"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/google/uuid"
)

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	value, err := h.system.GetStatus(r.Context())
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[systemdto.SystemStatusDTO]{Data: mapper.SystemStatusToDTO(value)})
}

func (h *Handler) Diagnostics(w http.ResponseWriter, r *http.Request) {
	idempotencyKey := uuid.New()
	if r.Header.Get("Idempotency-Key") != "" {
		parsed, err := request.IdempotencyKey(r)
		if err != nil {
			response.Invalid(w, r)
			return
		}
		idempotencyKey = parsed
	}
	operation, err := h.system.RunDiagnostics(r.Context(), idempotencyKey)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	w.Header().Set("Retry-After", "5")
	response.JSON(w, http.StatusAccepted, commondto.DataResponseDTO[operationdto.OperationDTO]{Data: mapper.OperationToDTO(*operation)})
}
