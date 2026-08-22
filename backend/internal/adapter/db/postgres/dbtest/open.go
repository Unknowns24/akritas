package dbtest

import (
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres"
	"github.com/Unknowns24/akritas/backend/internal/adapter/db/postgres/migrations"
	embeddedpostgres "github.com/fergusstrange/embedded-postgres"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func OpenMigrated(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("AKRITAS_POSTGRES_TEST_DSN")
	if dsn == "" {
		dsn = startEmbedded(t)
	}
	db, err := postgres.Open(dsn)
	if err != nil {
		t.Fatalf("open postgres: %v", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	sqlDB.SetMaxOpenConns(1)
	t.Cleanup(func() { _ = sqlDB.Close() })

	schemaName := "t_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	if err := db.Exec("CREATE SCHEMA " + schemaName).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}
	if err := db.Exec("SET search_path TO " + schemaName).Error; err != nil {
		t.Fatalf("search_path: %v", err)
	}
	t.Cleanup(func() {
		_ = db.Exec("DROP SCHEMA IF EXISTS " + schemaName + " CASCADE")
	})
	if err := migrations.Run(db); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

func startEmbedded(t *testing.T) string {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("allocate postgres port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	if err := listener.Close(); err != nil {
		t.Fatal(err)
	}
	dsn := fmt.Sprintf("postgres://akritas:akritas@127.0.0.1:%d/akritas?sslmode=disable", port)
	deadline := time.Now().Add(90 * time.Second)
	var startErr error
	for attempt := 1; time.Now().Before(deadline); attempt++ {
		database := embeddedpostgres.NewDatabase(embeddedpostgres.DefaultConfig().
			Username("akritas").
			Password("akritas").
			Database("akritas").
			Port(uint32(port)).
			RuntimePath(filepath.Join(t.TempDir(), fmt.Sprintf("pg-%d", attempt))).
			Locale("C").
			Logger(io.Discard).
			StartTimeout(60 * time.Second))
		startErr = database.Start()
		if startErr == nil {
			t.Cleanup(func() { _ = database.Stop() })
			return dsn
		}
		time.Sleep(time.Second)
	}
	t.Fatalf("start embedded postgres: %v (or set AKRITAS_POSTGRES_TEST_DSN)", startErr)
	return ""
}
