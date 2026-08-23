//go:build integration

package postgres_test

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/crypto/credentialcipher"
	dbadapter "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/migrations"
	administratorrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/credentialstore"
	githubrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubaccount"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
	portsout "github.com/Unknowns24/akritas/backend/internal/core/ports/out"
	"github.com/Unknowns24/akritas/backend/internal/core/ports/paging"
	"github.com/go-gormigrate/gormigrate/v2"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	postgrescontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	ukerpagination "github.com/unknowns24/uker/uker/pagination"
)

func TestMigrationsAndEncryptedCredentialStoreAgainstPostgreSQL(t *testing.T) {
	testcontainers.SkipIfProviderIsNotHealthy(t)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Minute)
	defer cancel()
	container, err := postgrescontainer.Run(ctx, "postgres:17-alpine",
		postgrescontainer.WithDatabase("akritas"), postgrescontainer.WithUsername("akritas"), postgrescontainer.WithPassword("test-password"),
		testcontainers.WithWaitStrategy(wait.ForLog("database system is ready to accept connections").WithOccurrence(2).WithStartupTimeout(2*time.Minute)),
	)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = testcontainers.TerminateContainer(container) }()
	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}
	db, err := dbadapter.Open(dbadapter.Config{DSN: dsn})
	if err != nil {
		t.Fatal(err)
	}
	if err := migrations.Run(db); err != nil {
		t.Fatal(err)
	}
	allTables := []string{"github_accounts", "dokploy_servers", "credentials", "github_app_registrations", "github_app_bindings", "administrators", "pending_enrollments", "administrator_sessions", "projects", "monitoring_checkpoints", "incidents", "log_events"}
	for _, table := range allTables {
		if !db.Migrator().HasTable(table) {
			t.Fatalf("migration did not create %s", table)
		}
	}

	cipher, _ := credentialcipher.NewFromKey(bytes.Repeat([]byte{7}, 32))
	credentials, _ := credentialstore.New(db, cipher)
	administrators := administratorrepo.NewRepository(db)
	transactor := dbadapter.NewTransactor(db)
	rollbackAdministrator, err := domain.NewAdministrator(uuid.New(), "rollback@example.com", "Rollback", time.Now().UTC())
	if err != nil {
		t.Fatal(err)
	}
	rollbackCause := errors.New("rollback shared transaction")
	err = transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		if createErr := administrators.Create(txCtx, rollbackAdministrator, "password-hash"); createErr != nil {
			return createErr
		}
		if putErr := credentials.Put(txCtx, portsout.CredentialOwnerAdministrator, rollbackAdministrator.ID, portsout.SecretValue{Kind: portsout.SecretKindAdministratorTOTP, Plaintext: []byte("JBSWY3DPEHPK3PXP")}); putErr != nil {
			return putErr
		}
		return rollbackCause
	})
	if !errors.Is(err, rollbackCause) {
		t.Fatalf("shared transaction error = %v", err)
	}
	if exists, existsErr := administrators.ExistsActive(ctx); existsErr != nil || exists {
		t.Fatalf("administrator survived rollback: exists=%v err=%v", exists, existsErr)
	}
	if _, getErr := credentials.Get(ctx, portsout.CredentialOwnerAdministrator, rollbackAdministrator.ID, portsout.SecretKindAdministratorTOTP); getErr == nil {
		t.Fatal("credential survived shared rollback")
	}
	repository, _ := githubrepo.New(db, credentials)
	now := time.Now().UTC()
	account, err := domain.NewGitHubAccount(uuid.New(), "Acme", domain.GitHubAccountOrganization, domain.GitHubAuthenticationPersonalAccessToken, "acme", domain.IntegrationStatusConnected, now)
	if err != nil {
		t.Fatal(err)
	}
	account.CredentialConfigured = true
	plaintext := []byte("github_pat_secret_value_123")
	if err := repository.CreateWithCredential(ctx, account, portsout.SecretValue{Kind: portsout.SecretKindGitHubPAT, Plaintext: plaintext}); err != nil {
		t.Fatal(err)
	}
	var stored struct {
		Ciphertext []byte
		Nonce      []byte
	}
	if err := db.Raw("SELECT ciphertext, nonce FROM credentials WHERE owner_id = ?", account.ID).Scan(&stored).Error; err != nil {
		t.Fatal(err)
	}
	if bytes.Contains(stored.Ciphertext, plaintext) || len(stored.Nonce) != 12 {
		t.Fatal("credential was not safely sealed")
	}
	decrypted, err := credentials.Get(ctx, portsout.CredentialOwnerGitHubAccount, account.ID, portsout.SecretKindGitHubPAT)
	if err != nil || !bytes.Equal(decrypted, plaintext) {
		t.Fatalf("credential round-trip failed: %v", err)
	}
	clear(decrypted)

	for index, createdAt := range []time.Time{now.Add(-time.Minute), now.Add(-2 * time.Minute)} {
		additional, createErr := domain.NewGitHubAccount(
			uuid.New(), "Account", domain.GitHubAccountPersonal,
			domain.GitHubAuthenticationPersonalAccessToken,
			[]string{"second", "third"}[index], domain.IntegrationStatusConnected, createdAt,
		)
		if createErr != nil {
			t.Fatal(createErr)
		}
		if createErr := db.Table("github_accounts").Create(additional).Error; createErr != nil {
			t.Fatal(createErr)
		}
	}

	pageSecret := bytes.Repeat([]byte{9}, 32)
	firstParams := paging.Params{
		Limit: 1,
		Sort: []ukerpagination.SortExpression{
			{Field: "created_at", Direction: ukerpagination.DirectionDesc},
			{Field: "id", Direction: ukerpagination.DirectionDesc},
		},
		Filters: map[string]string{"authentication_status_in": "connected"},
	}
	firstResult, err := repository.List(ctx, firstParams)
	if err != nil {
		t.Fatal(err)
	}
	if firstResult.Total != 3 || len(firstResult.Items) != 2 {
		t.Fatalf("filtered first page = %+v", firstResult)
	}
	firstPage, err := ukerpagination.BuildPageSigned(firstParams, firstResult.Items, firstParams.Limit, firstResult.Total, nil, pageSecret)
	if err != nil || firstPage.Paging.NextCursor == "" || firstPage.Paging.PrevCursor != "" {
		t.Fatalf("first page navigation = %+v, %v", firstPage.Paging, err)
	}
	nextValues := url.Values{"cursor": {firstPage.Paging.NextCursor}}
	nextParams, err := ukerpagination.ParseWithSecurity(nextValues, pageSecret, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	secondResult, err := repository.List(ctx, nextParams)
	if err != nil {
		t.Fatal(err)
	}
	secondPage, err := ukerpagination.BuildPageSigned(nextParams, secondResult.Items, nextParams.Limit, secondResult.Total, nil, pageSecret)
	if err != nil || secondPage.Paging.PrevCursor == "" {
		t.Fatalf("second page navigation = %+v, %v", secondPage.Paging, err)
	}
	prevValues := url.Values{"cursor": {secondPage.Paging.PrevCursor}}
	prevParams, err := ukerpagination.ParseWithSecurity(prevValues, pageSecret, time.Minute)
	if err != nil {
		t.Fatal(err)
	}
	previousResult, err := repository.List(ctx, prevParams)
	if err != nil {
		t.Fatal(err)
	}
	previousPage, err := ukerpagination.BuildPageSigned(prevParams, previousResult.Items, prevParams.Limit, previousResult.Total, nil, pageSecret)
	if err != nil || len(previousPage.Data) != 1 || previousPage.Data[0].ID != firstPage.Data[0].ID {
		t.Fatalf("previous navigation did not return first page: %+v, %v", previousPage, err)
	}

	migrator := gormigrate.New(db, gormigrate.DefaultOptions, migrations.All())
	for range migrations.All() {
		if err := migrator.RollbackLast(); err != nil {
			t.Fatal(err)
		}
	}
	for _, table := range allTables {
		if db.Migrator().HasTable(table) {
			t.Fatalf("rollback left table %s", table)
		}
	}
}
