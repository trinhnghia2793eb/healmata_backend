package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"userId"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	secretKey     string
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

func NewJWTManager(secretKey string, accessExpiry, refreshExpiry time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:     secretKey,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateToken creates a signed token (used for both Access and Refresh tokens)
func (jm *JWTManager) GenerateToken(userID string, duration time.Duration) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jm.secretKey))
}

// Generate Access Token and Refresh Token
func (jm *JWTManager) GenerateAccessAndRefreshToken(userID string) (string, string, int64, error) {
	accessToken, err := jm.GenerateToken(userID, jm.accessExpiry)
	if err != nil {
		return "", "", 0, err
	}
	refreshToken, err := jm.GenerateToken(userID, jm.refreshExpiry)
	if err != nil {
		return "", "", 0, err
	}
	return accessToken, refreshToken, int64(jm.accessExpiry.Seconds()), nil
}

// VerifyToken verifies a token and returns the claims if the token is valid
func (jm *JWTManager) VerifyToken(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jm.secretKey), nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token claims")
}
