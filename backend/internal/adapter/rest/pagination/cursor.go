package pagination

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
)

type Cursor struct {
	Offset             int      `json:"offset"`
	Limit              int      `json:"limit"`
	Sort               string   `json:"sort"`
	NameLike           string   `json:"name_like"`
	MonitoringStatusIn []string `json:"monitoring_status_in"`
}

func Sign(cursor Cursor, secret string) (string, error) {
	payload, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(payload)
	return base64.RawURLEncoding.EncodeToString(payload) + "." + base64.RawURLEncoding.EncodeToString(mac.Sum(nil)), nil
}

func Parse(raw, secret string) (Cursor, error) {
	parts := strings.Split(raw, ".")
	if len(parts) != 2 {
		return Cursor{}, errors.New("invalid cursor")
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Cursor{}, err
	}
	signature, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Cursor{}, err
	}
	mac := hmac.New(sha256.New, []byte(secret))
	_, _ = mac.Write(payload)
	if !hmac.Equal(signature, mac.Sum(nil)) {
		return Cursor{}, errors.New("invalid cursor signature")
	}
	var cursor Cursor
	if err := json.Unmarshal(payload, &cursor); err != nil {
		return Cursor{}, err
	}
	return cursor, nil
}

func Page(query paging.ListQuery, total int64, secret string) (paging.Page, error) {
	limit := paging.NormalizeLimit(query.Limit)
	page := paging.Page{Limit: limit, Total: total, HasMore: int64(query.Offset+limit) < total}
	if page.HasMore {
		next, err := Sign(Cursor{
			Offset: query.Offset + limit, Limit: limit, Sort: query.Sort,
			NameLike: query.NameLike, MonitoringStatusIn: statuses(query),
		}, secret)
		if err != nil {
			return paging.Page{}, err
		}
		page.NextCursor = next
	}
	if query.Offset > 0 {
		prevOffset := query.Offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		prev, err := Sign(Cursor{
			Offset: prevOffset, Limit: limit, Sort: query.Sort,
			NameLike: query.NameLike, MonitoringStatusIn: statuses(query),
		}, secret)
		if err != nil {
			return paging.Page{}, err
		}
		page.PrevCursor = prev
	}
	return page, nil
}

func statuses(query paging.ListQuery) []string {
	values := make([]string, 0, len(query.MonitoringStatusIn))
	for _, status := range query.MonitoringStatusIn {
		values = append(values, string(status))
	}
	return values
}

func FormatDebug(cursor Cursor) string {
	return fmt.Sprintf("%d/%d", cursor.Offset, cursor.Limit)
}
