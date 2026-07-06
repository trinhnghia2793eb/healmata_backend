// Package testhelper provides shared utilities for database integration tests.
// It sets up a real PostgreSQL connection and applies Goose migrations,
// then tears everything down after the test completes.
package testhelper

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver for database/sql (used by Goose)
	"github.com/pressly/goose/v3"
)

// migrationsDir returns the absolute path to the migrations folder,
// calculated relative to this source file so it works regardless of
// which directory `go test` is invoked from.
func migrationsDir() string {
	// runtime.Caller(0) gives the path of THIS file at compile time.
	_, filename, _, _ := runtime.Caller(0)
	// Go up from testhelper/ → db/ → internal/ → project root, then into migrations.
	root := filepath.Join(filepath.Dir(filename), "..", "migrations")
	abs, err := filepath.Abs(root)
	if err != nil {
		panic("testhelper: cannot resolve migrations dir: " + err.Error())
	}
	return abs
}

// getEnv returns the env var value or a fallback default.
func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// dsn builds a PostgreSQL DSN from environment variables.
// Reads TEST_DB_* first; falls back to DB_* so the same .env works
// for both app and tests.
func dsn() string {
	host := getEnv("TEST_DB_HOST", getEnv("DB_HOST", "localhost"))
	port := getEnv("TEST_DB_PORT", getEnv("DB_PORT", "5432"))
	user := getEnv("TEST_DB_USER", getEnv("DB_USER", ""))
	password := getEnv("TEST_DB_PASSWORD", getEnv("DB_PASSWORD", ""))
	name := getEnv("TEST_DB_NAME", getEnv("DB_NAME", ""))
	sslmode := getEnv("TEST_DB_SSLMODE", getEnv("DB_SSLMODE", "disable"))

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, name, sslmode,
	)
}

// SetupTestDB connects to PostgreSQL, runs all Goose migrations (Up),
// and registers a t.Cleanup handler that runs Goose Reset then closes
// the pool. Returns a *pgxpool.Pool ready for use.
//
// Usage:
//
//	pool := testhelper.SetupTestDB(t)
func SetupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	connStr := dsn()

	// --- 1. Run migrations via database/sql + goose ---
	// Goose requires a *sql.DB (not pgxpool). We open a separate
	// sql.DB just for migration management.
	sqlDB, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("testhelper: open sql.DB for goose: %v", err)
	}

	goose.SetDialect("postgres") //nolint:errcheck
	migrDir := migrationsDir()

	if err := goose.Up(sqlDB, migrDir); err != nil {
		_ = sqlDB.Close()
		t.Fatalf("testhelper: goose.Up failed: %v", err)
	}

	// --- 2. Connect pgxpool for the actual test queries ---
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, "postgresql://"+connStr)
	if err != nil {
		// Fallback: pgxpool also accepts keyword=value DSN via ParseConfig
		cfg, cfgErr := pgxpool.ParseConfig(connStr)
		if cfgErr != nil {
			_ = sqlDB.Close()
			t.Fatalf("testhelper: parse pgxpool config: %v", cfgErr)
		}
		pool, err = pgxpool.NewWithConfig(ctx, cfg)
		if err != nil {
			_ = sqlDB.Close()
			t.Fatalf("testhelper: create pgxpool: %v", err)
		}
	}

	// --- 3. Register cleanup: reset migrations + close connections ---
	t.Cleanup(func() {
		if err := goose.Reset(sqlDB, migrDir); err != nil {
			t.Logf("testhelper cleanup: goose.Reset error (non-fatal): %v", err)
		}
		pool.Close()
		_ = sqlDB.Close()
	})

	return pool
}
