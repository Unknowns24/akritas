package github

import (
	"net/http"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/google/uuid"
)

func (h *Handler) GetAccount(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("account_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	account, err := h.accounts.Get(r.Context(), id)
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.DataResponseDTO[dto.GitHubAccountDTO]{Data: mapper.GitHubAccountToDTO(*account)})
}
