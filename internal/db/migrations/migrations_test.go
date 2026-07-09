// Package migrations_test contains integration tests that verify all
// PostgreSQL constraints defined in the auth schema migrations.
//
// Each test runs against a real database. The testhelper.SetupTestDB helper
// applies all Goose migrations before tests and resets them after.
//
// PostgreSQL error code reference:
//
//	23505 — unique_violation (duplicate key)
//	23503 — foreign_key_violation (referencing non-existent row)
package migrations_test

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"healmata_backend/internal/db/testhelper"
)

// pgErrCode extracts the PostgreSQL error code string from a pgx error.
// Returns empty string if the error is not a *pgconn.PgError.
func pgErrCode(err error) string {
	if e, ok := err.(*pgconn.PgError); ok {
		return string(e.Code)
	}
	return ""
}

// ──────────────────────────────────────────────────────────────────────────────
// users table — 2 tests
// ──────────────────────────────────────────────────────────────────────────────

// TC-BE-GO-002-003 — Duplicate Email
func TestUsers_EmailUnique_CaseInsensitive(t *testing.T) {
	pool := testhelper.SetupTestDB(t)
	ctx := context.Background()

	_, err := pool.Exec(ctx,
		`INSERT INTO users (full_name, email) VALUES ($1, $2)`,
		"Alice", "alice@test.com",
	)
	require.NoError(t, err, "first insert should succeed")

	_, err = pool.Exec(ctx,
		`INSERT INTO users (full_name, email) VALUES ($1, $2)`,
		"Alice Duplicate", "ALICE@TEST.COM",
	)
	require.Error(t, err, "duplicate email (different case) must be rejected")
	assert.Equal(t, "23505", pgErrCode(err), "expect unique_violation (23505)")
}

// TC-BE-GO-002-004 — Duplicate Phone
func TestUsers_PhoneUnique(t *testing.T) {
	pool := testhelper.SetupTestDB(t)
	ctx := context.Background()

	_, err := pool.Exec(ctx,
		`INSERT INTO users (full_name, phone) VALUES ($1, $2)`,
		"Bob", "+84901234567",
	)
	require.NoError(t, err, "first insert should succeed")

	_, err = pool.Exec(ctx,
		`INSERT INTO users (full_name, phone) VALUES ($1, $2)`,
		"Bob Clone", "+84901234567",
	)
	require.Error(t, err, "duplicate phone must be rejected")
	assert.Equal(t, "23505", pgErrCode(err), "expect unique_violation (23505)")
}

// ──────────────────────────────────────────────────────────────────────────────
// social_accounts table — 3 tests
// ──────────────────────────────────────────────────────────────────────────────

// TC-BE-GO-002-005 — Foreign Key Social Account
func TestSocialAccounts_FK_InvalidUser(t *testing.T) {
	pool := testhelper.SetupTestDB(t)
	ctx := context.Background()

	_, err := pool.Exec(ctx,
		`INSERT INTO social_accounts (user_id, provider, provider_user_id, provider_email)
		 VALUES ($1, $2, $3, $4)`,
		"00000000-0000-0000-0000-000000000000", "google", "ghost-id", "ghost@g.com",
	)
	require.Error(t, err, "FK must reject non-existent user_id")
	assert.Equal(t, "23503", pgErrCode(err), "expect foreign_key_violation (23503)")
}

// TC-BE-GO-002-005-001 — Provider Unique
func TestSocialAccounts_ProviderUnique(t *testing.T) {
	pool := testhelper.SetupTestDB(t)
	ctx := context.Background()

	var userID string
	err := pool.QueryRow(ctx,
		`INSERT INTO users (full_name, email) VALUES ($1, $2) RETURNING id`,
		"Charlie Provider", "charlie.provider@test.com",
	).Scan(&userID)
	require.NoError(t, err)

	// Happy Path
	t.Run("Create_Social_Account", func(t *testing.T) {
		_, err = pool.Exec(ctx,
			`INSERT INTO social_accounts (user_id, provider, provider_user_id, provider_email)
		 VALUES ($1, $2, $3, $4)`,
			userID, "google", "google-uid-001", "charlie@g.com",
		)
		require.NoError(t, err, "first social_account insert should succeed")
	})

	// Same provider, different provider_user_id → different composite key → MUST succeed
	t.Run("DifferentProviderID_SameEmail_Success", func(t *testing.T) {
		_, err = pool.Exec(ctx,
			`INSERT INTO social_accounts (user_id, provider, provider_user_id, provider_email)
		 VALUES ($1, $2, $3, $4)`,
			userID, "google", "google-uid-002", "charlie@g.com",
		)
		require.NoError(t, err, "same provider but different provider_user_id must succeed (different composite key)")
	})

	// Same provider_user_id, diff email
	t.Run("SameProviderID_DiffEmail_Fails", func(t *testing.T) {
		_, err = pool.Exec(ctx,
			`INSERT INTO social_accounts (user_id, provider, provider_user_id, provider_email)
		 VALUES ($1, $2, $3, $4)`,
			userID, "google", "google-uid-001", "charlie-dup@g.com",
		)
		require.Error(t, err, "duplicate (provider, provider_user_id) must be rejected")
		assert.Equal(t, "23505", pgErrCode(err), "expect unique_violation (23505)")
	})

	// Different provider, same user_id, same provider_user_id
	t.Run("DiffProvider_SameProviderID_Success", func(t *testing.T) {
		_, err = pool.Exec(ctx,
			`INSERT INTO social_accounts (user_id, provider, provider_user_id, provider_email)
		 VALUES ($1, $2, $3, $4)`,
			userID, "apple", "google-uid-001", "charlie@g.com",
		)
		require.NoError(t, err, "different provider, same user_id, same provider_user_id should succeed")
	})
}

func TestSocialAccounts_CompositeUniqueConstraint(t *testing.T) {
	pool := testhelper.SetupTestDB(t)
	ctx := context.Background()

	var userA, userB string
	err := pool.QueryRow(ctx,
		`INSERT INTO users (full_name, email) VALUES ($1, $2) RETURNING id`,
		"User A", "user.a@test.com",
	).Scan(&userA)
	require.NoError(t, err)

	err = pool.QueryRow(ctx,
		`INSERT INTO users (full_name, email) VALUES ($1, $2) RETURNING id`,
		"User B", "user.b@test.com",
	).Scan(&userB)
	require.NoError(t, err)

	// Base Insert: User A links Google
	_, err = pool.Exec(ctx,
		`INSERT INTO social_accounts (user_id, provider, provider_user_id, provider_email)
		 VALUES ($1, $2, $3, $4)`,
		userA, "google", "google-uid-001", "user.a@g.com",
	)
	require.NoError(t, err, "Base insert must succeed")

	// Insert same provider_user_id for User A
	t.Run("DifferentUser_SameSocial_Fails", func(t *testing.T) {
		_, err := pool.Exec(ctx,
			`INSERT INTO social_accounts (user_id, provider, provider_user_id, provider_email)
			 VALUES ($1, $2, $3, $4)`,
			userB, "google", "google-uid-001", "user.b@g.com",
		)
		require.Error(t, err, "Different user, same provider_user_id must fail")
		assert.Equal(t, "23505", pgErrCode(err), "expect unique_violation (23505)")
	})

	t.Run("SameUser_SameProvider_DifferentID_Succeeds", func(t *testing.T) {
		_, err := pool.Exec(ctx,
			`INSERT INTO social_accounts (user_id, provider, provider_user_id, provider_email)
			 VALUES ($1, $2, $3, $4)`,
			userA, "google", "google-uid-002", "user.a@g.com",
		)
		require.NoError(t, err, "Same user, same provider, different provider_user_id must succeed")
	})

}

// TestSocialAccounts_Cascade_UserDelete verifies ON DELETE CASCADE:
// deleting a user must automatically remove all linked social_accounts.
func TestSocialAccounts_Cascade_UserDelete(t *testing.T) {
	pool := testhelper.SetupTestDB(t)
	ctx := context.Background()

	var userID string
	err := pool.QueryRow(ctx,
		`INSERT INTO users (full_name, email) VALUES ($1, $2) RETURNING id`,
		"Charlie Cascade", "charlie.cascade@test.com",
	).Scan(&userID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx,
		`INSERT INTO social_accounts (user_id, provider, provider_user_id, provider_email)
		 VALUES ($1, $2, $3, $4)`,
		userID, "apple", "apple-uid-001", "charlie@apple.com",
	)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	require.NoError(t, err)

	var count int
	err = pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM social_accounts WHERE user_id = $1`, userID,
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "social_account must be deleted when user is deleted (CASCADE)")
}

// ──────────────────────────────────────────────────────────────────────────────
// refresh_tokens table — 2 tests
// ──────────────────────────────────────────────────────────────────────────────

// TestRefreshTokens_FK_InvalidUser verifies that a refresh_token with a
// non-existent user_id is rejected by the FK constraint.
func TestRefreshTokens_FK_InvalidUser(t *testing.T) {
	pool := testhelper.SetupTestDB(t)
	ctx := context.Background()

	_, err := pool.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, device_id, expires_at)
		 VALUES ($1, $2, $3, $4)`,
		"00000000-0000-0000-0000-000000000000",
		"fake-hash", "device-001",
		time.Now().Add(30*24*time.Hour),
	)
	require.Error(t, err, "FK must reject non-existent user_id")
	assert.Equal(t, "23503", pgErrCode(err), "expect foreign_key_violation (23503)")
}

// TestRefreshTokens_Cascade_UserDelete verifies ON DELETE CASCADE:
// deleting a user must automatically remove all linked refresh_tokens.
func TestRefreshTokens_Cascade_UserDelete(t *testing.T) {
	pool := testhelper.SetupTestDB(t)
	ctx := context.Background()

	var userID string
	err := pool.QueryRow(ctx,
		`INSERT INTO users (full_name, email) VALUES ($1, $2) RETURNING id`,
		"Dave Cascade", "dave.cascade@test.com",
	).Scan(&userID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, device_id, expires_at)
		 VALUES ($1, $2, $3, $4)`,
		userID, "hash-abc", "device-dave", time.Now().Add(30*24*time.Hour),
	)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	require.NoError(t, err)

	var count int
	err = pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM refresh_tokens WHERE user_id = $1`, userID,
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "refresh_token must be deleted when user is deleted (CASCADE)")
}

// ──────────────────────────────────────────────────────────────────────────────
// user_sessions table — 4 tests (most complex: 2 FKs, 2 cascade paths)
// ──────────────────────────────────────────────────────────────────────────────

// setupUserAndToken is a shared test helper that inserts a user + refresh_token
// and returns both UUIDs. Used by user_sessions tests.
func setupUserAndToken(t *testing.T, ctx context.Context, pool *pgxpool.Pool, label string) (userID, tokenID string) {
	t.Helper()

	err := pool.QueryRow(ctx,
		`INSERT INTO users (full_name, email) VALUES ($1, $2) RETURNING id`,
		label+" User", label+"-sessions@test.com",
	).Scan(&userID)
	require.NoError(t, err, "setup: insert user for "+label)

	err = pool.QueryRow(ctx,
		`INSERT INTO refresh_tokens (user_id, token_hash, device_id, expires_at)
		 VALUES ($1, $2, $3, $4) RETURNING id`,
		userID, "hash-"+label, "device-"+label, time.Now().Add(30*24*time.Hour),
	).Scan(&tokenID)
	require.NoError(t, err, "setup: insert refresh_token for "+label)
	return
}

// TestUserSessions_FK_InvalidUser verifies that a session with a
// non-existent user_id is rejected by the FK constraint.
func TestUserSessions_FK_InvalidUser(t *testing.T) {
	pool := testhelper.SetupTestDB(t)
	ctx := context.Background()

	_, err := pool.Exec(ctx,
		`INSERT INTO user_sessions
		 (user_id, refresh_token_id, device_id, platform, ip_address, user_agent)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		"00000000-0000-0000-0000-000000000000",
		"00000000-0000-0000-0000-000000000001",
		"device-x", "ios", "127.0.0.1", "TestAgent/1.0",
	)
	require.Error(t, err, "FK must reject non-existent user_id in user_sessions")
	assert.Equal(t, "23503", pgErrCode(err), "expect foreign_key_violation (23503)")
}

// TestUserSessions_FK_InvalidRefreshToken verifies that a session with a valid
// user_id but non-existent refresh_token_id is rejected by the FK constraint.
func TestUserSessions_FK_InvalidRefreshToken(t *testing.T) {
	pool := testhelper.SetupTestDB(t)
	ctx := context.Background()

	var userID string
	err := pool.QueryRow(ctx,
		`INSERT INTO users (full_name, email) VALUES ($1, $2) RETURNING id`,
		"Eve FK Test", "eve.fk@test.com",
	).Scan(&userID)
	require.NoError(t, err)

	_, err = pool.Exec(ctx,
		`INSERT INTO user_sessions
		 (user_id, refresh_token_id, device_id, platform, ip_address, user_agent)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		userID,
		"00000000-0000-0000-0000-000000000099", // non-existent token
		"device-eve", "android", "192.168.1.1", "TestAgent/1.0",
	)
	require.Error(t, err, "FK must reject non-existent refresh_token_id in user_sessions")
	assert.Equal(t, "23503", pgErrCode(err), "expect foreign_key_violation (23503)")
}

// TestUserSessions_Cascade_UserDelete verifies that deleting a user cascades
// and removes all sessions belonging to that user.
func TestUserSessions_Cascade_UserDelete(t *testing.T) {
	pool := testhelper.SetupTestDB(t)
	ctx := context.Background()

	userID, tokenID := setupUserAndToken(t, ctx, pool, "UserCascade")

	_, err := pool.Exec(ctx,
		`INSERT INTO user_sessions
		 (user_id, refresh_token_id, device_id, platform, ip_address, user_agent)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		userID, tokenID, "device-1", "ios", "10.0.0.1", "TestAgent/1.0",
	)
	require.NoError(t, err)

	_, err = pool.Exec(ctx, `DELETE FROM users WHERE id = $1`, userID)
	require.NoError(t, err)

	var count int
	err = pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM user_sessions WHERE user_id = $1`, userID,
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "session must be deleted when user is deleted (CASCADE)")
}

// TestUserSessions_Cascade_TokenDelete verifies that deleting a refresh_token
// cascades and removes all sessions linked to that token.
func TestUserSessions_Cascade_TokenDelete(t *testing.T) {
	pool := testhelper.SetupTestDB(t)
	ctx := context.Background()

	userID, tokenID := setupUserAndToken(t, ctx, pool, "TokenCascade")

	_, err := pool.Exec(ctx,
		`INSERT INTO user_sessions
		 (user_id, refresh_token_id, device_id, platform, ip_address, user_agent)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		userID, tokenID, "device-2", "android", "10.0.0.2", "TestAgent/1.0",
	)
	require.NoError(t, err)

	// Delete only the token — user remains
	_, err = pool.Exec(ctx, `DELETE FROM refresh_tokens WHERE id = $1`, tokenID)
	require.NoError(t, err)

	var count int
	err = pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM user_sessions WHERE refresh_token_id = $1`, tokenID,
	).Scan(&count)
	require.NoError(t, err)
	assert.Equal(t, 0, count, "session must be deleted when refresh_token is deleted (CASCADE)")
}
