package dokploy

import "github.com/Unknowns24/akritas/backend/internal/core/ports/paging"

func providerBoundary(params paging.Params, field string) string {
	if params.Cursor == nil {
		return ""
	}
	if value := params.Cursor.After[field]; value != "" {
		return value
	}
	return params.Cursor.Before[field]
}
