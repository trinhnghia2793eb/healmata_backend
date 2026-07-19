package otp

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

type OTPService struct {
	
}

func NewOTPService() *OTPService {
	
	return &OTPService{}
}

// GenerateOTP generates a cryptographically secure numeric OTP of the specified length
func GenerateOTP(length int) (string, error) {
	const digits = "0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(digits))))
		if err != nil {
			return "", err
		}
		result[i] = digits[num.Int64()]
	}
	return string(result), nil
}

// HashOTP hashes the raw OTP using SHA-256
func HashOTP(otp string) string {
	hasher := sha256.New()
	hasher.Write([]byte(otp))
	return hex.EncodeToString(hasher.Sum(nil))
}