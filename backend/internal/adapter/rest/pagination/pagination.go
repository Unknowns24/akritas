package pagination

import (
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

var (
	ErrInvalidConfiguration = errors.New("invalid Uker pagination configuration")
	ErrInvalidQuery         = errors.New("invalid pagination query")
)

type Config struct {
	Secret []byte
	TTL    time.Duration
}

func NewConfig(secret []byte, ttl time.Duration) (Config, error) {
	if len(secret) < 32 || ttl <= 0 {
		return Config{}, ErrInvalidConfiguration
	}
	return Config{Secret: append([]byte(nil), secret...), TTL: ttl}, nil
}

type Policy struct {
	AllowedFilters map[string]struct{}
	AllowedSorts   map[string]struct{}
	DefaultSort    []ukerpagination.SortExpression
	Scope          string
	BoundaryField  string
}

func Parse(r *http.Request, config Config, policy Policy) (paging.Params, error) {
	if r == nil || len(config.Secret) < 32 || config.TTL <= 0 {
		return paging.Params{}, ErrInvalidConfiguration
	}
	values := cloneValues(r.URL.Query())
	if err := validateRawQuery(values, policy); err != nil {
		return paging.Params{}, err
	}
	if values.Get("cursor") == "" && values.Get("sort") == "" {
		sorts := policy.DefaultSort
		if policy.BoundaryField != "" {
			sorts = []ukerpagination.SortExpression{{Field: policy.BoundaryField, Direction: ukerpagination.DirectionAsc}}
		}
		if len(sorts) > 0 {
			values.Set("sort", encodeSort(sorts))
		}
	}
	params, err := ukerpagination.ParseWithSecurity(values, config.Secret, config.TTL)
	if err != nil {
		return paging.Params{}, err
	}
	if err := validateParams(&params, policy); err != nil {
		return paging.Params{}, err
	}
	if policy.Scope != "" {
		const scopeFilter = "integration_scope_eq"
		if current, exists := params.Filters[scopeFilter]; exists && current != policy.Scope {
			return paging.Params{}, ErrInvalidQuery
		}
		params.Filters[scopeFilter] = policy.Scope
	}
	return params, nil
}

func BuildPage[T any](params paging.Params, result paging.Slice[T], config Config, extractor ukerpagination.CursorExtractor[T]) (ukerpagination.PagingResponse[T], error) {
	return ukerpagination.BuildPageSigned(params, result.Items, params.Limit, result.Total, extractor, config.Secret)
}

func BuildProviderPage[T any](params paging.Params, result paging.Slice[T], config Config) (ukerpagination.PagingResponse[T], error) {
	var nextCursor, prevCursor string
	var err error
	if len(result.NextBoundary) > 0 {
		nextCursor, err = ukerpagination.BuildNextCursorSigned(params, result.NextBoundary, config.Secret)
		if err != nil {
			return ukerpagination.PagingResponse[T]{}, err
		}
	}
	if len(result.PrevBoundary) > 0 {
		prevCursor, err = ukerpagination.BuildPrevCursorSigned(params, result.PrevBoundary, config.Secret)
		if err != nil {
			return ukerpagination.PagingResponse[T]{}, err
		}
	}
	return ukerpagination.NewPage(result.Items, params.Limit, result.Total, nextCursor != "", nextCursor, prevCursor), nil
}

func validateRawQuery(values url.Values, policy Policy) error {
	for key := range values {
		if key == "cursor" || key == "limit" || key == "sort" {
			continue
		}
		if _, allowed := policy.AllowedFilters[key]; !allowed {
			return ErrInvalidQuery
		}
	}
	return nil
}

func validateParams(params *paging.Params, policy Policy) error {
	for key := range params.Filters {
		if key == "integration_scope_eq" && policy.Scope != "" {
			continue
		}
		if _, allowed := policy.AllowedFilters[key]; !allowed {
			return ErrInvalidQuery
		}
	}
	for _, sort := range params.Sort {
		if policy.BoundaryField != "" && sort.Field == policy.BoundaryField {
			continue
		}
		if _, allowed := policy.AllowedSorts[sort.Field]; !allowed {
			return ErrInvalidQuery
		}
	}
	if policy.BoundaryField == "" && len(params.Sort) > 0 && !hasSort(params.Sort, "id") {
		direction := params.Sort[len(params.Sort)-1].Direction
		params.Sort = append(params.Sort, ukerpagination.SortExpression{Field: "id", Direction: direction})
	}
	return nil
}

func hasSort(sorts []ukerpagination.SortExpression, field string) bool {
	for _, sort := range sorts {
		if sort.Field == field {
			return true
		}
	}
	return false
}

func encodeSort(sorts []ukerpagination.SortExpression) string {
	values := make([]string, 0, len(sorts))
	for _, sort := range sorts {
		values = append(values, sort.Field+":"+string(sort.Direction))
	}
	return strings.Join(values, ",")
}

func cloneValues(input url.Values) url.Values {
	output := make(url.Values, len(input))
	for key, values := range input {
		output[key] = append([]string(nil), values...)
	}
	return output
}
