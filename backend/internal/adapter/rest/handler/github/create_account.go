package github

import (
	"net/http"
	"strings"

	commondto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/common"
	githubdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/github"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/request"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/response"
)

func (h *Handler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	var body githubdto.CreateGitHubPATAccountRequestDTO
	if request.DecodeJSON(w, r, &body) != nil || len(strings.TrimSpace(body.DisplayName)) < 1 || len(body.DisplayName) > 100 || len(strings.TrimSpace(body.AccountIdentifier)) < 1 || len(body.AccountIdentifier) > 100 || len(body.PersonalAccessToken) < 20 || len(body.PersonalAccessToken) > 512 {
		response.Invalid(w, r)
		return
	}
	account, err := h.accounts.CreatePAT(r.Context(), mapper.CreateGitHubPATAccountToCommand(body))
	if err != nil {
		response.Error(w, r, err)
		return
	}
	response.JSON(w, http.StatusCreated, commondto.DataResponseDTO[githubdto.GitHubAccountDTO]{Data: mapper.GitHubAccountToDTO(*account)})
}
