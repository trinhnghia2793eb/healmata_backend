package bootstrap

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	AppName string
	AppPort string
	AppEnv  string

	GinMode string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	SMTPHost        string
	SMTPPort        int
	SMTPUser        string
	SMTPPassword    string
	MailFromAddress string
	MailFromName    string
}

func LoadEnv() (*Config, error) {
	// Không lỗi nếu không có file .env
	_ = godotenv.Load()

	// string --> int for smtpPort
	smtpPort, err := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	if err != nil {
		smtpPort = 587 // fallback
	}

	cfg := &Config{
		AppName:    getEnv("APP_NAME", "Go Application"),
		AppPort:    getEnv("APP_PORT", "8080"),
		AppEnv:     getEnv("APP_ENV", "development"),
		GinMode:    getEnv("GIN_MODE", "debug"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		SMTPHost:        getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:        smtpPort,
		SMTPUser:        os.Getenv("SMTP_USER"),
		SMTPPassword:    os.Getenv("SMTP_PASSWORD"),
		MailFromAddress: getEnv("MAIL_FROM_ADDRESS", "no-reply@myproject.com"),
		MailFromName:    getEnv("MAIL_FROM_NAME", "My Project"),
	}

	if cfg.DBUser == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}

	if cfg.DBName == "" {
		return nil, fmt.Errorf("DB_NAME is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
