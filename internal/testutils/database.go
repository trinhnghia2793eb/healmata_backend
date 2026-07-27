package testutils

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/pressly/goose/v3"
)

var (
	dbSetupOnce sync.Once
	globalPool  *pgxpool.Pool
)

func loadDotEnv() {
	dir, _ := os.Getwd()
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			_ = godotenv.Load(envPath)
			return
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
}

func migrationsDir() string {
	// runtime.Caller(0) gives the path of THIS file at compile time.
	_, filename, _, _ := runtime.Caller(0)
	// Go up from testhelper/ (which is in internal/app/db/testhelper/) to internal/app/db/migrations/
	root := filepath.Join(filepath.Dir(filename), "..", "db\\migrations")
	abs, err := filepath.Abs(root)
	if err != nil {
		panic("testhelper: cannot resolve migrations dir: " + err.Error())
	}
	return abs
}

func getEnv(key, fallback string) string {
	loadDotEnv()
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

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

func SetupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	// 1. Chỉ chạy Migration và tạo Pool đúng 1 lần duy nhất
	dbSetupOnce.Do(func() {
		connStr := dsn()

		sqlDB, err := sql.Open("pgx", connStr)
		if err != nil {
			panic(fmt.Sprintf("testhelper: open sql.DB for goose: %v", err))
		}

		goose.SetDialect("postgres")
		goose.SetLogger(goose.NopLogger())
		migrDir := migrationsDir()

		if err := goose.Up(sqlDB, migrDir); err != nil {
			_ = sqlDB.Close()
			panic(fmt.Sprintf("testhelper: goose.Up failed: %v", err))
		}
		// Đóng sqlDB vì chúng ta sẽ dùng pgxpool cho các tác vụ sau
		_ = sqlDB.Close()

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		pool, err := pgxpool.New(ctx, connStr)
		if err != nil {
			panic(fmt.Sprintf("testhelper: create pgxpool: %v", err))
		}

		globalPool = pool
	})

	// 2. Dọn dẹp dữ liệu của TẤT CẢ các bảng trước khi chạy test case này
	TruncateTables(t, globalPool)

	// Không dùng goose.Reset() trong Cleanup nữa, để tái sử dụng schema
	return globalPool
}

// TruncateTables xóa toàn bộ dữ liệu nhưng giữ nguyên cấu trúc bảng
func TruncateTables(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Liệt kê các bảng cần dọn dẹp (sử dụng CASCADE để tự động clear các bảng có khóa ngoại)
	query := `TRUNCATE TABLE users, social_accounts, refresh_tokens, otp_requests, user_sessions CASCADE;`

	_, err := pool.Exec(ctx, query)
	if err != nil {
		t.Fatalf("testhelper: failed to truncate tables: %v", err)
	}
}
