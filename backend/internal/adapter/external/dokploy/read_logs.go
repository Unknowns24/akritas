package dokploy

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
)

func (c *Client) FetchLogs(ctx context.Context, request portsout.LogFetchRequest) ([]portsout.RawLogRecord, error) {
	if c.credentials == nil || request.Tail < 1 || request.Tail > 10000 || strings.TrimSpace(request.Since) == "" || request.Application.ApplicationIdentifier == "" || request.Server.ID != request.Application.DokployServerID {
		return nil, domain.ErrIntegrationUnavailable
	}
	credential, err := c.credentials.Get(ctx, portsout.CredentialOwnerDokployServer, request.Server.ID, portsout.SecretKindDokployAPIKey)
	if err != nil {
		return nil, domain.ErrIntegrationUnavailable
	}
	defer wipe(credential)
	base, err := c.normalizeAndValidateURL(ctx, request.Server.BaseURL)
	if err != nil {
		return nil, domain.ErrIntegrationUnavailable
	}
	query := url.Values{"applicationId": {request.Application.ApplicationIdentifier}, "tail": {strconv.Itoa(request.Tail)}, "since": {request.Since}}
	body, err := c.do(ctx, base+"/api/application.readLogs?"+query.Encode(), string(credential))
	if err != nil {
		return nil, normalizeProviderError(err)
	}
	records, err := ParseLogs(body)
	if err != nil {
		return nil, domain.ErrIntegrationUnavailable.Wrap(err)
	}
	return records, nil
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
