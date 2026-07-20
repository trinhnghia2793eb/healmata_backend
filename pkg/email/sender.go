package email

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

//go:embed templates/otp.html
var templateFS embed.FS

type EmailSender interface {
	SendOTP(to string, otpCode string) error
}

type gomailSender struct {
	host            string
	port            int
	user            string
	password        string
	mailFromAddress string
	mailFromName    string
}

// OTP data --> inject into HTML Template
type otpTemplateData struct {
	OTPCode      string
	MailFromName string
}

// constructor
func NewEmailSender(host string, port int, user, password, fromAddress, fromName string) EmailSender {
	return &gomailSender{
		host:            host,
		port:            port,
		user:            user,
		password:        password,
		mailFromAddress: fromAddress,
		mailFromName:    fromName,
	}
}

// SendOTP
func (s *gomailSender) SendOTP(to string, otpCode string) error {
	// parse file template from embed folder
	tmpl, err := template.ParseFS(templateFS, "templates/otp.html")
	if err != nil {
		return fmt.Errorf("cannot parse email template: %w", err)
	}

	// prepare data --> execute email template into a buffer
	data := otpTemplateData{
		OTPCode:      otpCode,
		MailFromName: s.mailFromName,
	}
	var bodyBuffer bytes.Buffer
	if err := tmpl.Execute(&bodyBuffer, data); err != nil {
		return fmt.Errorf("cannot execute email template: %w", err)
	}

	// build image content
	m := gomail.NewMessage()
	m.SetAddressHeader("From", s.mailFromAddress, s.mailFromName)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Healmata - Xác thực hệ thống")

	// set body format: text/html --> assign html content with injected data
	m.SetBody("text/html", bodyBuffer.String())

	// initialize dialer to connect to SMTP server --> send mail
	d := gomail.NewDialer(s.host, s.port, s.user, s.password)
	return d.DialAndSend(m)
}
