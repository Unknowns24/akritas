package apperr

import (
	"errors"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestApplicationErrorContract(t *testing.T) {
	t.Parallel()

	codePattern := regexp.MustCompile(`^0x5[0-9A-F]{2}[0-9A-F]{3}[VUFCNI]$`)
	seen := make(map[string]string, len(Catalog()))
	for name, appErr := range Catalog() {
		if !codePattern.MatchString(appErr.Code) {
			t.Fatalf("%s has invalid code %q", name, appErr.Code)
		}
		if previous, exists := seen[appErr.Code]; exists {
			t.Fatalf("%s and %s share code %s", name, previous, appErr.Code)
		}
		seen[appErr.Code] = name
		if appErr.Message == "" || appErr.UserMessage == "" {
			t.Fatalf("%s must expose safe messages", name)
		}
	}
}

func TestApplicationErrorCatalogMatchesDocumentation(t *testing.T) {
	t.Parallel()

	documentation, err := os.ReadFile("../../../docs/errors/aaa-map.md")
	if err != nil {
		t.Fatalf("read error catalog: %v", err)
	}
	contents := string(documentation)
	for name, appErr := range Catalog() {
		if !strings.Contains(contents, "`"+name+"`") || !strings.Contains(contents, "`"+appErr.Code+"`") {
			t.Fatalf("catalog does not contain %s (%s)", name, appErr.Code)
		}
	}
}

func TestNotFoundAndConflictTypes(t *testing.T) {
	t.Parallel()

	if ErrProjectNotFound.Type() != 'N' {
		t.Fatalf("project not found must be type N, got %c", ErrProjectNotFound.Type())
	}
	if ErrProjectNameConflict.Type() != 'C' {
		t.Fatalf("name conflict must be type C, got %c", ErrProjectNameConflict.Type())
	}
	if !errors.Is(ErrGitHubAccountNotFound.Wrap(errors.New("missing")), ErrGitHubAccountNotFound) {
		t.Fatal("wrapped application error must match its sentinel")
	}
}
