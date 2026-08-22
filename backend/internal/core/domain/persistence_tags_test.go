package domain

import (
	"reflect"
	"strings"
	"testing"
)

func TestPersistedIntegrationEntitiesExposeOnlyDeclarativeGORMMetadata(t *testing.T) {
	t.Parallel()

	for _, entity := range []any{GitHubAccount{}, DokployServer{}} {
		typeOf := reflect.TypeOf(entity)
		id, ok := typeOf.FieldByName("ID")
		if !ok || !strings.Contains(id.Tag.Get("gorm"), "primaryKey") {
			t.Fatalf("%s.ID must declare its GORM primary key metadata", typeOf.Name())
		}
		createdAt, ok := typeOf.FieldByName("CreatedAt")
		if !ok || !strings.Contains(createdAt.Tag.Get("gorm"), "column:created_at") {
			t.Fatalf("%s.CreatedAt must map the existing schema", typeOf.Name())
		}
	}
}
