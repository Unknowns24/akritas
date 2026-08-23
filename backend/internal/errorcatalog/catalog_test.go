package errorcatalog_test

import (
	"os"
	"regexp"
	"strings"
	"testing"

	dberrors "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/errors"
	githubexternal "github.com/Unknowns24/akritas/backend/internal/adapter/external/github"
	resterrors "github.com/Unknowns24/akritas/backend/internal/adapter/rest/errors"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func TestStableErrorCatalogsAreGloballyUniqueAndDocumentedOnce(t *testing.T) {
	t.Parallel()
	documentation, err := os.ReadFile("../../docs/errors/aaa-map.md")
	if err != nil {
		t.Fatal(err)
	}
	catalogs := []map[string]*domain.Error{
		domain.DomainErrors(), domain.AuthenticationErrors(), domain.IntegrationErrors(), domain.ProjectErrors(),
		domain.IncidentErrors(), domain.MonitoringErrors(), domain.InvestigationErrors(), domain.OperationErrors(),
		dberrors.Catalog(), resterrors.Catalog(), githubexternal.Catalog(),
	}
	pattern := regexp.MustCompile(`^[012F]x[0-9A-F]{6}[VUFCNIR]$`)
	seenCode := make(map[string]string)
	seenName := make(map[string]struct{})
	contents := string(documentation)
	for _, catalog := range catalogs {
		for name, stable := range catalog {
			if stable == nil || !pattern.MatchString(stable.Code) || stable.Message == "" || stable.UserMessage == "" {
				t.Fatalf("invalid stable error %s: %#v", name, stable)
			}
			if previous := seenCode[stable.Code]; previous != "" {
				t.Fatalf("duplicate exact stable code %s: %s and %s", stable.Code, previous, name)
			}
			if _, duplicate := seenName[name]; duplicate {
				t.Fatalf("sentinel cataloged more than once: %s", name)
			}
			seenCode[stable.Code] = name
			seenName[name] = struct{}{}
			if strings.Count(contents, "`"+name+"`") != 1 || strings.Count(contents, "`"+stable.Code+"`") != 1 {
				t.Fatalf("documentation must contain %s (%s) exactly once", name, stable.Code)
			}
		}
	}
}
