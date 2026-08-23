//go:build integration

package postgres_test

import (
	"bytes"
	"context"
	"errors"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/crypto/credentialcipher"
	dbadapter "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/migrations"
	administratorrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator"
	administratorsessionrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator_session"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/credentialstore"
	githubrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/githubaccount"
	pendingenrollmentrepo "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/pending_enrollment"
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
	allTables := []string{"github_accounts", "dokploy_servers", "credentials", "github_app_registrations", "github_app_bindings", "administrators", "pending_enrollments", "administrator_sessions", "projects", "monitoring_checkpoints", "incidents", "log_events", "investigations", "operations", "evidence"}
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

	// Recovery rotates metadata, encrypted TOTP ownership and all sessions as
	// one PostgreSQL invariant. First prove rollback, then commit the same shape.
	recoveryAdministrator, err := domain.NewAdministrator(uuid.New(), "recovery@example.com", "Recovery", time.Now().UTC().Add(-time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if err := administrators.Create(ctx, recoveryAdministrator, "old-password-hash"); err != nil {
		t.Fatal(err)
	}
	if err := credentials.Put(ctx, portsout.CredentialOwnerAdministrator, recoveryAdministrator.ID, portsout.SecretValue{Kind: portsout.SecretKindAdministratorTOTP, Plaintext: []byte("OLDTOTPSEEDVALUE")}); err != nil {
		t.Fatal(err)
	}
	sessions := administratorsessionrepo.NewRepository(db)
	for index := 0; index < 2; index++ {
		authenticatedAt := time.Now().UTC()
		session, sessionErr := domain.NewAdministratorSession(uuid.New(), recoveryAdministrator.ID, authenticatedAt, authenticatedAt.Add(time.Hour), authenticatedAt.Add(2*time.Hour))
		if sessionErr != nil {
			t.Fatal(sessionErr)
		}
		if saveErr := sessions.Save(ctx, session, strings.Repeat(string(rune('a'+index)), 64)); saveErr != nil {
			t.Fatal(saveErr)
		}
	}
	pendingEnrollments := pendingenrollmentrepo.NewRepository(db)
	pending, err := domain.NewPendingEnrollment(uuid.New(), recoveryAdministrator.Email, recoveryAdministrator.DisplayName, time.Now().UTC(), time.Now().UTC().Add(10*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pendingEnrollments.Replace(ctx, pending, "new-password-hash"); err != nil {
		t.Fatal(err)
	}
	if err := credentials.Put(ctx, portsout.CredentialOwnerPendingEnrollment, pending.ID, portsout.SecretValue{Kind: portsout.SecretKindAdministratorTOTP, Plaintext: []byte("NEWTOTPSEEDVALUE")}); err != nil {
		t.Fatal(err)
	}

	err = transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		consumed, consumeErr := pendingEnrollments.Consume(txCtx, pending.ID)
		if consumeErr != nil || consumed == nil {
			return consumeErr
		}
		rotated, rotateErr := administrators.RotateCredentials(txCtx, recoveryAdministrator.ID, "old-password-hash", consumed.PasswordHash, 123, time.Now().UTC())
		if rotateErr != nil || !rotated {
			return rotateErr
		}
		if moveErr := credentials.MoveOwner(txCtx, portsout.CredentialOwnerPendingEnrollment, pending.ID, portsout.CredentialOwnerAdministrator, recoveryAdministrator.ID); moveErr != nil {
			return moveErr
		}
		if revokeErr := sessions.RevokeAll(txCtx, recoveryAdministrator.ID, time.Now().UTC()); revokeErr != nil {
			return revokeErr
		}
		return rollbackCause
	})
	if !errors.Is(err, rollbackCause) {
		t.Fatalf("recovery rollback error=%v", err)
	}
	if current, findErr := administrators.FindByEmail(ctx, recoveryAdministrator.Email); findErr != nil || current == nil || current.PasswordHash != "old-password-hash" {
		t.Fatalf("password survived rollback incorrectly: current=%+v err=%v", current, findErr)
	}
	if currentPending, findErr := pendingEnrollments.FindByID(ctx, pending.ID); findErr != nil || currentPending == nil {
		t.Fatalf("pending enrollment not restored: pending=%+v err=%v", currentPending, findErr)
	}
	var revokedCount int64
	if countErr := db.Table("administrator_sessions").Where("administrator_id = ? AND revoked_at IS NOT NULL", recoveryAdministrator.ID).Count(&revokedCount).Error; countErr != nil || revokedCount != 0 {
		t.Fatalf("rollback revoked sessions: count=%d err=%v", revokedCount, countErr)
	}

	freshAt := time.Now().UTC()
	freshSession, err := domain.NewAdministratorSession(uuid.New(), recoveryAdministrator.ID, freshAt, freshAt.Add(time.Hour), freshAt.Add(2*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	err = transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
		consumed, consumeErr := pendingEnrollments.Consume(txCtx, pending.ID)
		if consumeErr != nil || consumed == nil {
			return consumeErr
		}
		rotated, rotateErr := administrators.RotateCredentials(txCtx, recoveryAdministrator.ID, "old-password-hash", consumed.PasswordHash, 456, freshAt)
		if rotateErr != nil || !rotated {
			return rotateErr
		}
		if moveErr := credentials.MoveOwner(txCtx, portsout.CredentialOwnerPendingEnrollment, pending.ID, portsout.CredentialOwnerAdministrator, recoveryAdministrator.ID); moveErr != nil {
			return moveErr
		}
		if revokeErr := sessions.RevokeAll(txCtx, recoveryAdministrator.ID, freshAt); revokeErr != nil {
			return revokeErr
		}
		return sessions.Save(txCtx, freshSession, strings.Repeat("f", 64))
	})
	if err != nil {
		t.Fatalf("commit recovery: %v", err)
	}
	current, err := administrators.FindByEmail(ctx, recoveryAdministrator.Email)
	if err != nil || current == nil || current.PasswordHash != "new-password-hash" || current.Administrator.LastAcceptedTOTPPeriod != 456 {
		t.Fatalf("committed credentials=%+v err=%v", current, err)
	}
	rotatedSecret, err := credentials.Get(ctx, portsout.CredentialOwnerAdministrator, recoveryAdministrator.ID, portsout.SecretKindAdministratorTOTP)
	if err != nil || string(rotatedSecret) != "NEWTOTPSEEDVALUE" {
		t.Fatalf("committed TOTP rotation mismatch: value=%q err=%v", rotatedSecret, err)
	}
	clear(rotatedSecret)
	if currentPending, findErr := pendingEnrollments.FindByID(ctx, pending.ID); findErr != nil || currentPending != nil {
		t.Fatalf("pending enrollment survived commit: pending=%+v err=%v", currentPending, findErr)
	}
	if countErr := db.Table("administrator_sessions").Where("administrator_id = ? AND revoked_at IS NOT NULL", recoveryAdministrator.ID).Count(&revokedCount).Error; countErr != nil || revokedCount != 2 {
		t.Fatalf("old sessions not revoked: count=%d err=%v", revokedCount, countErr)
	}
	if fresh, findErr := sessions.FindByTokenHash(ctx, strings.Repeat("f", 64)); findErr != nil || fresh == nil || !fresh.IsActive(freshAt) {
		t.Fatalf("fresh session invalid: session=%+v err=%v", fresh, findErr)
	}
	if consumed, consumeErr := administrators.ConsumeTOTPPeriod(ctx, recoveryAdministrator.ID, "old-password-hash", 999); consumeErr != nil || consumed {
		t.Fatalf("stale login credentials advanced TOTP state: consumed=%v err=%v", consumed, consumeErr)
	}

	// DELETE ... RETURNING is the linearization point for confirmation. Two
	// simultaneous requests may race, but exactly one can consume the slot.
	concurrentPending, err := domain.NewPendingEnrollment(uuid.New(), recoveryAdministrator.Email, recoveryAdministrator.DisplayName, freshAt, freshAt.Add(10*time.Minute))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := pendingEnrollments.Replace(ctx, concurrentPending, "another-password-hash"); err != nil {
		t.Fatal(err)
	}
	consumeResults := make(chan bool, 2)
	consumeErrors := make(chan error, 2)
	for range 2 {
		go func() {
			consumed, consumeErr := pendingEnrollments.Consume(ctx, concurrentPending.ID)
			consumeErrors <- consumeErr
			consumeResults <- consumed != nil
		}()
	}
	consumeSuccesses := 0
	for range 2 {
		if consumeErr := <-consumeErrors; consumeErr != nil {
			t.Fatal(consumeErr)
		}
		if <-consumeResults {
			consumeSuccesses++
		}
	}
	if consumeSuccesses != 1 {
		t.Fatalf("concurrent pending consumers succeeded=%d want 1", consumeSuccesses)
	}

	// Hold the refresh transaction open after it has locked the row. Revoke-all
	// must wait, then revoke the refreshed session once that transaction commits.
	refreshStarted := make(chan struct{})
	releaseRefresh := make(chan struct{})
	refreshDone := make(chan error, 1)
	go func() {
		refreshDone <- transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
			refreshed, refreshErr := sessions.RefreshActive(txCtx, strings.Repeat("f", 64), freshAt.Add(time.Minute), freshAt.Add(90*time.Minute))
			if refreshErr != nil {
				return refreshErr
			}
			if refreshed == nil {
				return errors.New("fresh session was not refreshable")
			}
			close(refreshStarted)
			<-releaseRefresh
			return nil
		})
	}()
	<-refreshStarted
	revokeDone := make(chan error, 1)
	go func() { revokeDone <- sessions.RevokeAll(ctx, recoveryAdministrator.ID, freshAt.Add(2*time.Minute)) }()
	close(releaseRefresh)
	if err := <-refreshDone; err != nil {
		t.Fatal(err)
	}
	if err := <-revokeDone; err != nil {
		t.Fatal(err)
	}
	if refreshed, refreshErr := sessions.RefreshActive(ctx, strings.Repeat("f", 64), freshAt.Add(3*time.Minute), freshAt.Add(2*time.Hour)); refreshErr != nil || refreshed != nil {
		t.Fatalf("revoke did not win after in-flight refresh: session=%+v err=%v", refreshed, refreshErr)
	}

	// Reverse the ordering: revoke holds the row lock before refresh starts, so
	// the conditional UPDATE observes revoked_at after the lock is released.
	revokedFirstAt := freshAt.Add(4 * time.Minute)
	revokedFirstSession, err := domain.NewAdministratorSession(uuid.New(), recoveryAdministrator.ID, freshAt, freshAt.Add(time.Hour), freshAt.Add(2*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	revokedFirstHash := strings.Repeat("r", 64)
	if err := sessions.Save(ctx, revokedFirstSession, revokedFirstHash); err != nil {
		t.Fatal(err)
	}
	revokeStarted := make(chan struct{})
	releaseRevoke := make(chan struct{})
	revokeTransactionDone := make(chan error, 1)
	go func() {
		revokeTransactionDone <- transactor.WithinTransaction(ctx, func(txCtx context.Context) error {
			if revokeErr := sessions.RevokeAll(txCtx, recoveryAdministrator.ID, revokedFirstAt); revokeErr != nil {
				return revokeErr
			}
			close(revokeStarted)
			<-releaseRevoke
			return nil
		})
	}()
	<-revokeStarted
	blockedRefresh := make(chan *domain.AdministratorSession, 1)
	blockedRefreshErr := make(chan error, 1)
	go func() {
		refreshed, refreshErr := sessions.RefreshActive(ctx, revokedFirstHash, revokedFirstAt, revokedFirstAt.Add(time.Hour))
		blockedRefresh <- refreshed
		blockedRefreshErr <- refreshErr
	}()
	close(releaseRevoke)
	if err := <-revokeTransactionDone; err != nil {
		t.Fatal(err)
	}
	if refreshErr := <-blockedRefreshErr; refreshErr != nil {
		t.Fatal(refreshErr)
	}
	if refreshed := <-blockedRefresh; refreshed != nil {
		t.Fatalf("refresh survived revoke-first ordering: %+v", refreshed)
	}
	repository, _ := githubrepo.New(db, credentials)
	// Cursor timestamps are serialized at second precision. Keep the fixture on
	// that same boundary so previous-page assertions are deterministic.
	now := time.Now().UTC().Truncate(time.Second)
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
