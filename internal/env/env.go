package env

import (
	"os"
	"strconv"
)

// Config holds all environment-driven configuration for the application.
type Config struct {
	DBAddr        string
	DBMaxOpenConn int
	DBMaxIdleConn int
	DBMaxIdleTime string
	Port          string
}

// Load reads environment variables and returns a populated Config.
// Default values are used when a variable is not set.
func Load() Config {
	return Config{
		DBAddr:        getString("DB_ADDR", "postgres://postgres:postgres@localhost:5432/healmata?sslmode=disable"),
		DBMaxOpenConn: getInt("DB_MAX_OPEN_CONN", 25),
		DBMaxIdleConn: getInt("DB_MAX_IDLE_CONN", 25),
		DBMaxIdleTime: getString("DB_MAX_IDLE_TIME", "15m"),
		Port:          getString("PORT", "8080"),
	}
}

func getString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func getInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}
	return valAsInt
}
