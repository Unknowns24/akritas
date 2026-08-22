package github

import (
	"net/http"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
	"github.com/google/uuid"
)

func (h *Handler) UpdateAccount(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(r.PathValue("account_id"))
	if err != nil {
		response.Invalid(w, r)
		return
	}
	var body dto.UpdateGitHubAccountRequestDTO
	if request.DecodeJSON(w, r, &body) != nil || (body.DisplayName == nil && body.PersonalAccessToken == nil) || (body.DisplayName != nil && (len(strings.TrimSpace(*body.DisplayName)) < 1 || len(*body.DisplayName) > 100)) || (body.PersonalAccessToken != nil && (len(*body.PersonalAccessToken) < 20 || len(*body.PersonalAccessToken) > 512)) {
		response.Invalid(w, r)
		return
	}
	account, err := h.accounts.Update(r.Context(), id, mapper.UpdateGitHubAccountToCommand(body))
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusOK, dto.DataResponseDTO[dto.GitHubAccountDTO]{Data: mapper.GitHubAccountToDTO(*account)})
}
