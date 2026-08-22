package common_test

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/rest/mapper"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	"github.com/google/uuid"
)

func TestSafeIntegrationDTOsCannotSerializeCredentialMaterial(t *testing.T) {
	now := time.Now().UTC()
	account, err := domain.NewGitHubAccount(uuid.New(), "Acme", domain.GitHubAccountOrganization, domain.GitHubAuthenticationPersonalAccessToken, "acme", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	account.CredentialConfigured = true
	server, err := domain.NewDokployServer(uuid.New(), "Production", "https://dokploy.example.com", strings.Repeat("a", 64), domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	server.CredentialConfigured = true
	encoded, err := json.Marshal(map[string]any{"github": mapper.GitHubAccountToDTO(*account), "dokploy": mapper.DokployServerToDTO(*server)})
	if err != nil {
		t.Fatal(err)
	}
	serialized := string(encoded)
	for _, forbidden := range []string{"personal_access_token", "api_credential", "private_key", "webhook_secret", "ciphertext", "nonce", "installation_token"} {
		if strings.Contains(serialized, `"`+forbidden+`":`) {
			t.Fatalf("safe DTO leaked forbidden property %q: %s", forbidden, serialized)
		}
	}
	if !strings.Contains(serialized, "credential_configured") {
		t.Fatal("safe credential projection is missing")
	}
}
