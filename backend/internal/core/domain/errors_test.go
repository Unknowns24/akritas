package domain

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"
)

func TestDomainErrorContract(t *testing.T) {
	t.Parallel()

	codePattern := regexp.MustCompile(`^0x4[0-9A-F]{2}[0-9A-F]{3}[VUFCNI]$`)
	seen := make(map[string]string, len(DomainErrors()))
	for name, domainErr := range DomainErrors() {
		if !codePattern.MatchString(domainErr.Code) {
			t.Fatalf("%s has invalid code %q", name, domainErr.Code)
		}
		if previous, exists := seen[domainErr.Code]; exists {
			t.Fatalf("%s and %s share code %s", name, previous, domainErr.Code)
		}
		seen[domainErr.Code] = name
		if domainErr.Message == "" || domainErr.UserMessage == "" {
			t.Fatalf("%s must expose safe messages", name)
		}
	}
}

func TestDomainErrorCatalogMatchesDocumentation(t *testing.T) {
	t.Parallel()

	documentation, err := os.ReadFile("../../../docs/errors/aaa-map.md")
	if err != nil {
		t.Fatalf("read error catalog: %v", err)
	}
	contents := string(documentation)
	for name, domainErr := range DomainErrors() {
		if !strings.Contains(contents, "`"+name+"`") || !strings.Contains(contents, "`"+domainErr.Code+"`") {
			t.Fatalf("catalog does not contain %s (%s)", name, domainErr.Code)
		}
	}
}

func TestDomainErrorWrapPreservesIdentityAndCause(t *testing.T) {
	t.Parallel()

	cause := errors.New("provider detail")
	wrapped := ErrInvalidProject.Wrap(cause)
	if !errors.Is(wrapped, ErrInvalidProject) {
		t.Fatal("wrapped error must match its sentinel")
	}
	if !errors.Is(wrapped, cause) {
		t.Fatal("wrapped error must unwrap to its cause")
	}
	if got := fmt.Sprint(wrapped); got != ErrInvalidProject.Message {
		t.Fatalf("Error() leaked cause: %q", got)
	}
}
