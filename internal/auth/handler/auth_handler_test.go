package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"

	"healmata_backend/internal/auth/dto"
	"healmata_backend/internal/testutils"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success cases", func(t *testing.T) {
		tests := []struct {
			name         string
			reqBody      dto.RegisterRequestDTO
			checkDBEmail string
			checkDBPhone string
		}{
			{
				name: "Register with email",
				reqBody: dto.RegisterRequestDTO{
					FullName:        "Nguyen Van A",
					Identifier:      "testuser@example.com",
					Password:        "Password123!",
					ConfirmPassword: "Password123!",
				},
				checkDBEmail: "testuser@example.com",
			},
			{
				name: "Register with phone",
				reqBody: dto.RegisterRequestDTO{
					FullName:        "Nguyen Van B",
					Identifier:      "+84901234567",
					Password:        "Password123!",
					ConfirmPassword: "Password123!",
				},
				checkDBPhone: "+84901234567",
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				//=======================================================================
				// 1. Setup DB và Mock
				pool := testutils.SetupTestDB(t)
				emailSender := testutils.NewMockEmailSender()

				// 2. Khởi tạo test router
				r := testutils.SetupTestRouter(pool, emailSender)
				//=======================================================================

				bodyBytes, err := json.Marshal(tc.reqBody)
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

				data := resp["data"].(map[string]interface{})
				assert.NotEmpty(t, data["accessToken"])
				assert.NotEmpty(t, data["refreshToken"])
				assert.Greater(t, int64(data["expiresIn"].(float64)), int64(0))

				// DB Verification
				ctx := context.Background()
				var dbUserID string
				if tc.checkDBEmail != "" {
					err = pool.QueryRow(ctx, "SELECT id FROM users WHERE email = $1", tc.checkDBEmail).Scan(&dbUserID)
				} else {
					err = pool.QueryRow(ctx, "SELECT id FROM users WHERE phone = $1", tc.checkDBPhone).Scan(&dbUserID)
				}
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
		}
	})

	t.Run("Failure cases", func(t *testing.T) {
		tests := []struct {
			name            string
			reqBody         dto.RegisterRequestDTO
			setupFunc       func(t *testing.T, pool *pgxpool.Pool)
			isValidationErr bool
			expectedCode    string
			expectedMsg     string
			expectedHTTP    int
		}{
			{
				name: "Duplicate Email",
				reqBody: dto.RegisterRequestDTO{
					FullName:        "Nguyen Van C",
					Identifier:      "existing@example.com",
					Password:        "Password123!",
					ConfirmPassword: "Password123!",
				},
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO users (full_name, email, password_hash)
						VALUES ($1, $2, $3)
					`, "Existing User", "existing@example.com", "dummyhash")
					require.NoError(t, err)
				},
				isValidationErr: false,
				expectedCode:    "AUTH_REG_001",
				expectedMsg:     "EMAIL_EXISTS",
				expectedHTTP:    http.StatusConflict,
			},
			{
				name: "Duplicate Phone",
				reqBody: dto.RegisterRequestDTO{
					FullName:        "Nguyen Van D",
					Identifier:      "+84901234567",
					Password:        "Password123!",
					ConfirmPassword: "Password123!",
				},
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO users (full_name, phone, password_hash)
						VALUES ($1, $2, $3)
					`, "Existing User", "+84901234567", "dummyhash")
					require.NoError(t, err)
				},
				isValidationErr: false,
				expectedCode:    "AUTH_REG_002",
				expectedMsg:     "PHONE_EXISTS",
				expectedHTTP:    http.StatusConflict,
			},
			{
				name: "Password mismatch",
				reqBody: dto.RegisterRequestDTO{
					FullName:        "Nguyen Van E",
					Identifier:      "testuser2@example.com",
					Password:        "Password123!",
					ConfirmPassword: "DifferentPassword123!",
				},
				isValidationErr: true,
				expectedCode:    "VR-CONFIRM-002",
				expectedHTTP:    http.StatusBadRequest,
			},
			{
				name: "Invalid email format",
				reqBody: dto.RegisterRequestDTO{
					FullName:        "Nguyen Van F",
					Identifier:      "invalid-email@",
					Password:        "Password123!",
					ConfirmPassword: "Password123!",
				},
				isValidationErr: true,
				expectedCode:    "VR-ID-002",
				expectedHTTP:    http.StatusBadRequest,
			},
			{
				name: "Invalid phone format",
				reqBody: dto.RegisterRequestDTO{
					FullName:        "Nguyen Van G",
					Identifier:      "123abc456",
					Password:        "Password123!",
					ConfirmPassword: "Password123!",
				},
				isValidationErr: true,
				expectedCode:    "VR-ID-002",
				expectedHTTP:    http.StatusBadRequest,
			},
			{
				name: "Invalid name length",
				reqBody: dto.RegisterRequestDTO{
					FullName:        "A",
					Identifier:      "validuser@example.com",
					Password:        "Password123!",
					ConfirmPassword: "Password123!",
				},
				isValidationErr: true,
				expectedCode:    "VR-NAME-004",
				expectedHTTP:    http.StatusBadRequest,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				pool := testutils.SetupTestDB(t)
				emailSender := testutils.NewMockEmailSender()
				r := testutils.SetupTestRouter(pool, emailSender)

				if tc.setupFunc != nil {
					tc.setupFunc(t, pool)
				}

				bodyBytes, err := json.Marshal(tc.reqBody)
				require.NoError(t, err)

				req, err := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(bodyBytes))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedHTTP, w.Code)

				var resp map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.False(t, resp["success"].(bool))
				errMap := resp["error"].(map[string]interface{})

				// Xử lý logic check lỗi tùy theo loại lỗi (Validation vs Business)
				if tc.isValidationErr {
					assert.Equal(t, "VALIDATION_ERROR", errMap["code"])
					details := errMap["details"].([]interface{})
					assert.NotEmpty(t, details)
					firstDetail := details[0].(map[string]interface{})
					assert.Equal(t, tc.expectedCode, firstDetail["code"])
				} else {
					assert.Equal(t, tc.expectedCode, errMap["code"])
					if tc.expectedMsg != "" {
						assert.Equal(t, tc.expectedMsg, errMap["message"])
					}
				}
			})
		}
	})
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	hashPassword := func(pwd string) string {
		h, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		return string(h)
	}

	t.Run("Success cases", func(t *testing.T) {
		tests := []struct {
			name      string
			setupFunc func(t *testing.T, pool *pgxpool.Pool)
			reqBody   dto.LoginRequestDTO
		}{
			{
				name: "Login with email",
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO users (full_name, email, password_hash, status)
						VALUES ($1, $2, $3, $4)
					`, "Nguyen Van A", "testlogin@example.com", hashPassword("Password123!"), "active")
					require.NoError(t, err)
				},
				reqBody: dto.LoginRequestDTO{
					Identifier: "testlogin@example.com",
					Password:   "Password123!",
				},
			},
			{
				name: "Login with phone",
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO users (full_name, phone, password_hash, status)
						VALUES ($1, $2, $3, $4)
					`, "Nguyen Van B", "+84901234567", hashPassword("Password123!"), "active")
					require.NoError(t, err)
				},
				reqBody: dto.LoginRequestDTO{
					Identifier: "+84901234567",
					Password:   "Password123!",
				},
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				pool := testutils.SetupTestDB(t)
				emailSender := testutils.NewMockEmailSender()
				r := testutils.SetupTestRouter(pool, emailSender)

				if tc.setupFunc != nil {
					tc.setupFunc(t, pool)
				}

				bodyBytes, err := json.Marshal(tc.reqBody)
				require.NoError(t, err)

				req, err := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(bodyBytes))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)

				var resp map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.True(t, resp["success"].(bool))
				assert.Equal(t, "LOGIN_SUCCESS", resp["message"])

				data := resp["data"].(map[string]interface{})
				assert.NotEmpty(t, data["accessToken"])
				assert.NotEmpty(t, data["refreshToken"])
				assert.Greater(t, int64(data["expiresIn"].(float64)), int64(0))
			})
		}
	})

	t.Run("Failure cases", func(t *testing.T) {
		tests := []struct {
			name            string
			setupFunc       func(t *testing.T, pool *pgxpool.Pool)
			reqBody         dto.LoginRequestDTO
			isValidationErr bool
			expectedCode    string
			expectedMsg     string
			expectedHTTP    int
		}{
			{
				name: "Invalid credentials (wrong password)",
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO users (full_name, email, password_hash, status)
						VALUES ($1, $2, $3, $4)
					`, "Nguyen Van A", "testlogin@example.com", hashPassword("Password123!"), "active")
					require.NoError(t, err)
				},
				reqBody: dto.LoginRequestDTO{
					Identifier: "testlogin@example.com",
					Password:   "WrongPassword!",
				},
				isValidationErr: false,
				expectedCode:    "AUTH_LOGIN_001",
				expectedMsg:     "INVALID_CREDENTIAL",
				expectedHTTP:    http.StatusUnauthorized,
			},
			{
				name: "User not found",
				reqBody: dto.LoginRequestDTO{
					Identifier: "nonexistent@example.com",
					Password:   "Password123!",
				},
				isValidationErr: false,
				expectedCode:    "AUTH_LOGIN_002",
				expectedMsg:     "USER_NOT_FOUND",
				expectedHTTP:    http.StatusNotFound,
			},
			{
				name: "User disabled",
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO users (full_name, email, password_hash, status)
						VALUES ($1, $2, $3, $4)
					`, "Nguyen Van A", "testlogin@example.com", hashPassword("Password123!"), "disabled")
					require.NoError(t, err)
				},
				reqBody: dto.LoginRequestDTO{
					Identifier: "testlogin@example.com",
					Password:   "Password123!",
				},
				isValidationErr: false,
				expectedCode:    "AUTH_LOGIN_003",
				expectedMsg:     "USER_DISABLED",
				expectedHTTP:    http.StatusForbidden,
			},
			{
				name: "Validation error (invalid email format)",
				reqBody: dto.LoginRequestDTO{
					Identifier: "invalid-email@",
					Password:   "Password123!",
				},
				isValidationErr: true,
				expectedCode:    "VR-ID-002",
				expectedHTTP:    http.StatusBadRequest,
			},
			{
				name: "Validation error (invalid password format)",
				reqBody: dto.LoginRequestDTO{
					Identifier: "testlogin@example.com",
					Password:   "short",
				},
				isValidationErr: true,
				expectedCode:    "VR-PWD-002",
				expectedHTTP:    http.StatusBadRequest,
			},
		}

		for _, tc := range tests {
			t.Run(tc.name, func(t *testing.T) {
				pool := testutils.SetupTestDB(t)
				emailSender := testutils.NewMockEmailSender()
				r := testutils.SetupTestRouter(pool, emailSender)

				if tc.setupFunc != nil {
					tc.setupFunc(t, pool)
				}

				bodyBytes, err := json.Marshal(tc.reqBody)
				require.NoError(t, err)

				req, err := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(bodyBytes))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				assert.Equal(t, tc.expectedHTTP, w.Code)

				var resp map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.False(t, resp["success"].(bool))
				errMap := resp["error"].(map[string]interface{})

				if tc.isValidationErr {
					assert.Equal(t, "VALIDATION_ERROR", errMap["code"])
					details := errMap["details"].([]interface{})
					assert.NotEmpty(t, details)
					firstDetail := details[0].(map[string]interface{})
					assert.Equal(t, tc.expectedCode, firstDetail["code"])
				} else {
					assert.Equal(t, tc.expectedCode, errMap["code"])
					assert.Equal(t, tc.expectedMsg, errMap["message"])
				}
			})
		}
	})
}
