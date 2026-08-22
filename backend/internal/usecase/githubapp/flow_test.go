package githubapp

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsin "github.com/Unknowns24/akritas/backend/internal/core/ports/in"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/google/uuid"
)

func TestManifestFlowUsesSeparateOneTimeStatesAndVerifiedInstallation(t *testing.T) {
	now := time.Date(2026, 8, 22, 12, 0, 0, 0, time.UTC)
	store := &registrationStoreFake{}
	gateway := &appGatewayFake{conversion: portsout.GitHubManifestConversion{
		AppID: 42, AppSlug: "akritas-test", AppName: "Akritas", ClientID: "Iv1.client", PrivateKey: []byte("private-key"), WebhookSecret: []byte("webhook-secret"),
	}, installation: portsout.GitHubInstallation{InstallationID: 99, AccountLogin: "acme", AccountType: domain.GitHubAccountOrganization}}
	uc, err := New(store, gateway, "https://akritas.example.com", uuid.New, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	states := []string{"conversion-state-value-that-is-long-enough-0001", "installation-state-value-that-is-long-enough-02"}
	uc.newState = func() (string, error) {
		state := states[0]
		states = states[1:]
		return state, nil
	}

	started, err := uc.StartRegistration(context.Background(), portsin.StartGitHubAppRegistrationCommand{DisplayName: "Acme", OwnerType: domain.GitHubAccountOrganization, Organization: "acme"})
	if err != nil {
		t.Fatal(err)
	}
	if string(store.registration.ConversionStateDigest) == started.State {
		t.Fatal("raw state was persisted")
	}
	var manifest map[string]any
	if err := json.Unmarshal([]byte(started.Manifest), &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest["public"] != false || manifest["hook_attributes"].(map[string]any)["active"] != false {
		t.Fatalf("manifest is not private with disabled webhooks: %#v", manifest)
	}

	converted, err := uc.CompleteManifest(context.Background(), "manifest-code-at-least-twenty", started.State)
	if err != nil {
		t.Fatal(err)
	}
	if converted.RedirectURL == "" || store.registration.Status != portsout.GitHubAppRegistrationConverted {
		t.Fatal("conversion was not completed")
	}
	if len(store.secrets) != 2 || string(store.secrets[0].Plaintext) == "" {
		t.Fatal("App secrets were not handed to the credential transaction")
	}
	if _, err := uc.CompleteManifest(context.Background(), "manifest-code-at-least-twenty", started.State); !errors.Is(err, domain.ErrManifestStateConflict) {
		t.Fatalf("expected replay conflict, got %v", err)
	}

	installed, err := uc.CompleteInstallation(context.Background(), 99, "installation-state-value-that-is-long-enough-02")
	if err != nil {
		t.Fatal(err)
	}
	if installed.Account.AccountIdentifier != "acme" || installed.Account.AuthenticationMethod != domain.GitHubAuthenticationGitHubApp || store.binding.InstallationID != 99 {
		t.Fatalf("unexpected installed account: %#v", installed.Account)
	}
}

func TestInstallationRejectsUnverifiedOwnerWithoutCreatingAccount(t *testing.T) {
	now := time.Now().UTC()
	state := "installation-state-value-that-is-long-enough-99"
	digest := sha256.Sum256([]byte(state))
	appID := int64(42)
	store := &registrationStoreFake{registration: &portsout.GitHubAppRegistration{
		ID: uuid.New(), DisplayName: "Acme", AccountType: domain.GitHubAccountOrganization, AccountIdentifier: "acme",
		InstallationStateDigest: digest[:], Status: portsout.GitHubAppRegistrationConverted, AppID: &appID, ExpiresAt: now.Add(time.Hour),
	}}
	gateway := &appGatewayFake{installation: portsout.GitHubInstallation{InstallationID: 99, AccountLogin: "attacker", AccountType: domain.GitHubAccountOrganization}}
	uc, _ := New(store, gateway, "https://akritas.example.com", uuid.New, func() time.Time { return now })
	_, err := uc.CompleteInstallation(context.Background(), 99, state)
	if !errors.Is(err, domain.ErrGitHubCredentialRejected) || store.completedAccount != nil {
		t.Fatalf("unverified installation was not rejected atomically: %v", err)
	}
}

type registrationStoreFake struct {
	registration     *portsout.GitHubAppRegistration
	secrets          []portsout.SecretValue
	binding          portsout.GitHubAppBinding
	completedAccount *domain.GitHubAccount
	conversionUsed   bool
	installationUsed bool
}

func (s *registrationStoreFake) CreateRegistration(_ context.Context, registration *portsout.GitHubAppRegistration) error {
	copy := *registration
	s.registration = &copy
	return nil
}
func (s *registrationStoreFake) ConsumeConversionState(_ context.Context, digest []byte, now time.Time) (*portsout.GitHubAppRegistration, error) {
	if s.registration == nil || s.conversionUsed || now.After(s.registration.ExpiresAt) || !equalBytes(digest, s.registration.ConversionStateDigest) {
		return nil, domain.ErrManifestStateConflict
	}
	s.conversionUsed = true
	return s.registration, nil
}
func (s *registrationStoreFake) CompleteConversion(_ context.Context, registration *portsout.GitHubAppRegistration, secrets []portsout.SecretValue) error {
	s.registration = registration
	s.secrets = secrets
	return nil
}
func (s *registrationStoreFake) ConsumeInstallationState(_ context.Context, digest []byte, now time.Time) (*portsout.GitHubAppRegistration, error) {
	if s.registration == nil || s.installationUsed || now.After(s.registration.ExpiresAt) || !equalBytes(digest, s.registration.InstallationStateDigest) {
		return nil, domain.ErrManifestStateConflict
	}
	s.installationUsed = true
	return s.registration, nil
}
func (s *registrationStoreFake) CompleteInstallation(_ context.Context, registration *portsout.GitHubAppRegistration, account *domain.GitHubAccount, binding portsout.GitHubAppBinding) error {
	s.registration, s.completedAccount, s.binding = registration, account, binding
	return nil
}

type appGatewayFake struct {
	conversion   portsout.GitHubManifestConversion
	installation portsout.GitHubInstallation
}

func (g *appGatewayFake) ExchangeManifest(context.Context, string) (portsout.GitHubManifestConversion, error) {
	return g.conversion, nil
}
func (g *appGatewayFake) VerifyInstallation(context.Context, portsout.GitHubAppRegistration, int64) (portsout.GitHubInstallation, error) {
	return g.installation, nil
}

func equalBytes(left, right []byte) bool {
	if len(left) != len(right) {
		return false
	}
	for i := range left {
		if left[i] != right[i] {
			return false
		}
	}
	return true
}
