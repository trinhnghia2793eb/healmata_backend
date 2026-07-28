package config

import (
	"path/filepath"
	"time"

	"github.com/spf13/viper"
)

// config struct
type Config struct {
	AppName string `mapstructure:"APP_NAME"`
	AppPort string `mapstructure:"APP_PORT"`
	AppEnv  string `mapstructure:"APP_ENV"`
	
	GinMode            string   `mapstructure:"GIN_MODE"`
	CorsAllowedOrigins []string `mapstructure:"CORS_ALLOWED_ORIGINS"`

	DBHost     string `mapstructure:"DB_HOST"`
	DBPort     string `mapstructure:"DB_PORT"`
	DBUser     string `mapstructure:"DB_USER"`
	DBPassword string `mapstructure:"DB_PASSWORD"`
	DBName     string `mapstructure:"DB_NAME"`
	DBSSLMode  string `mapstructure:"DB_SSLMODE"`
	
	JWTSecret        string        `mapstructure:"JWT_SECRET"`
	JWTAccessExpiry  time.Duration `mapstructure:"JWT_ACCESS_EXPIRY"`
	JWTRefreshExpiry time.Duration `mapstructure:"JWT_REFRESH_EXPIRY"`

	SMTPHost        string `mapstructure:"SMTP_HOST"`
	SMTPPort        int    `mapstructure:"SMTP_PORT"`
	SMTPUser        string `mapstructure:"SMTP_USER"`
	SMTPPassword    string `mapstructure:"SMTP_PASSWORD"`
	MailFromAddress string `mapstructure:"MAIL_FROM_ADDRESS"`
	MailFromName    string `mapstructure:"MAIL_FROM_NAME"`
}

// load config
func LoadConfig(path string) (*Config, error) {
	viper.SetDefault("APP_NAME", "Go Application")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("GIN_MODE", "debug")
	viper.SetDefault("CORS_ALLOWED_ORIGINS", []string{"*"})
	viper.SetDefault("JWT_SECRET", "my-project-secret-key-change-me-in-production")
	viper.SetDefault("JWT_ACCESS_EXPIRY", "1h")
	viper.SetDefault("JWT_REFRESH_EXPIRY", "720h")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("SMTP_PORT", 587)

	// read .env file
	envPath := filepath.Join(path, ".env")
	viper.SetConfigFile(envPath)

	// allows overriding operating system environment variables (higher priority than the .env file).
	viper.AutomaticEnv()

	// read the configuration file (proceed even if the .env file is missing, as execution is possible using OS environment variables)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}
	
	// bring data into struct
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
