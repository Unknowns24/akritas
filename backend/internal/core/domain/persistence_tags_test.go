package domain

import (
	"reflect"
	"strings"
	"testing"
)

func TestPersistedIntegrationEntitiesExposeOnlyDeclarativeGORMMetadata(t *testing.T) {
	t.Parallel()

	for _, entity := range []any{GitHubAccount{}, DokployServer{}, Administrator{}, AdministratorSession{}, PendingEnrollment{}, Project{}} {
		typeOf := reflect.TypeOf(entity)
		id, ok := typeOf.FieldByName("ID")
		if !ok || !strings.Contains(id.Tag.Get("gorm"), "primaryKey") {
			t.Fatalf("%s.ID must declare its GORM primary key metadata", typeOf.Name())
		}
		if createdAt, ok := typeOf.FieldByName("CreatedAt"); ok && !strings.Contains(createdAt.Tag.Get("gorm"), "column:created_at") {
			t.Fatalf("%s.CreatedAt must map the existing schema", typeOf.Name())
		}
	}
}

func TestProjectDomainEntityDoesNotContainSecretMaterial(t *testing.T) {
	t.Parallel()

	typeOf := reflect.TypeOf(Project{})
	for index := 0; index < typeOf.NumField(); index++ {
		name := strings.ToLower(typeOf.Field(index).Name)
		if strings.Contains(name, "password") || strings.Contains(name, "secret") || strings.Contains(name, "token") || strings.Contains(name, "credential") || strings.Contains(name, "cipher") {
			t.Fatalf("Project leaks secret field %s", typeOf.Field(index).Name)
		}
	}
}

func TestAuthenticationDomainEntitiesDoNotContainSecretMaterial(t *testing.T) {
	t.Parallel()
	for _, entity := range []any{Administrator{}, AdministratorSession{}, PendingEnrollment{}} {
		typeOf := reflect.TypeOf(entity)
		for index := 0; index < typeOf.NumField(); index++ {
			name := strings.ToLower(typeOf.Field(index).Name)
			if strings.Contains(name, "password") || strings.Contains(name, "secret") || strings.Contains(name, "tokenhash") || strings.Contains(name, "cipher") {
				t.Fatalf("%s leaks persistence-only secret field %s", typeOf.Name(), typeOf.Field(index).Name)
			}
		}
	}
}
