package project

import (
	"net/http"
	"strconv"
	"strings"

	projectdto "github.com/Unknowns24/akritas/backend/internal/adapter/rest/dto/project"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/httperr"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/middleware"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/pagination"
	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/utils"
	"github.com/Unknowns24/akritas/backend/internal/core/apperr"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
)

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	query, err := h.listQuery(r)
	if err != nil {
		httperr.Write(w, apperr.ErrInvalidProjectCommand.Wrap(err), middleware.RequestID(r))
		return
	}
	projects, total, err := h.lister.List(r.Context(), query)
	if err != nil {
		httperr.Write(w, err, middleware.RequestID(r))
		return
	}
	page, err := pagination.Page(query, total, h.paginationSecret)
	if err != nil {
		httperr.Write(w, err, middleware.RequestID(r))
		return
	}
	utils.WriteJSON(w, http.StatusOK, projectdto.FromList(projects, page))
}

func (h *Handler) listQuery(r *http.Request) (paging.ListQuery, error) {
	values := r.URL.Query()
	cursor := values.Get("cursor")
	if cursor != "" {
		if values.Get("limit") != "" || values.Get("sort") != "" || values.Get("name_like") != "" || values.Get("monitoring_status_in") != "" {
			return paging.ListQuery{}, apperr.ErrInvalidProjectCommand
		}
		parsed, err := pagination.Parse(cursor, h.paginationSecret)
		if err != nil {
			return paging.ListQuery{}, err
		}
		return paging.ListQuery{
			Limit: parsed.Limit, Offset: parsed.Offset, Sort: parsed.Sort, NameLike: parsed.NameLike,
			MonitoringStatusIn: parseStatuses(strings.Join(parsed.MonitoringStatusIn, ",")),
		}, nil
	}
	limit := 20
	if raw := values.Get("limit"); raw != "" {
		parsed, err := strconv.Atoi(raw)
		if err != nil {
			return paging.ListQuery{}, err
		}
		limit = parsed
	}
	return paging.ListQuery{
		Limit: limit, Sort: values.Get("sort"), NameLike: values.Get("name_like"),
		MonitoringStatusIn: parseStatuses(values.Get("monitoring_status_in")),
	}, nil
}

func parseStatuses(raw string) []domain.MonitoringStatus {
	if strings.TrimSpace(raw) == "" {
		return nil
	}
	parts := strings.Split(raw, ",")
	statuses := make([]domain.MonitoringStatus, 0, len(parts))
	for _, part := range parts {
		statuses = append(statuses, domain.MonitoringStatus(strings.TrimSpace(part)))
	}
	return statuses
}
