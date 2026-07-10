package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"healmata_backend/internal/app/router"
	"healmata_backend/internal/auth/dto"
	"healmata_backend/internal/db/testhelper"
)

func TestRegister(t *testing.T) {
	// Set Gin to test mode to avoid excessive log output
	gin.SetMode(gin.TestMode)

	t.Run("Success - Register with email", func(t *testing.T) {
		pool := testhelper.SetupTestDB(t)
		r := gin.New()
		router.RegisterRoutes(r, pool)

		reqBody := dto.RegisterRequestDTO{
			FullName:        "Nguyen Van A",
			Identifier:      "testuser@example.com",
			Password:        "Password123!",
			ConfirmPassword: "Password123!",
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("Success test failed with code %d, body: %s", w.Code, w.Body.String())
		}

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.True(t, resp["success"].(bool))
		assert.Equal(t, "REGISTER_SUCCESS", resp["message"])

		data := resp["data"].(map[string]interface{})
		assert.NotEmpty(t, data["accessToken"])
		assert.NotEmpty(t, data["refreshToken"])
		assert.Greater(t, int64(data["expiresIn"].(float64)), int64(0))

		// Explicit Database Verifications
		ctx := context.Background()
		var dbUserID string
		err = pool.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", "testuser@example.com").Scan(&dbUserID)
		assert.NoError(t, err)
		assert.NotEmpty(t, dbUserID)

		var tokenHash string
		var tokenID string
		err = pool.QueryRow(ctx, "SELECT id, token_hash FROM refresh_tokens WHERE user_id = $1", dbUserID).Scan(&tokenID, &tokenHash)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenID)
		assert.NotEmpty(t, tokenHash)

		var sessionID string
		var refTokenID string
		err = pool.QueryRow(ctx, "SELECT id, refresh_token_id FROM user_sessions WHERE user_id = $1", dbUserID).Scan(&sessionID, &refTokenID)
		assert.NoError(t, err)
		assert.NotEmpty(t, sessionID)
		assert.Equal(t, tokenID, refTokenID)
	})

	t.Run("Success - Register with phone", func(t *testing.T) {
		pool := testhelper.SetupTestDB(t)
		r := gin.New()
		router.RegisterRoutes(r, pool)

		reqBody := dto.RegisterRequestDTO{
			FullName:        "Nguyen Van B",
			Identifier:      "+84901234567",
			Password:        "Password123!",
			ConfirmPassword: "Password123!",
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.True(t, resp["success"].(bool))
		assert.Equal(t, "REGISTER_SUCCESS", resp["message"])

		// Explicit Database Verifications
		ctx := context.Background()
		var dbUserID string
		err = pool.QueryRow(ctx, "SELECT id FROM users WHERE phone = $1", "+84901234567").Scan(&dbUserID)
		assert.NoError(t, err)
		assert.NotEmpty(t, dbUserID)

		var tokenHash string
		var tokenID string
		err = pool.QueryRow(ctx, "SELECT id, token_hash FROM refresh_tokens WHERE user_id = $1", dbUserID).Scan(&tokenID, &tokenHash)
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenID)
		assert.NotEmpty(t, tokenHash)

		var sessionID string
		var refTokenID string
		err = pool.QueryRow(ctx, "SELECT id, refresh_token_id FROM user_sessions WHERE user_id = $1", dbUserID).Scan(&sessionID, &refTokenID)
		assert.NoError(t, err)
		assert.NotEmpty(t, sessionID)
		assert.Equal(t, tokenID, refTokenID)
	})

	t.Run("Failure - Duplicate Email", func(t *testing.T) {
		pool := testhelper.SetupTestDB(t)
		r := gin.New()
		router.RegisterRoutes(r, pool)

		// Seed user directly into database
		ctx := context.Background()
		_, err := pool.Exec(ctx, `
			INSERT INTO users (full_name, email, password_hash)
			VALUES ($1, $2, $3)
		`, "Existing User", "existing@example.com", "dummyhash")
		require.NoError(t, err)

		reqBody := dto.RegisterRequestDTO{
			FullName:        "Nguyen Van C",
			Identifier:      "existing@example.com",
			Password:        "Password123!",
			ConfirmPassword: "Password123!",
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.False(t, resp["success"].(bool))
		errMap := resp["error"].(map[string]interface{})
		assert.Equal(t, "AUTH_REG_001", errMap["code"])
		assert.Equal(t, "EMAIL_EXISTS", errMap["message"])
	})

	t.Run("Failure - Duplicate Phone", func(t *testing.T) {
		pool := testhelper.SetupTestDB(t)
		r := gin.New()
		router.RegisterRoutes(r, pool)

		// Seed user directly into database
		ctx := context.Background()
		_, err := pool.Exec(ctx, `
			INSERT INTO users (full_name, phone, password_hash)
			VALUES ($1, $2, $3)
		`, "Existing User", "+84901234567", "dummyhash")
		require.NoError(t, err)

		reqBody := dto.RegisterRequestDTO{
			FullName:        "Nguyen Van D",
			Identifier:      "+84901234567",
			Password:        "Password123!",
			ConfirmPassword: "Password123!",
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusConflict, w.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.False(t, resp["success"].(bool))
		errMap := resp["error"].(map[string]interface{})
		assert.Equal(t, "AUTH_REG_002", errMap["code"])
		assert.Equal(t, "PHONE_EXISTS", errMap["message"])
	})

	t.Run("Failure - Validation error (password mismatch)", func(t *testing.T) {
		pool := testhelper.SetupTestDB(t)
		r := gin.New()
		router.RegisterRoutes(r, pool)

		reqBody := dto.RegisterRequestDTO{
			FullName:        "Nguyen Van E",
			Identifier:      "testuser2@example.com",
			Password:        "Password123!",
			ConfirmPassword: "DifferentPassword123!",
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.False(t, resp["success"].(bool))
		errMap := resp["error"].(map[string]interface{})
		assert.Equal(t, "AUTH_REG_008", errMap["code"])
	})

	t.Run("Failure - Validation error (invalid email format)", func(t *testing.T) {
		pool := testhelper.SetupTestDB(t)
		r := gin.New()
		router.RegisterRoutes(r, pool)

		reqBody := dto.RegisterRequestDTO{
			FullName:        "Nguyen Van F",
			Identifier:      "invalid-email@",
			Password:        "Password123!",
			ConfirmPassword: "Password123!",
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.False(t, resp["success"].(bool))
		errMap := resp["error"].(map[string]interface{})
		assert.Equal(t, "AUTH_REG_006", errMap["code"])
		assert.Equal(t, "INVALID_EMAIL", errMap["message"])
	})

	t.Run("Failure - Validation error (invalid phone format)", func(t *testing.T) {
		pool := testhelper.SetupTestDB(t)
		r := gin.New()
		router.RegisterRoutes(r, pool)

		reqBody := dto.RegisterRequestDTO{
			FullName:        "Nguyen Van G",
			Identifier:      "123abc456",
			Password:        "Password123!",
			ConfirmPassword: "Password123!",
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.False(t, resp["success"].(bool))
		errMap := resp["error"].(map[string]interface{})
		assert.Equal(t, "AUTH_REG_007", errMap["code"])
		assert.Equal(t, "INVALID_PHONE", errMap["message"])
	})

	t.Run("Failure - Validation error (invalid name length)", func(t *testing.T) {
		pool := testhelper.SetupTestDB(t)
		r := gin.New()
		router.RegisterRoutes(r, pool)

		reqBody := dto.RegisterRequestDTO{
			FullName:        "A",
			Identifier:      "validuser@example.com",
			Password:        "Password123!",
			ConfirmPassword: "Password123!",
		}

		bodyBytes, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req, err := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(bodyBytes))
		require.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)

		var resp map[string]interface{}
		err = json.Unmarshal(w.Body.Bytes(), &resp)
		require.NoError(t, err)

		assert.False(t, resp["success"].(bool))
		errMap := resp["error"].(map[string]interface{})
		assert.Equal(t, "AUTH_REG_004", errMap["code"])
		assert.Equal(t, "INVALID_NAME", errMap["message"])
	})
}
