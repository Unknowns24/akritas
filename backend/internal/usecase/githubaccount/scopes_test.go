package githubaccount

import (
	"testing"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func TestClassicAndFineGrainedPATPermissionPolicy(t *testing.T) {
	for _, test := range []struct {
		name        string
		accountType domain.GitHubAccountType
		scopes      []string
		want        bool
	}{
		{name: "fine grained has no introspection header", accountType: domain.GitHubAccountPersonal, scopes: nil, want: true},
		{name: "classic personal requires repo", accountType: domain.GitHubAccountPersonal, scopes: []string{"repo"}, want: true},
		{name: "classic public repo is insufficient", accountType: domain.GitHubAccountPersonal, scopes: []string{"public_repo"}, want: false},
		{name: "classic organization requires membership visibility", accountType: domain.GitHubAccountOrganization, scopes: []string{"repo", "read:org"}, want: true},
		{name: "classic organization missing read org", accountType: domain.GitHubAccountOrganization, scopes: []string{"repo"}, want: false},
	} {
		t.Run(test.name, func(t *testing.T) {
			if got := classicScopesSatisfy(test.accountType, test.scopes); got != test.want {
				t.Fatalf("classicScopesSatisfy() = %v, want %v", got, test.want)
			}
		})
	}
}
