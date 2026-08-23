package timeline

import (
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func TestTimelineOrderQualifiesCalculatedColumnsWithTimelineAlias(t *testing.T) {
	params := paging.Params{Sort: []ukerpagination.SortExpression{
		{Field: "occurred_at", Direction: ukerpagination.DirectionAsc},
		{Field: "id", Direction: ukerpagination.DirectionDesc},
	}}

	got := timelineOrder(params)
	want := "timeline.occurred_at ASC, timeline.id DESC"
	if got != want {
		t.Fatalf("order = %q, want %q", got, want)
	}
}

func TestTimelineOrderIgnoresUnsupportedFields(t *testing.T) {
	params := paging.Params{Sort: []ukerpagination.SortExpression{
		{Field: "project_id", Direction: ukerpagination.DirectionDesc},
		{Field: "occurred_at", Direction: ukerpagination.DirectionDesc},
	}}

	got := timelineOrder(params)
	want := "timeline.occurred_at DESC"
	if got != want {
		t.Fatalf("order = %q, want %q", got, want)
	}
}
