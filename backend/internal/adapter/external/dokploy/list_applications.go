package dokploy

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
)

func (c *Client) ListApplications(ctx context.Context, server domain.DokployServer, params paging.Params) (paging.Slice[domain.DokployApplication], error) {
	if params.Limit < 1 {
		params.Limit = 25
	}
	if params.Limit > 100 || c.credentials == nil {
		return paging.Slice[domain.DokployApplication]{}, domain.ErrInvalidDokployApplication
	}
	offset := 0
	if position := providerBoundary(params, "provider_offset"); position != "" {
		parsed, err := strconv.Atoi(position)
		if err != nil || parsed < 0 {
			return paging.Slice[domain.DokployApplication]{}, domain.ErrInvalidDokployApplication
		}
		offset = parsed
	}
	credential, err := c.credentials.Get(ctx, portsout.CredentialOwnerDokployServer, server.ID, portsout.SecretKindDokployAPIKey)
	if err != nil {
		return paging.Slice[domain.DokployApplication]{}, domain.ErrIntegrationUnavailable
	}
	defer wipe(credential)
	base, err := c.normalizeAndValidateURL(ctx, server.BaseURL)
	if err != nil {
		return paging.Slice[domain.DokployApplication]{}, domain.ErrIntegrationUnavailable
	}
	values := url.Values{"limit": {strconv.Itoa(params.Limit)}, "offset": {strconv.Itoa(offset)}}
	if nameLike := params.Filters["name_like"]; nameLike != "" {
		values.Set("q", nameLike)
	}
	response, err := c.do(ctx, base+"/api/application.search?"+values.Encode(), string(credential))
	if err != nil {
		return paging.Slice[domain.DokployApplication]{}, normalizeProviderError(err)
	}
	var payload applicationSearchResponse
	if err := json.Unmarshal(response, &payload); err != nil {
		var direct []applicationDTO
		if directErr := json.Unmarshal(response, &direct); directErr != nil {
			return paging.Slice[domain.DokployApplication]{}, domain.ErrIntegrationUnavailable
		}
		payload.Items = direct
		payload.Total = len(direct)
	}
	items := make([]domain.DokployApplication, 0, len(payload.Items))
	for _, item := range payload.Items {
		if environment := params.Filters["environment_eq"]; environment != "" && !strings.EqualFold(item.environmentName(), environment) {
			continue
		}
		mapped, mapErr := domain.NewDokployApplication(server.ID, item.ApplicationID, item.AppName, item.displayName(), item.environmentName(), mapStatus(item.status()))
		if mapErr == nil {
			items = append(items, mapped)
		}
	}
	if len(payload.Items) > 0 && len(items) == 0 {
		return paging.Slice[domain.DokployApplication]{}, domain.ErrIntegrationUnavailable
	}
	total := payload.Total
	if total < len(items) {
		total = len(items)
	}
	result := paging.Slice[domain.DokployApplication]{Items: items, Total: int64(total)}
	if offset+len(payload.Items) < total || len(payload.Items) == params.Limit {
		result.NextBoundary = map[string]string{"provider_offset": strconv.Itoa(offset + len(payload.Items))}
	}
	if offset > 0 {
		previous := offset - params.Limit
		if previous < 0 {
			previous = 0
		}
		result.PrevBoundary = map[string]string{"provider_offset": strconv.Itoa(previous)}
	}
	return result, nil
}

type applicationSearchResponse struct {
	Items []applicationDTO `json:"items"`
	Total int              `json:"total"`
}

type applicationDTO struct {
	ApplicationID     string          `json:"applicationId"`
	AppName           string          `json:"appName"`
	Name              string          `json:"name"`
	Status            string          `json:"status"`
	ApplicationStatus string          `json:"applicationStatus"`
	Environment       json.RawMessage `json:"environment"`
}

func (a applicationDTO) displayName() string {
	if strings.TrimSpace(a.Name) != "" {
		return a.Name
	}
	return a.AppName
}

func (a applicationDTO) status() string {
	if a.ApplicationStatus != "" {
		return a.ApplicationStatus
	}
	return a.Status
}

func (a applicationDTO) environmentName() string {
	if len(a.Environment) == 0 || string(a.Environment) == "null" {
		return ""
	}
	var value string
	if json.Unmarshal(a.Environment, &value) == nil {
		return value
	}
	var object struct {
		Name string `json:"name"`
	}
	_ = json.Unmarshal(a.Environment, &object)
	return object.Name
}

func mapStatus(status string) domain.DokployApplicationStatus {
	switch strings.ToLower(status) {
	case "running":
		return domain.DokployApplicationRunning
	case "stopped":
		return domain.DokployApplicationStopped
	case "degraded":
		return domain.DokployApplicationDegraded
	default:
		return domain.DokployApplicationUnknown
	}
}
