package dokploy

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
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
		log.Printf("dokploy: read logs request failed source_type=%s resource_id=%s service=%s instance=%s endpoint=%s error=%v", request.Source.Type, request.Source.ResourceIdentifier, request.Source.ServiceName, request.Source.InstanceIdentifier, endpoint, err)
		return nil, normalizeProviderError(err)
	}
	records, err := ParseLogs(body)
	if err != nil {
		log.Printf("dokploy: parse logs failed source_type=%s resource_id=%s service=%s instance=%s endpoint=%s bytes=%d error=%v", request.Source.Type, request.Source.ResourceIdentifier, request.Source.ServiceName, request.Source.InstanceIdentifier, endpoint, len(body), err)
		return nil, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	log.Printf("dokploy: parsed logs source_type=%s resource_id=%s service=%s instance=%s endpoint=%s records=%d", request.Source.Type, request.Source.ResourceIdentifier, request.Source.ServiceName, request.Source.InstanceIdentifier, endpoint, len(records))
	return records, nil
}

type containerDTO struct {
	ID     string
	State  string
	Labels map[string]string
}

func (c containerDTO) identifier() string {
	return c.ID
}

func (c *containerDTO) UnmarshalJSON(body []byte) error {
	var fields map[string]any
	if err := json.Unmarshal(body, &fields); err != nil {
		return err
	}
	c.ID = firstString(fields, "Id", "ID", "id", "containerId", "containerID")
	c.State = firstString(fields, "State", "state", "status", "Status")
	c.Labels = firstStringMap(fields, "Labels", "labels")
	return nil
}

func firstString(fields map[string]any, keys ...string) string {
	for _, key := range keys {
		if value, ok := fields[key].(string); ok {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func firstStringMap(fields map[string]any, keys ...string) map[string]string {
	for _, key := range keys {
		raw, ok := fields[key]
		if !ok {
			continue
		}
		result := map[string]string{}
		switch value := raw.(type) {
		case map[string]string:
			return value
		case map[string]any:
			for label, labelValue := range value {
				if text, ok := labelValue.(string); ok {
					result[label] = text
				}
			}
		}
		if len(result) > 0 {
			return result
		}
	}
	return nil
}

func (c *Client) resolveComposeContainer(ctx context.Context, base, credential string, source domain.DokploySource) (string, error) {
	values := url.Values{"appName": {source.InstanceIdentifier}, "appType": {string(source.RuntimeType)}}
	if source.ProviderServerID != "" {
		values.Set("serverId", source.ProviderServerID)
	}
	log.Printf("dokploy: resolving compose container compose_id=%s service=%s app_name=%s runtime=%s provider_server_id_configured=%t", source.ResourceIdentifier, source.ServiceName, source.InstanceIdentifier, source.RuntimeType, source.ProviderServerID != "")
	body, err := c.do(ctx, base+"/api/docker.getContainersByAppNameMatch?"+values.Encode(), credential)
	if err != nil {
		log.Printf("dokploy: resolve compose container request failed compose_id=%s service=%s app_name=%s error=%v", source.ResourceIdentifier, source.ServiceName, source.InstanceIdentifier, err)
		return "", normalizeProviderError(err)
	}
	containers, err := decodeContainers(body)
	if err != nil {
		log.Printf("dokploy: resolve compose container response was not parseable compose_id=%s service=%s app_name=%s bytes=%d", source.ResourceIdentifier, source.ServiceName, source.InstanceIdentifier, len(body))
		return "", domain.ErrIntegrationUnavailable.Wrap(err)
	}
	matchingCandidates := make([]string, 0)
	runningMatchingCandidates := make([]string, 0)
	allCandidates := make([]string, 0)
	runningCandidates := make([]string, 0)
	serviceLabels := make([]string, 0)
	states := make([]string, 0)
	running := 0
	matching := 0
	for _, container := range containers {
		state := strings.TrimSpace(container.State)
		identifier := strings.TrimSpace(container.identifier())
		if identifier != "" {
			allCandidates = append(allCandidates, identifier)
		}
		if label := composeServiceLabel(container.Labels, source); label != "" {
			serviceLabels = append(serviceLabels, label)
		}
		if state != "" {
			states = append(states, state)
		}
		isRunning := strings.EqualFold(state, "running")
		if isRunning {
			running++
			if identifier != "" {
				runningCandidates = append(runningCandidates, identifier)
			}
		}
		if matchesComposeService(container.Labels, source) {
			matching++
			if identifier != "" {
				matchingCandidates = append(matchingCandidates, identifier)
				if isRunning {
					runningMatchingCandidates = append(runningMatchingCandidates, identifier)
				}
			}
		}
	}
	if len(runningMatchingCandidates) > 0 {
		sort.Strings(runningMatchingCandidates)
		log.Printf("dokploy: resolved running compose container compose_id=%s service=%s app_name=%s containers=%d running=%d service_matches=%d selected_container=%s states=%q service_labels=%q", source.ResourceIdentifier, source.ServiceName, source.InstanceIdentifier, len(containers), running, matching, runningMatchingCandidates[0], strings.Join(states, ","), strings.Join(serviceLabels, ","))
		return runningMatchingCandidates[0], nil
	}
	if len(matchingCandidates) > 0 {
		sort.Strings(matchingCandidates)
		log.Printf("dokploy: resolved non-running compose container for logs compose_id=%s service=%s app_name=%s containers=%d running=%d service_matches=%d selected_container=%s states=%q service_labels=%q", source.ResourceIdentifier, source.ServiceName, source.InstanceIdentifier, len(containers), running, matching, matchingCandidates[0], strings.Join(states, ","), strings.Join(serviceLabels, ","))
		return matchingCandidates[0], nil
	}
	if len(runningCandidates) == 1 {
		log.Printf("dokploy: using single running compose container fallback compose_id=%s service=%s app_name=%s selected_container=%s states=%q service_labels=%q", source.ResourceIdentifier, source.ServiceName, source.InstanceIdentifier, runningCandidates[0], strings.Join(states, ","), strings.Join(serviceLabels, ","))
		return runningCandidates[0], nil
	}
	if len(allCandidates) == 1 {
		log.Printf("dokploy: using single compose container fallback for logs compose_id=%s service=%s app_name=%s selected_container=%s states=%q service_labels=%q", source.ResourceIdentifier, source.ServiceName, source.InstanceIdentifier, allCandidates[0], strings.Join(states, ","), strings.Join(serviceLabels, ","))
		return allCandidates[0], nil
	}
	log.Printf("dokploy: no compose container could be selected compose_id=%s service=%s app_name=%s containers=%d running=%d service_matches=%d states=%q service_labels=%q", source.ResourceIdentifier, source.ServiceName, source.InstanceIdentifier, len(containers), running, matching, strings.Join(states, ","), strings.Join(serviceLabels, ","))
	return "", domain.ErrDokployContainerUnavailable
}

func decodeContainers(body []byte) ([]containerDTO, error) {
	var containers []containerDTO
	if err := json.Unmarshal(body, &containers); err == nil {
		return containers, nil
	}
	var wrapper struct {
		Items      []containerDTO `json:"items"`
		Data       []containerDTO `json:"data"`
		Containers []containerDTO `json:"containers"`
	}
	if err := json.Unmarshal(body, &wrapper); err != nil {
		return nil, err
	}
	switch {
	case wrapper.Items != nil:
		return wrapper.Items, nil
	case wrapper.Data != nil:
		return wrapper.Data, nil
	default:
		return wrapper.Containers, nil
	}
}

func matchesComposeService(labels map[string]string, source domain.DokploySource) bool {
	service := composeServiceLabel(labels, source)
	if source.RuntimeType == domain.DokployRuntimeStack {
		return service == source.InstanceIdentifier+"_"+source.ServiceName || service == source.ServiceName
	}
	return service == source.ServiceName
}

func composeServiceLabel(labels map[string]string, source domain.DokploySource) string {
	if source.RuntimeType == domain.DokployRuntimeStack {
		return strings.TrimSpace(labels["com.docker.swarm.service.name"])
	}
	return strings.TrimSpace(labels["com.docker.compose.service"])
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
