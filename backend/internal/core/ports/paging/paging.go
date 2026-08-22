package paging

import ukerpagination "github.com/unknowns24/uker/uker/pagination"

// Params is the project-wide pagination contract. It intentionally aliases
// Uker so every port uses the same parsing, cursor and sort/filter semantics.
type Params = ukerpagination.Params

// Slice contains the provider/repository result before REST builds signed
// cursors. Boundaries are values, never encoded or signed tokens.
type Slice[T any] struct {
	Items        []T
	Total        int64
	NextBoundary map[string]string
	PrevBoundary map[string]string
}
