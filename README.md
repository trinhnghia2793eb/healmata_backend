# HEALMATA

Backend service built with Go and Gin.

---

## Prerequisites

* Go (latest stable version recommended)
* PostgreSQL

---

## Getting Started

### 1. Clone the repository

```bash
git clone <repository-url>
cd <project-folder>
```

### 2. Create the environment file

Copy the example configuration:

```bash
cp .env.example .env
```

Then update the database configuration inside `.env`

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSLMODE=disable
```

---

## Environment

### Development

```env
APP_ENV=development
GIN_MODE=debug
```

### Production

```env
APP_ENV=production
GIN_MODE=release
```

---

## Run

Start the server:

```bash
go run ./cmd/server
```

---

## Health Check

After the server starts successfully, open:

```
http://localhost:8080/auth/health
```

---

## Testing

### Database Constraint Integration Tests

These tests verify that all PostgreSQL constraints (unique indexes, foreign keys, and cascade deletes) in the auth schema are working correctly. They run against a **real database** — no mocks.

**Prerequisites:** 

Before running the tests, ensure you have the following installed on your system:
- **Go** (version 1.20+)
- **PostgreSQL** instance running and accepting connections
- **GNU Make** (installed by default on macOS/Linux; for Windows use Git Bash or chocolatey)
- **Goose CLI** (`go install github.com/pressly/goose/v3/cmd/goose@latest` to apply migrations or run make commands)

Your database must be reachable, and the `.env` file must be properly configured (see [Getting Started](#getting-started) and the `TEST_DB_*` optional variables below).

Run all DB constraint tests:

```bash
make test-db
```

This runs `go test ./internal/app/db/migrations/... -v -count=1` and covers:

| Table | Constraint | Test |
|---|---|---|
| `users` | Email unique (case-insensitive) | `TestUsers_EmailUnique_CaseInsensitive` |
| `users` | Phone unique | `TestUsers_PhoneUnique` |
| `social_accounts` | FK → users | `TestSocialAccounts_FK_InvalidUser` |
| `social_accounts` | Provider + provider_user_id unique | `TestSocialAccounts_ProviderUnique` |
| `social_accounts` | Cascade delete on user | `TestSocialAccounts_Cascade_UserDelete` |
| `refresh_tokens` | FK → users | `TestRefreshTokens_FK_InvalidUser` |
| `refresh_tokens` | Cascade delete on user | `TestRefreshTokens_Cascade_UserDelete` |
| `user_sessions` | FK → users | `TestUserSessions_FK_InvalidUser` |
| `user_sessions` | FK → refresh_tokens | `TestUserSessions_FK_InvalidRefreshToken` |
| `user_sessions` | Cascade delete on user | `TestUserSessions_Cascade_UserDelete` |
| `user_sessions` | Cascade delete on token | `TestUserSessions_Cascade_TokenDelete` |

> **Note:** Each test runs Goose `migrate up` before and `migrate reset` after, so the DB is always left in a clean state.

To use a separate test database, set `TEST_DB_*` environment variables (they override `DB_*`):

```env
TEST_DB_HOST=localhost
TEST_DB_PORT=5432
TEST_DB_USER=your_test_user
TEST_DB_PASSWORD=your_test_password
TEST_DB_NAME=healmata_test
TEST_DB_SSLMODE=disable
```

### All Tests

```bash
go test ./...
```

---
