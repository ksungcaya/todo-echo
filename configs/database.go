package configs

import "os"

type dbConnection struct {
	Dialect  string `json:"dialect"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// DatabaseConfig definition
type DatabaseConfig struct {
	MySQL dbConnection
}

// NewDatabaseConfig creates DatabaseConfig
func NewDatabaseConfig() DatabaseConfig {
	return DatabaseConfig{
		MySQL: dbConnection{
			Dialect:  "mysql",
			Host:     os.Getenv("DB_HOST"),
			Port:     GetEnvInt("DB_PORT", 3306),
			Database: os.Getenv("DB_NAME"),
			Username: os.Getenv("DB_USERNAME"),
			Password: os.Getenv("DB_PASSWORD"),
		},
	}
}
