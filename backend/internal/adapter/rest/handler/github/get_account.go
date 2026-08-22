package github

import (
	"net/http"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	githubdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/github"
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
	response.JSON(w, http.StatusOK, commondto.DataResponseDTO[githubdto.GitHubAccountDTO]{Data: mapper.GitHubAccountToDTO(*account)})
}
