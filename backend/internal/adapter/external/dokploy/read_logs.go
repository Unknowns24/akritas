package dokploy

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (c *Client) FetchLogs(ctx context.Context, request portsout.LogFetchRequest) ([]portsout.RawLogRecord, error) {
	if c.credentials == nil || request.Tail < 1 || request.Tail > 10000 || strings.TrimSpace(request.Since) == "" || request.Source.Validate() != nil || request.Server.ID != request.Source.DokployServerID {
		return nil, domain.ErrIntegrationUnavailable
	}
	base, credential, err := c.providerAccess(ctx, request.Server)
	if err != nil {
		return nil, err
	}
	defer wipe(credential)
	query := url.Values{"tail": {strconv.Itoa(request.Tail)}, "since": {request.Since}}
	endpoint := "/api/application.readLogs"
	if request.Source.Type == domain.DokploySourceApplication {
		query.Set("applicationId", request.Source.ResourceIdentifier)
	} else {
		containerID, resolveErr := c.resolveComposeContainer(ctx, base, string(credential), request.Source)
		if resolveErr != nil {
			return nil, resolveErr
		}
		endpoint = "/api/compose.readLogs"
		query.Set("composeId", request.Source.ResourceIdentifier)
		query.Set("containerId", containerID)
	}
	body, err := c.do(ctx, base+endpoint+"?"+query.Encode(), string(credential))
	if err != nil {
		return nil, normalizeProviderError(err)
	}
	records, err := ParseLogs(body)
	if err != nil {
		return nil, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	return records, nil
}

type containerDTO struct {
	ID     string            `json:"Id"`
	IDAlt  string            `json:"ID"`
	State  string            `json:"State"`
	Labels map[string]string `json:"Labels"`
}

func (c containerDTO) identifier() string {
	if c.ID != "" {
		return c.ID
	}
	return c.IDAlt
}

func (c *Client) resolveComposeContainer(ctx context.Context, base, credential string, source domain.DokploySource) (string, error) {
	values := url.Values{"appName": {source.InstanceIdentifier}, "appType": {string(source.RuntimeType)}}
	if source.ProviderServerID != "" {
		values.Set("serverId", source.ProviderServerID)
	}
	body, err := c.do(ctx, base+"/api/docker.getContainersByAppNameMatch?"+values.Encode(), credential)
	if err != nil {
		return "", normalizeProviderError(err)
	}
	var containers []containerDTO
	if json.Unmarshal(body, &containers) != nil {
		var wrapper struct {
			Items []containerDTO `json:"items"`
		}
		if json.Unmarshal(body, &wrapper) != nil {
			return "", domain.ErrIntegrationUnavailable
		}
		containers = wrapper.Items
	}
	candidates := make([]string, 0)
	for _, container := range containers {
		if !strings.EqualFold(strings.TrimSpace(container.State), "running") || !matchesComposeService(container.Labels, source) {
			continue
		}
		if identifier := strings.TrimSpace(container.identifier()); identifier != "" {
			candidates = append(candidates, identifier)
		}
	}
	if len(candidates) == 0 {
		return "", domain.ErrDokployContainerUnavailable
	}
	sort.Strings(candidates)
	return candidates[0], nil
}

func matchesComposeService(labels map[string]string, source domain.DokploySource) bool {
	if source.RuntimeType == domain.DokployRuntimeStack {
		return labels["com.docker.stack.namespace"] == source.InstanceIdentifier &&
			labels["com.docker.swarm.service.name"] == source.InstanceIdentifier+"_"+source.ServiceName
	}
	return labels["com.docker.compose.project"] == source.InstanceIdentifier && labels["com.docker.compose.service"] == source.ServiceName
}

func ParseLogs(body []byte) ([]portsout.RawLogRecord, error) {
	var payload string
	if err := json.Unmarshal(body, &payload); err != nil {
		payload = string(body)
	}
	scanner := bufio.NewScanner(strings.NewReader(payload))
	scanner.Buffer(make([]byte, 64*1024), maximumBodyBytes)
	records := make([]portsout.RawLogRecord, 0)
	ordinals := make(map[int64]int)
	for scanner.Scan() {
		line := scanner.Text()
		separator := strings.IndexByte(line, ' ')
		if separator <= 0 {
			if len(records) > 0 {
				records[len(records)-1].Message += "\n" + line
				records[len(records)-1].ContentHash = hashContent(records[len(records)-1].Message)
			}
			continue
		}
		timestamp, err := time.Parse(time.RFC3339Nano, line[:separator])
		if err != nil {
			if len(records) > 0 {
				records[len(records)-1].Message += "\n" + line
				records[len(records)-1].ContentHash = hashContent(records[len(records)-1].Message)
			}
			continue
		}
		message := strings.TrimSuffix(line[separator+1:], "\r")
		key := timestamp.UnixNano()
		ordinal := ordinals[key]
		ordinals[key] = ordinal + 1
		records = append(records, portsout.RawLogRecord{Timestamp: timestamp.UTC(), Ordinal: ordinal, ContentHash: hashContent(message), Message: message})
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return records, nil
}

func hashContent(value string) string {
	sum := sha256.Sum256([]byte(value))
	return hex.EncodeToString(sum[:])
}

var _ portsout.LogSource = (*Client)(nil)
