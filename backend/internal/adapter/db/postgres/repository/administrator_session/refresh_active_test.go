package administratorsession_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	administratorsession "github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator_session"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func TestRefreshActiveEnforcesExpiryAndRevocationAtomically(t *testing.T) {
	db := dbtest.Connect(t)
	repo := administratorsession.NewRepository(db)
	administratorID := seedAdministrator(t, db)
	now := time.Now().UTC()
	active, _ := domain.NewAdministratorSession(uuid.New(), administratorID, now.Add(-time.Minute), now.Add(time.Minute), now.Add(time.Hour))
	if err := repo.Save(context.Background(), active, "active-hash"); err != nil {
		t.Fatal(err)
	}
	refreshed, err := repo.RefreshActive(context.Background(), "active-hash", now, now.Add(2*time.Hour))
	if err != nil || refreshed == nil {
		t.Fatalf("refresh: session=%v err=%v", refreshed, err)
	}
	if !refreshed.IdleExpiresAt.Equal(refreshed.AbsoluteExpiresAt) {
		t.Fatalf("idle not capped: %+v", refreshed)
	}
	if err := repo.Revoke(context.Background(), active.ID, now); err != nil {
		t.Fatal(err)
	}
	if found, err := repo.RefreshActive(context.Background(), "active-hash", now, now.Add(time.Minute)); err != nil || found != nil {
		t.Fatalf("revoked refresh: session=%v err=%v", found, err)
	}
	if found, err := repo.RefreshActive(context.Background(), "random", now, now.Add(time.Minute)); err != nil || found != nil {
		t.Fatalf("random refresh: session=%v err=%v", found, err)
	}
}

func TestRevokeAllRevokesEveryExistingSession(t *testing.T) {
	db := dbtest.Connect(t)
	repo := administratorsession.NewRepository(db)
	administratorID := seedAdministrator(t, db)
	now := time.Now().UTC()
	for i := 0; i < 2; i++ {
		session, _ := domain.NewAdministratorSession(uuid.New(), administratorID, now, now.Add(time.Hour), now.Add(2*time.Hour))
		if err := repo.Save(context.Background(), session, uuid.NewString()); err != nil {
			t.Fatal(err)
		}
	}
	if err := repo.RevokeAll(context.Background(), administratorID, now.Add(time.Minute)); err != nil {
		t.Fatal(err)
	}
	var active int64
	if err := db.Table("administrator_sessions").Where("administrator_id = ? AND revoked_at IS NULL", administratorID).Count(&active).Error; err != nil {
		t.Fatal(err)
	}
	if active != 0 {
		t.Fatalf("active sessions=%d want 0", active)
	}
}
