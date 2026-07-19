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

	"healmata_backend/internal/app/router"
	"healmata_backend/internal/auth/dto"
	"healmata_backend/internal/db/testhelper"
)

func TestForgotPassword(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// for validate OTP API
	validOTP := "123456"
	hash := sha256.Sum256([]byte(validOTP))
	validOTPHash := hex.EncodeToString(hash[:])
	validRequestID := "123e4567-e89b-12d3-a456-426614174000"

	// for reset token API
	validResetToken := "valid-reset-token-xyz"
	resetHash := sha256.Sum256([]byte(validResetToken))
	validResetTokenHash := hex.EncodeToString(resetHash[:])

	// testCase struct
	type testCase struct {
		name string
		method string
		url string
		setupDB func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) // Chuẩn bị dữ liệu mẫu
		requestBody interface{}                                             // DTO hoặc raw map để test lỗi format
		expectedStatus int
		checkResponse func(t *testing.T, resp map[string]interface{})       // Check returned JSON
		checkDB func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) // Check DB change
	}

	// ==============================================================================================================================================
	// test cases 
	tests := []testCase{
		// ======================================================================= 
		// POST /v1/auth/forgot-password
		{
			name: "TC1.1 - Success (Email)",
			method: "POST",
			url: "/v1/auth/forgot-password",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				_, err := pool.Exec(ctx, `
					INSERT INTO users (full_name, email, password_hash)
					VALUES ('User Email', 'valid@example.com', 'dummyhash')
				`)
				require.NoError(t, err)
			},
			requestBody: dto.ForgotPasswordRequestDTO{
				Identifier: "valid@example.com",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				assert.Equal(t, "OTP_SENT", resp["message"])

				data := resp["data"].(map[string]interface{})
				assert.NotEmpty(t, data["resetRequestId"])
				assert.Equal(t, float64(6), data["otpLength"])
				assert.Equal(t, float64(300), data["expiresIn"])
				assert.Equal(t, float64(60), data["resendAfter"])
			},
			checkDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				var count int
				err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM otp_requests WHERE identifier = 'valid@example.com' AND purpose = 'reset_password'`).Scan(&count)
				require.NoError(t, err)
				assert.Equal(t, 1, count, "Nên có 1 record OTP được tạo trong DB")
			},
		},
		// {
		// 	name: "TC1.2 - Success (Phone)",
		// 	method: "POST",
		// 	url: "/v1/auth/forgot-password",
		// 	setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
		// 		_, err := pool.Exec(ctx, `
		// 			INSERT INTO users (full_name, phone, password_hash)
		// 			VALUES ('User Phone', '+84987654321', 'dummyhash')
		// 		`)
		// 		require.NoError(t, err)
		// 	},
		// 	requestBody: dto.ForgotPasswordRequestDTO{
		// 		Identifier: "+84987654321",
		// 	},
		// 	expectedStatus: http.StatusOK,
		// 	checkResponse: func(t *testing.T, resp map[string]interface{}) {
		// 		assert.True(t, resp["success"].(bool))
		// 		assert.Equal(t, "OTP_SENT", resp["message"])
		// 		assert.NotEmpty(t, resp["data"])
		// 	},
		// 	checkDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
		// 		var count int
		// 		err := pool.QueryRow(ctx, `SELECT COUNT(*) FROM otp_requests WHERE identifier = '+84987654321' AND purpose = 'reset_password'`).Scan(&count)
		// 		require.NoError(t, err)
		// 		assert.Equal(t, 1, count)
		// 	},
		// },
		{
			name: "TC1.3 - Failure (User Not Found)",
			method: "POST",
			url: "/v1/auth/forgot-password",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				// Không insert user nào
			},
			requestBody: dto.ForgotPasswordRequestDTO{
				Identifier: "notfound@example.com",
			},
			expectedStatus: http.StatusNotFound,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				errMap := resp["error"].(map[string]interface{})
				assert.Equal(t, "AUTH_FORGOT_001", errMap["code"])
			},
			checkDB: nil,
		},
		{
			name: "TC1.4 - Failure (Validation - Empty Identifier)",
			method: "POST",
			url: "/v1/auth/forgot-password",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {},
			requestBody: map[string]string{ // use 
				"identifier": "   ",
			},
			expectedStatus: http.StatusUnprocessableEntity, // customErrors trả về HTTP 422
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				// Middleware bắt lỗi cấu trúc/trường required
				errMap := resp["error"].(map[string]interface{})
				assert.NotEmpty(t, errMap["code"])
			},
			checkDB: nil,
		},
		{
			name: "TC1.5 - Failure (Validation - Invalid Format)",
			method: "POST",
			url: "/v1/auth/forgot-password",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {},
			requestBody: dto.ForgotPasswordRequestDTO{
				Identifier: "invalid-email-format@",
			},
			expectedStatus: http.StatusUnprocessableEntity, // customErrors trả về HTTP 422
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				errMap := resp["error"].(map[string]interface{})
				assert.Equal(t, "AUTH_REG_006", errMap["code"]) // ErrInvalidEmail
			},
			checkDB: nil,
		},
		{
			name: "TC1.6 - Failure (Too Many Requests)",
			method: "POST",
			url: "/v1/auth/forgot-password",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				// 1. Tạo user
				_, err := pool.Exec(ctx, `
					INSERT INTO users (full_name, email, password_hash)
					VALUES ('Spam User', 'spam@example.com', 'dummyhash')
				`)
				require.NoError(t, err)

				// 2. Tạo 1 request OTP vừa mới sinh ra (dưới 60s)
				_, err = pool.Exec(ctx, `
					INSERT INTO otp_requests (identifier, otp_hash, purpose, expires_at, created_at)
					VALUES ('spam@example.com', 'dummyhash', 'reset_password', NOW() + INTERVAL '5 minutes', NOW())
				`)
				require.NoError(t, err)
			},
			requestBody: dto.ForgotPasswordRequestDTO{
				Identifier: "spam@example.com",
			},
			expectedStatus: http.StatusTooManyRequests,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				errMap := resp["error"].(map[string]interface{})
				assert.Equal(t, "AUTH_FORGOT_002", errMap["code"]) // TOO_MANY_REQUESTS
			},
			checkDB: nil,
		},

		// ======================================================================= 
		// POST /v1/auth/verify-reset-otp
		{
			name: "TC2.1 - Success",
			method: "POST",
			url: "/v1/auth/verify-reset-otp",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				_, err := pool.Exec(ctx, `
					INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, created_at)
					VALUES ($1, 'test@example.com', $2, 'reset_password', NOW() + INTERVAL '5 minutes', 0, NOW())
				`, validRequestID, validOTPHash)
				require.NoError(t, err)
			},
			requestBody: dto.VerifyResetOtpRequestDTO{
				ResetRequestId: validRequestID,
				Otp:            validOTP,
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
				data := resp["data"].(map[string]interface{})
				assert.NotEmpty(t, data["resetToken"])
				assert.NotZero(t, data["expiresIn"])
			},
			checkDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				var verifiedAt *time.Time
				err := pool.QueryRow(ctx, `SELECT verified_at FROM otp_requests WHERE id = $1`, validRequestID).Scan(&verifiedAt)
				require.NoError(t, err)
				assert.NotNil(t, verifiedAt, "Trường verified_at phải được cập nhật sau khi xác thực thành công")
			},
		},
		{
			name: "TC2.2 - Failure (Invalid OTP)",
			method: "POST",
			url: "/v1/auth/verify-reset-otp",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				_, err := pool.Exec(ctx, `
					INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, created_at)
					VALUES ($1, 'test@example.com', $2, 'reset_password', NOW() + INTERVAL '5 minutes', 0, NOW())
				`, validRequestID, validOTPHash)
				require.NoError(t, err)
			},
			requestBody: dto.VerifyResetOtpRequestDTO{
				ResetRequestId: validRequestID,
				Otp:            "654321", // OTP sai
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				errMap := resp["error"].(map[string]interface{})
				assert.Equal(t, "AUTH_OTP_001", errMap["code"])
			},
			checkDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				var attempts int
				err := pool.QueryRow(ctx, `SELECT attempts FROM otp_requests WHERE id = $1`, validRequestID).Scan(&attempts)
				require.NoError(t, err)
				assert.Equal(t, 1, attempts, "Trường attempts phải được cộng thêm 1 khi nhập sai OTP")
			},
		},
		{
			name: "TC2.3 - Failure (Validation - Bad OTP Format)",
			method: "POST",
			url: "/v1/auth/verify-reset-otp",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {},
			requestBody: map[string]interface{}{
				"resetRequestId": validRequestID,
				"otp":            "12AB", // Thiếu ký tự và không phải số
			},
			expectedStatus: http.StatusBadRequest, // Middleware chặn
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				errMap := resp["error"].(map[string]interface{})
				assert.NotEmpty(t, errMap["code"])
			},
			checkDB: nil,
		},
		{
			name: "TC2.4 - Failure (Expired OTP)",
			method: "POST",
			url: "/v1/auth/verify-reset-otp",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				_, err := pool.Exec(ctx, `
					INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, created_at)
					VALUES ($1, 'test@example.com', $2, 'reset_password', NOW() - INTERVAL '1 minutes', 0, NOW() - INTERVAL '6 minutes')
				`, validRequestID, validOTPHash) // expires_at nằm trong quá khứ
				require.NoError(t, err)
			},
			requestBody: dto.VerifyResetOtpRequestDTO{
				ResetRequestId: validRequestID,
				Otp:            validOTP,
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				errMap := resp["error"].(map[string]interface{})
				assert.Equal(t, "AUTH_OTP_002", errMap["code"]) // EXPIRED_OTP
			},
			checkDB: nil,
		},
		{
			name: "TC2.5 - Failure (Too Many Attempts)",
			method: "POST",
			url: "/v1/auth/verify-reset-otp",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				_, err := pool.Exec(ctx, `
					INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, created_at)
					VALUES ($1, 'test@example.com', $2, 'reset_password', NOW() + INTERVAL '5 minutes', 5, NOW())
				`, validRequestID, validOTPHash) // attempts đã đạt giới hạn (>=5)
				require.NoError(t, err)
			},
			requestBody: dto.VerifyResetOtpRequestDTO{
				ResetRequestId: validRequestID,
				Otp:            validOTP, // Dù OTP đúng cũng phải chặn
			},
			expectedStatus: http.StatusTooManyRequests, // 429
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				errMap := resp["error"].(map[string]interface{})
				assert.Equal(t, "AUTH_OTP_003", errMap["code"]) // MAX_ATTEMPTS_EXCEEDED
			},
			checkDB: nil,
		},
		{
			name: "TC2.6 - Failure (Already Verified)",
			method: "POST",
			url: "/v1/auth/verify-reset-otp",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				_, err := pool.Exec(ctx, `
					INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, verified_at, created_at)
					VALUES ($1, 'test@example.com', $2, 'reset_password', NOW() + INTERVAL '5 minutes', 0, NOW(), NOW())
				`, validRequestID, validOTPHash) // verified_at khác NULL
				require.NoError(t, err)
			},
			requestBody: dto.VerifyResetOtpRequestDTO{
				ResetRequestId: validRequestID,
				Otp:            validOTP,
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				errMap := resp["error"].(map[string]interface{})
				assert.Equal(t, "AUTH_OTP_001", errMap["code"]) // INVALID_OTP hoặc lỗi tương đương tùy logic của bạn
			},
			checkDB: nil,
		},

		// ======================================================================= 
		// POST /v1/auth/reset-password
		{
			name:   "TC3.1 - Success (Reset Password)",
			method: "POST",
			url:    "/v1/auth/reset-password",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				// 1. Tạo user cũ
				_, err := pool.Exec(ctx, `
					INSERT INTO users (full_name, email, password_hash)
					VALUES ('Test User', 'reset@example.com', 'old_password_hash')
				`)
				require.NoError(t, err)

				// 2. Tạo record OTP hợp lệ đã cấp token (còn hạn)
				_, err = pool.Exec(ctx, `
					INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, reset_token_hash, token_expires_at, created_at)
					VALUES ($1, 'reset@example.com', 'dummy_otp_hash', 'reset_password', NOW(), 0, $2, NOW() + INTERVAL '15 minutes', NOW())
				`, "00000000-0000-0000-0000-000000000001", validResetTokenHash)
				require.NoError(t, err)
			},
			requestBody: dto.ResetPasswordRequestDTO{
				ResetToken:      validResetToken,
				NewPassword:     "NewStrongPass123!",
				ConfirmPassword: "NewStrongPass123!",
			},
			expectedStatus: http.StatusOK,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.True(t, resp["success"].(bool))
			},
			checkDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				// 1. Kiểm tra password của user đã được đổi (hash khác với ban đầu)
				var newPassHash string
				err := pool.QueryRow(ctx, `SELECT password_hash FROM users WHERE email = 'reset@example.com'`).Scan(&newPassHash)
				require.NoError(t, err)
				assert.NotEqual(t, "old_password_hash", newPassHash, "Mật khẩu phải được băm lại và cập nhật")

				// 2. Kiểm tra token đã bị vô hiệu hóa (set về NULL)
				var tokenHash *string
				err = pool.QueryRow(ctx, `SELECT reset_token_hash FROM otp_requests WHERE id = '00000000-0000-0000-0000-000000000001'`).Scan(&tokenHash)
				require.NoError(t, err)
				assert.Nil(t, tokenHash, "reset_token_hash phải được set thành NULL để tránh dùng lại")
			},
		},
		{
			name:   "TC3.2 - Failure (Validation - Password Too Short)",
			method: "POST",
			url:    "/v1/auth/reset-password",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {},
			requestBody: dto.ResetPasswordRequestDTO{
				ResetToken:      validResetToken,
				NewPassword:     "short",
				ConfirmPassword: "short",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				errMap := resp["error"].(map[string]interface{})
				assert.Equal(t, "AUTH_RESET_001", errMap["code"])
			},
			checkDB: nil,
		},
		{
			name:   "TC3.3 - Failure (Validation - Password Mismatch)",
			method: "POST",
			url:    "/v1/auth/reset-password",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {},
			requestBody: dto.ResetPasswordRequestDTO{
				ResetToken:      validResetToken,
				NewPassword:     "NewStrongPass123!",
				ConfirmPassword: "MismatchPass123!",
			},
			expectedStatus: http.StatusBadRequest,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				errMap := resp["error"].(map[string]interface{})
				assert.Equal(t, "AUTH_RESET_002", errMap["code"]) // Mismatch error
			},
			checkDB: nil,
		},
		{
			name:   "TC3.4 - Failure (Token Expired)",
			method: "POST",
			url:    "/v1/auth/reset-password",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				// Cấp token nhưng token_expires_at nằm ở quá khứ
				_, err := pool.Exec(ctx, `
					INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, reset_token_hash, token_expires_at, created_at)
					VALUES ($1, 'reset@example.com', 'dummy_otp_hash', 'reset_password', NOW(), 0, $2, NOW() - INTERVAL '1 minutes', NOW() - INTERVAL '16 minutes')
				`, "00000000-0000-0000-0000-000000000002", validResetTokenHash)
				require.NoError(t, err)
			},
			requestBody: dto.ResetPasswordRequestDTO{
				ResetToken:      validResetToken,
				NewPassword:     "NewStrongPass123!",
				ConfirmPassword: "NewStrongPass123!",
			},
			expectedStatus: http.StatusGone, // Hoặc http.StatusBadRequest tùy thuộc vào định nghĩa của mã AUTH_RESET_003
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				errMap := resp["error"].(map[string]interface{})
				assert.Equal(t, "AUTH_RESET_003", errMap["code"]) // TOKEN_INVALID_OR_EXPIRED
			},
			checkDB: nil,
		},
		{
			name:   "TC3.5 - Failure (Token Not Found or Already Used)",
			method: "POST",
			url:    "/v1/auth/reset-password",
			setupDB: func(t *testing.T, pool *pgxpool.Pool, ctx context.Context) {
				// Mô phỏng token đã sử dụng (reset_token_hash đã bị set thành NULL)
				_, err := pool.Exec(ctx, `
					INSERT INTO otp_requests (id, identifier, otp_hash, purpose, expires_at, attempts, reset_token_hash, token_expires_at, created_at)
					VALUES ($1, 'reset@example.com', 'dummy_otp_hash', 'reset_password', NOW(), 0, NULL, NOW() + INTERVAL '15 minutes', NOW())
				`, "00000000-0000-0000-0000-000000000003")
				require.NoError(t, err)
			},
			requestBody: dto.ResetPasswordRequestDTO{
				ResetToken:      validResetToken, // Vẫn gửi token đúng
				NewPassword:     "NewStrongPass123!",
				ConfirmPassword: "NewStrongPass123!",
			},
			expectedStatus: http.StatusGone,
			checkResponse: func(t *testing.T, resp map[string]interface{}) {
				assert.False(t, resp["success"].(bool))
				errMap := resp["error"].(map[string]interface{})
				assert.Equal(t, "AUTH_RESET_003", errMap["code"]) 
			},
			checkDB: nil,
		},

	}

	// ==============================================================================================================================================
	// run test
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Thiết lập DB độc lập cho mỗi test case thông qua testhelper
			pool := testhelper.SetupTestDB(t)
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Chạy logic setup dữ liệu mẫu
			if tc.setupDB != nil {
				tc.setupDB(t, pool, ctx)
			}

			// Cấu hình Router và Mock
			r := gin.New()
			emailSender := testhelper.NewMockEmailSender()
			router.RegisterRoutes(r, pool, emailSender)

			// Chuẩn bị Request Body
			bodyBytes, err := json.Marshal(tc.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest(tc.method, tc.url, bytes.NewBuffer(bodyBytes))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Khởi chạy Request
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			// Kiểm tra HTTP Status
			assert.Equal(t, tc.expectedStatus, w.Code)

			// Decode và kiểm tra JSON Response
			var resp map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &resp)
			require.NoError(t, err)

			if tc.checkResponse != nil {
				tc.checkResponse(t, resp)
			}

			// Kiểm tra tác động xuống DB sau khi API chạy xong
			if tc.checkDB != nil {
				tc.checkDB(t, pool, ctx)
			}
		})
	}
}