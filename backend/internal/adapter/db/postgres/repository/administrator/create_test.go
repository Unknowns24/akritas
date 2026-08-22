package administrator_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/dbtest"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/model"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/repository/administrator"
	"github.com/Unknowns24/akritas/backend/internal/core/domain"
)

func newDomainAdministrator(t *testing.T, email string) *domain.Administrator {
	t.Helper()
	admin, err := domain.NewAdministrator(uuid.New(), email, "Akritas Administrator", time.Now().UTC())
	if err != nil {
		t.Fatalf("build administrator: %v", err)
	}
	return admin
}

func TestCreatePersistsAdministrator(t *testing.T) {
	db := dbtest.Connect(t)
	repo := administrator.NewRepository(db)

	admin := newDomainAdministrator(t, "admin@example.com")
	passwordHash := "$argon2id$v=19$m=19456,t=2,p=1$salt$hash"
	encryptedSecret := []byte("ciphertext")

	if err := repo.Create(context.Background(), admin, passwordHash, encryptedSecret); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var record model.Administrator
	if err := db.First(&record, "id = ?", admin.ID).Error; err != nil {
		t.Fatalf("query administrator: %v", err)
	}
	if record.Email != admin.Email || record.PasswordHash != passwordHash || string(record.EncryptedTOTPSecret) != string(encryptedSecret) {
		t.Fatalf("persisted row does not match input: %+v", record)
	}
}

func TestCreateRejectsDuplicateEmail(t *testing.T) {
	db := dbtest.Connect(t)
	repo := administrator.NewRepository(db)

	first := newDomainAdministrator(t, "admin@example.com")
	if err := repo.Create(context.Background(), first, "hash-one", []byte("cipher-one")); err != nil {
		t.Fatalf("unexpected error on first create: %v", err)
	}

	second := newDomainAdministrator(t, "admin@example.com")
	err := repo.Create(context.Background(), second, "hash-two", []byte("cipher-two"))
	if !errors.Is(err, domain.ErrAdministratorAlreadyExists) {
		t.Fatalf("expected ErrAdministratorAlreadyExists, got %v", err)
	}

	var count int64
	if err := db.Model(&model.Administrator{}).Count(&count).Error; err != nil {
		t.Fatalf("count administrators: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 administrator row, got %d", count)
	}
}
