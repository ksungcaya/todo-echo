package configs

import (
	"os"
	"strconv"
)

// AppConfig definition
type AppConfig struct {
	Port     int            `json:"port"`
	Env      string         `json:"env"`
	Database DatabaseConfig `json:"database"`
}

// IsProd determines if current app env is in production
func (c AppConfig) IsProd() bool {
	return c.Env == "production"
}

// New creates new AppConfig
func New() AppConfig {
	return AppConfig{
		Port:     GetEnvInt("APP_PORT", 5050),
		Env:      GetEnv("APP_ENV", "development"),
		Database: NewDatabaseConfig(),
	}
}

// GetEnv with a fallback
func GetEnv(key string, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetEnvInt gets env and parses it to int
func GetEnvInt(key string, fallback int) int {
	env, ok := os.LookupEnv(key)
	if ok == false {
		return fallback
	}

	v, err := strconv.Atoi(env)
	if err != nil {
		return fallback
	}

	return v
}

// GetEnvBool gets env and parses it to boolean
func GetEnvBool(key string, fallback bool) bool {
	env, ok := os.LookupEnv(key)
	if ok == false {
		return fallback
	}

	v, err := strconv.ParseBool(env)
	if err != nil {
		return fallback
	}

	return v
}
