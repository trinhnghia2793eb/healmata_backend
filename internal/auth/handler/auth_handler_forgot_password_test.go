package handler_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"healmata_backend/internal/auth/dto"
	"healmata_backend/internal/testutils"
)

var (
	validOTP        = "123456"
	validRequestID  = "123e4567-e89b-12d3-a456-426614174000"
	validResetToken = "valid-reset-token-xyz"
)

func getValidOTPHash() string {
	hash := sha256.Sum256([]byte(validOTP))
	return hex.EncodeToString(hash[:])
}

func getValidResetTokenHash() string {
	resetHash := sha256.Sum256([]byte(validResetToken))
	return hex.EncodeToString(resetHash[:])
}

// ====================================================================================
// TEST FORGOT PASSWORD
// ====================================================================================
func TestForgotPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success cases", func(t *testing.T) {
		tests := []struct {
			name         string
			reqBody      dto.ForgotPasswordRequestDTO
			setupFunc    func(t *testing.T, pool *pgxpool.Pool)
			checkDBEmail string
		}{
			{
				name: "TC1.1 - Success (Email)",
				reqBody: dto.ForgotPasswordRequestDTO{
					Identifier: "valid@example.com",
				},
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO users (full_name, email, password_hash)
						VALUES ('User Email', 'valid@example.com', 'dummyhash')
					`)
					require.NoError(t, err)
				},
				checkDBEmail: "valid@example.com",
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

				req, err := http.NewRequest("POST", "/v1/auth/forgot-password", bytes.NewBuffer(bodyBytes))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)

				var resp map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.True(t, resp["success"].(bool))
				assert.Equal(t, "OTP_SENT", resp["message"])

				data, ok := resp["data"].(map[string]interface{})
				require.True(t, ok, "Trường data bị thiếu hoặc không đúng định dạng")
				assert.NotEmpty(t, data["resetRequestId"])
				assert.Equal(t, float64(6), data["otpLength"])
				assert.Equal(t, float64(300), data["expiresIn"])
				assert.Equal(t, float64(60), data["resendAfter"])

				// DB Verification
				if tc.checkDBEmail != "" {
					var count int
					err := pool.QueryRow(context.Background(), `SELECT COUNT(*) FROM otp_requests WHERE identifier = $1 AND purpose = 'reset_password'`, tc.checkDBEmail).Scan(&count)
					require.NoError(t, err)
					assert.Equal(t, 1, count, "Nên có 1 record OTP được tạo trong DB")
				}
			})
		}
	})

	t.Run("Failure cases", func(t *testing.T) {
		tests := []struct {
			name            string
			reqBody         dto.ForgotPasswordRequestDTO
			setupFunc       func(t *testing.T, pool *pgxpool.Pool)
			isValidationErr bool
			expectedCode    string
			expectedMsg     string
			expectedHTTP    int
		}{
			{
				name: "TC1.3 - Failure (User Not Found)",
				reqBody: dto.ForgotPasswordRequestDTO{
					Identifier: "notfound@example.com",
				},
				isValidationErr: false,
				expectedCode:    "AUTH_FORGOT_001",
				expectedMsg:     "USER_NOT_FOUND",
				expectedHTTP:    http.StatusNotFound,
			},
			{
				name: "TC1.4 - Failure (Validation - Empty Identifier)",
				reqBody: dto.ForgotPasswordRequestDTO{
					Identifier: "",
				},
				isValidationErr: true,
				expectedCode:    "VR-ID-001",
				expectedHTTP:    http.StatusBadRequest,
			},
			{
				name: "TC1.5 - Failure (Validation - Invalid Format)",
				reqBody: dto.ForgotPasswordRequestDTO{
					Identifier: "invalid-email-format@",
				},
				isValidationErr: true,
				expectedCode:    "VR-ID-002",
				expectedHTTP:    http.StatusBadRequest,
			},
			{
				name: "TC1.6 - Failure (Too Many Requests)",
				reqBody: dto.ForgotPasswordRequestDTO{
					Identifier: "spam@example.com",
				},
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO users (full_name, email, password_hash)
						VALUES ('Spam User', 'spam@example.com', 'dummyhash')
					`)
					require.NoError(t, err)

					_, err = pool.Exec(context.Background(), `
						INSERT INTO otp_requests (identifier, otp_hash, purpose, expires_at, created_at)
						VALUES ('spam@example.com', 'dummyhash', 'reset_password', NOW() + INTERVAL '5 minutes', NOW())
					`)
					require.NoError(t, err)
				},
				isValidationErr: false,
				expectedCode:    "AUTH_FORGOT_002",
				expectedMsg:     "TOO_MANY_REQUESTS",
				expectedHTTP:    http.StatusTooManyRequests,
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

				req, err := http.NewRequest("POST", "/v1/auth/forgot-password", bytes.NewBuffer(bodyBytes))
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
					if tc.expectedMsg != "" {
						assert.Equal(t, tc.expectedMsg, errMap["message"])
					}
				}
			})
		}
	})
}

// ====================================================================================
// TEST VERIFY RESET OTP
// ====================================================================================
func TestVerifyResetOtp(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success cases", func(t *testing.T) {
		tests := []struct {
			name      string
			reqBody   dto.VerifyResetOtpRequestDTO
			setupFunc func(t *testing.T, pool *pgxpool.Pool)
		}{
			{
				name: "TC2.1 - Success",
				reqBody: dto.VerifyResetOtpRequestDTO{
					ResetRequestId: validRequestID,
					Otp:            validOTP,
				},
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, created_at)
						VALUES ($1, 'test@example.com', $2, 'reset_password', NOW() + INTERVAL '5 minutes', 0, NOW())
					`, validRequestID, getValidOTPHash())
					require.NoError(t, err)
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

				req, err := http.NewRequest("POST", "/v1/auth/verify-reset-otp", bytes.NewBuffer(bodyBytes))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)

				var resp map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.NotEmpty(t, data["resetToken"])
				assert.NotZero(t, data["expiresIn"])

				// DB Verification
				var verifiedAt *time.Time
				err = pool.QueryRow(context.Background(), `SELECT verified_at FROM otp_requests WHERE id = $1`, validRequestID).Scan(&verifiedAt)
				require.NoError(t, err)
				assert.NotNil(t, verifiedAt, "Trường verified_at phải được cập nhật sau khi xác thực thành công")
			})
		}
	})

	t.Run("Failure cases", func(t *testing.T) {
		tests := []struct {
			name            string
			reqBody         dto.VerifyResetOtpRequestDTO
			setupFunc       func(t *testing.T, pool *pgxpool.Pool)
			isValidationErr bool
			expectedCode    string
			expectedMsg     string
			expectedHTTP    int
			checkDBFunc     func(t *testing.T, pool *pgxpool.Pool)
		}{
			{
				name: "TC2.2 - Failure (Invalid OTP)",
				reqBody: dto.VerifyResetOtpRequestDTO{
					ResetRequestId: validRequestID,
					Otp:            "654321", // OTP sai
				},
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, created_at)
						VALUES ($1, 'test@example.com', $2, 'reset_password', NOW() + INTERVAL '5 minutes', 0, NOW())
					`, validRequestID, getValidOTPHash())
					require.NoError(t, err)
				},
				isValidationErr: false,
				expectedCode:    "AUTH_OTP_001",
				expectedMsg:     "INVALID_OTP",
				expectedHTTP:    http.StatusBadRequest,
				checkDBFunc: func(t *testing.T, pool *pgxpool.Pool) {
					var attempts int
					err := pool.QueryRow(context.Background(), `SELECT attempts FROM otp_requests WHERE id = $1`, validRequestID).Scan(&attempts)
					require.NoError(t, err)
					assert.Equal(t, 1, attempts, "Trường attempts phải được cộng thêm 1 khi nhập sai OTP")
				},
			},
			{
				name: "TC2.3 - Failure (Validation - Bad OTP Format)",
				reqBody: dto.VerifyResetOtpRequestDTO{
					ResetRequestId: validRequestID,
					Otp:            "12AB", // Thiếu ký tự và không phải số
				},
				isValidationErr: true,
				expectedCode:    "VR-OTP-002",
				expectedHTTP:    http.StatusBadRequest,
			},
			{
				name: "TC2.4 - Failure (Expired OTP)",
				reqBody: dto.VerifyResetOtpRequestDTO{
					ResetRequestId: validRequestID,
					Otp:            validOTP,
				},
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, created_at)
						VALUES ($1, 'test@example.com', $2, 'reset_password', NOW() - INTERVAL '1 minutes', 0, NOW() - INTERVAL '6 minutes')
					`, validRequestID, getValidOTPHash())
					require.NoError(t, err)
				},
				isValidationErr: false,
				expectedCode:    "AUTH_OTP_002",
				expectedMsg:     "EXPIRED_OTP",
				expectedHTTP:    http.StatusBadRequest,
			},
			{
				name: "TC2.5 - Failure (Too Many Attempts)",
				reqBody: dto.VerifyResetOtpRequestDTO{
					ResetRequestId: validRequestID,
					Otp:            validOTP, // Nhập đúng nhưng đã quá lần thử
				},
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, created_at)
						VALUES ($1, 'test@example.com', $2, 'reset_password', NOW() + INTERVAL '5 minutes', 5, NOW())
					`, validRequestID, getValidOTPHash())
					require.NoError(t, err)
				},
				isValidationErr: false,
				expectedCode:    "AUTH_OTP_003",
				expectedMsg:     "TOO_MANY_ATTEMPTS",
				expectedHTTP:    http.StatusTooManyRequests,
			},
			{
				name: "TC2.6 - Failure (Already Verified)",
				reqBody: dto.VerifyResetOtpRequestDTO{
					ResetRequestId: validRequestID,
					Otp:            validOTP,
				},
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, verified_at, created_at)
						VALUES ($1, 'test@example.com', $2, 'reset_password', NOW() + INTERVAL '5 minutes', 0, NOW(), NOW())
					`, validRequestID, getValidOTPHash())
					require.NoError(t, err)
				},
				isValidationErr: false,
				expectedCode:    "AUTH_OTP_001",
				expectedMsg:     "INVALID_OTP",
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

				req, err := http.NewRequest("POST", "/v1/auth/verify-reset-otp", bytes.NewBuffer(bodyBytes))
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
					if tc.expectedMsg != "" {
						assert.Equal(t, tc.expectedMsg, errMap["message"])
					}
				}

				if tc.checkDBFunc != nil {
					tc.checkDBFunc(t, pool)
				}
			})
		}
	})
}

// ====================================================================================
// TEST RESET PASSWORD
// ====================================================================================
func TestResetPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Success cases", func(t *testing.T) {
		tests := []struct {
			name      string
			reqBody   dto.ResetPasswordRequestDTO
			setupFunc func(t *testing.T, pool *pgxpool.Pool)
		}{
			{
				name: "TC3.1 - Success (Reset Password)",
				reqBody: dto.ResetPasswordRequestDTO{
					ResetToken:      validResetToken,
					NewPassword:     "NewStrongPass123!",
					ConfirmPassword: "NewStrongPass123!",
				},
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO users (full_name, email, password_hash)
						VALUES ('Test User', 'reset@example.com', 'old_password_hash')
					`)
					require.NoError(t, err)

					_, err = pool.Exec(context.Background(), `
						INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, reset_token_hash, token_expires_at, created_at)
						VALUES ($1, 'reset@example.com', 'dummy_otp_hash', 'reset_password', NOW(), 0, $2, NOW() + INTERVAL '15 minutes', NOW())
					`, "00000000-0000-0000-0000-000000000001", getValidResetTokenHash())
					require.NoError(t, err)
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

				req, err := http.NewRequest("POST", "/v1/auth/reset-password", bytes.NewBuffer(bodyBytes))
				require.NoError(t, err)
				req.Header.Set("Content-Type", "application/json")

				w := httptest.NewRecorder()
				r.ServeHTTP(w, req)

				assert.Equal(t, http.StatusOK, w.Code)

				var resp map[string]interface{}
				err = json.Unmarshal(w.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.True(t, resp["success"].(bool))

				// DB Verification
				var newPassHash string
				err = pool.QueryRow(context.Background(), `SELECT password_hash FROM users WHERE email = 'reset@example.com'`).Scan(&newPassHash)
				require.NoError(t, err)
				assert.NotEqual(t, "old_password_hash", newPassHash, "Mật khẩu phải được băm lại và cập nhật")

				var tokenHash *string
				err = pool.QueryRow(context.Background(), `SELECT reset_token_hash FROM otp_requests WHERE id = '00000000-0000-0000-0000-000000000001'`).Scan(&tokenHash)
				require.NoError(t, err)
				assert.Nil(t, tokenHash, "reset_token_hash phải được set thành NULL để tránh dùng lại")
			})
		}
	})

	t.Run("Failure cases", func(t *testing.T) {
		tests := []struct {
			name            string
			reqBody         dto.ResetPasswordRequestDTO
			setupFunc       func(t *testing.T, pool *pgxpool.Pool)
			isValidationErr bool
			expectedCode    string
			expectedMsg     string
			expectedHTTP    int
		}{
			{
				name: "TC3.2 - Failure (Validation - Password Too Short)",
				reqBody: dto.ResetPasswordRequestDTO{
					ResetToken:      validResetToken,
					NewPassword:     "short",
					ConfirmPassword: "short",
				},
				isValidationErr: true,
				expectedCode:    "VR-NEW-PWD-002",
				expectedHTTP:    http.StatusBadRequest,
			},
			{
				name: "TC3.3 - Failure (Validation - Password Mismatch)",
				reqBody: dto.ResetPasswordRequestDTO{
					ResetToken:      validResetToken,
					NewPassword:     "NewStrongPass123!",
					ConfirmPassword: "MismatchPass123!",
				},
				isValidationErr: true,
				expectedCode:    "VR-CONFIRM-002",
				expectedHTTP:    http.StatusBadRequest,
			},
			{
				name: "TC3.4 - Failure (Token Expired)",
				reqBody: dto.ResetPasswordRequestDTO{
					ResetToken:      validResetToken,
					NewPassword:     "NewStrongPass123!",
					ConfirmPassword: "NewStrongPass123!",
				},
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, reset_token_hash, token_expires_at, created_at)
						VALUES ($1, 'reset@example.com', 'dummy_otp_hash', 'reset_password', NOW(), 0, $2, NOW() - INTERVAL '1 minutes', NOW() - INTERVAL '16 minutes')
					`, "00000000-0000-0000-0000-000000000002", getValidResetTokenHash())
					require.NoError(t, err)
				},
				isValidationErr: false,
				expectedCode:    "AUTH_RESET_003",
				expectedMsg:     "RESET_TOKEN_EXPIRED",
				expectedHTTP:    http.StatusGone,
			},
			{
				name: "TC3.5 - Failure (Token Not Found or Already Used)",
				reqBody: dto.ResetPasswordRequestDTO{
					ResetToken:      validResetToken, // Vẫn gửi token đúng
					NewPassword:     "NewStrongPass123!",
					ConfirmPassword: "NewStrongPass123!",
				},
				setupFunc: func(t *testing.T, pool *pgxpool.Pool) {
					_, err := pool.Exec(context.Background(), `
						INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, reset_token_hash, token_expires_at, created_at)
						VALUES ($1, 'reset@example.com', 'dummy_otp_hash', 'reset_password', NOW(), 0, NULL, NOW() + INTERVAL '15 minutes', NOW())
					`, "00000000-0000-0000-0000-000000000003")
					require.NoError(t, err)
				},
				isValidationErr: false,
				expectedCode:    "AUTH_RESET_003",
				expectedMsg:     "RESET_TOKEN_EXPIRED",
				expectedHTTP:    http.StatusGone,
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

				req, err := http.NewRequest("POST", "/v1/auth/reset-password", bytes.NewBuffer(bodyBytes))
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
					if tc.expectedMsg != "" {
						assert.Equal(t, tc.expectedMsg, errMap["message"])
					}
				}
			})
		}
	})
}
