package testutils

import "fmt"

type MockEmailSender struct{
	SendCalled bool
}

func NewMockEmailSender() *MockEmailSender {
	return &MockEmailSender{}
}

func (m *MockEmailSender) SendOTP(to string, otpCode string) error {
	// Chỉ in ra log để biết là hàm đã được gọi, không gửi email thực
	fmt.Printf("[Test] Mock Email: Gửi OTP %s tới %s\n", otpCode, to)
	m.SendCalled = true
	return nil
}
