package automation

import (
	"net/http"

	automationdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/automation"
	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

func (h *Handler) GetPolicy(w http.ResponseWriter, r *http.Request) {
	value, err := h.automation.GetPolicy(r.Context())
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[automationdto.PolicyDTO]{Data: mapper.AutomationPolicyToDTO(value)})
}

func (h *Handler) PutPolicy(w http.ResponseWriter, r *http.Request) {
	var body automationdto.PolicyDTO
	if err := request.DecodeJSON(w, r, &body); err != nil {
		response.Invalid(w, r)
		return
	}
	policy, err := mapper.AutomationPolicyToDomain(body)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	value, err := h.automation.PutPolicy(r.Context(), policy)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[automationdto.PolicyDTO]{Data: mapper.AutomationPolicyToDTO(value)})
}
