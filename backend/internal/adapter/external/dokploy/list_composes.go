package dokploy

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/google/uuid"
)

type composeSearchResponse struct {
	Items []composeDTO `json:"items"`
	Total int          `json:"total"`
}

type composeDTO struct {
	ComposeID     string          `json:"composeId"`
	AppName       string          `json:"appName"`
	Name          string          `json:"name"`
	EnvironmentID string          `json:"environmentId"`
	Environment   json.RawMessage `json:"environment"`
	ComposeStatus string          `json:"composeStatus"`
	Status        string          `json:"status"`
	ComposeType   string          `json:"composeType"`
	ServerID      string          `json:"serverId"`
}

func (c *Client) ListComposes(ctx context.Context, server domain.DokployServer, params paging.Params) (paging.Slice[domain.DokployCompose], error) {
	if params.Limit < 1 {
		params.Limit = 25
	}
	if params.Limit > 100 || c.credentials == nil {
		return paging.Slice[domain.DokployCompose]{}, domain.ErrInvalidDokployCompose
	}
	offset := 0
	if position := providerBoundary(params, "provider_offset"); position != "" {
		parsed, err := strconv.Atoi(position)
		if err != nil || parsed < 0 {
			return paging.Slice[domain.DokployCompose]{}, domain.ErrInvalidDokployCompose
		}
		offset = parsed
	}
	base, credential, err := c.providerAccess(ctx, server)
	if err != nil {
		return paging.Slice[domain.DokployCompose]{}, err
	}
	defer wipe(credential)
	values := url.Values{"limit": {strconv.Itoa(params.Limit)}, "offset": {strconv.Itoa(offset)}}
	if value := strings.TrimSpace(params.Filters["name_like"]); value != "" {
		values.Set("q", value)
	}
	body, err := c.do(ctx, base+"/api/compose.search?"+values.Encode(), string(credential))
	if err != nil {
		return paging.Slice[domain.DokployCompose]{}, normalizeProviderError(err)
	}
	var payload composeSearchResponse
	if err := json.Unmarshal(body, &payload); err != nil {
		if directErr := json.Unmarshal(body, &payload.Items); directErr != nil {
			return paging.Slice[domain.DokployCompose]{}, domain.ErrIntegrationUnavailable
		}
		payload.Total = len(payload.Items)
	}
	items := make([]domain.DokployCompose, 0, len(payload.Items))
	for _, item := range payload.Items {
		mapped, mapErr := item.domain(server.ID)
		if mapErr == nil {
			items = append(items, mapped)
		}
	}
	if len(payload.Items) > 0 && len(items) == 0 {
		return paging.Slice[domain.DokployCompose]{}, domain.ErrIntegrationUnavailable
	}
	total := payload.Total
	if total < len(items) {
		total = len(items)
	}
	result := paging.Slice[domain.DokployCompose]{Items: items, Total: int64(total)}
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

func (c *Client) GetCompose(ctx context.Context, server domain.DokployServer, identifier string) (domain.DokployCompose, error) {
	identifier = strings.TrimSpace(identifier)
	if identifier == "" || c.credentials == nil {
		return domain.DokployCompose{}, domain.ErrIntegrationNotFound
	}
	base, credential, err := c.providerAccess(ctx, server)
	if err != nil {
		return domain.DokployCompose{}, err
	}
	defer wipe(credential)
	body, err := c.do(ctx, base+"/api/compose.one?"+url.Values{"composeId": {identifier}}.Encode(), string(credential))
	if err != nil {
		return domain.DokployCompose{}, normalizeComposeProviderError(err)
	}
	var payload composeDTO
	if json.Unmarshal(body, &payload) != nil || payload.ComposeID != identifier {
		return domain.DokployCompose{}, domain.ErrIntegrationUnavailable
	}
	return payload.domain(server.ID)
}

func (c *Client) ListComposeServices(ctx context.Context, server domain.DokployServer, composeID string, refresh bool) ([]domain.DokployComposeService, error) {
	composeID = strings.TrimSpace(composeID)
	if composeID == "" || c.credentials == nil {
		return nil, domain.ErrInvalidDokployCompose
	}
	base, credential, err := c.providerAccess(ctx, server)
	if err != nil {
		return nil, err
	}
	defer wipe(credential)
	loadType := "cache"
	if refresh {
		loadType = "fetch"
	}
	values := url.Values{"composeId": {composeID}, "type": {loadType}}
	body, err := c.do(ctx, base+"/api/compose.loadServices?"+values.Encode(), string(credential))
	if err != nil {
		return nil, normalizeComposeProviderError(err)
	}
	var names []string
	if json.Unmarshal(body, &names) != nil {
		return nil, domain.ErrIntegrationUnavailable
	}
	unique := make(map[string]struct{}, len(names))
	for _, name := range names {
		name = strings.TrimSpace(name)
		if name != "" {
			unique[name] = struct{}{}
		}
	}
	names = names[:0]
	for name := range unique {
		names = append(names, name)
	}
	sort.Strings(names)
	result := make([]domain.DokployComposeService, 0, len(names))
	for _, name := range names {
		value, mapErr := domain.NewDokployComposeService(server.ID, composeID, name)
		if mapErr != nil {
			return nil, mapErr
		}
		result = append(result, value)
	}
	return result, nil
}

func (c *Client) providerAccess(ctx context.Context, server domain.DokployServer) (string, []byte, error) {
	credential, err := c.credentials.Get(ctx, portsout.CredentialOwnerDokployServer, server.ID, portsout.SecretKindDokployAPIKey)
	if err != nil {
		return "", nil, domain.ErrIntegrationUnavailable
	}
	base, err := c.normalizeAndValidateURL(ctx, server.BaseURL)
	if err != nil {
		wipe(credential)
		return "", nil, domain.ErrIntegrationUnavailable
	}
	return base, credential, nil
}

func normalizeComposeProviderError(err error) error {
	var providerErr *providerError
	if errors.As(err, &providerErr) && providerErr.Status == http.StatusNotFound {
		return domain.ErrIntegrationNotFound
	}
	return normalizeProviderError(err)
}

func (d composeDTO) domain(serverID uuid.UUID) (domain.DokployCompose, error) {
	displayName := strings.TrimSpace(d.Name)
	if displayName == "" {
		displayName = strings.TrimSpace(d.AppName)
	}
	statusValue := d.ComposeStatus
	if statusValue == "" {
		statusValue = d.Status
	}
	return domain.NewDokployCompose(serverID, d.ComposeID, d.AppName, displayName, d.EnvironmentID,
		d.environmentName(), mapComposeStatus(statusValue), mapRuntimeType(d.ComposeType), d.ServerID)
}

func (d composeDTO) environmentName() string {
	if len(d.Environment) == 0 || string(d.Environment) == "null" {
		return ""
	}
	var value string
	if json.Unmarshal(d.Environment, &value) == nil {
		return value
	}
	var object struct {
		Name string `json:"name"`
	}
	_ = json.Unmarshal(d.Environment, &object)
	return object.Name
}

func mapComposeStatus(value string) domain.DokploySourceStatus {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "running":
		return domain.DokploySourceRunning
	case "stopped":
		return domain.DokploySourceStopped
	case "error", "degraded":
		return domain.DokploySourceDegraded
	default:
		return domain.DokploySourceUnknown
	}
}

func mapRuntimeType(value string) domain.DokployRuntimeType {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "stack":
		return domain.DokployRuntimeStack
	case "docker-compose", "compose":
		return domain.DokployRuntimeCompose
	default:
		return ""
	}
}
