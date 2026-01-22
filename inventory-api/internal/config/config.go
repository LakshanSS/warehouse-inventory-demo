package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port        int
	Environment string
	DatabaseURL string
	Hostname    string
}

func LoadConfig() (*Config, error) {
	port := 9090
	if p := os.Getenv("PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			port = parsed
		}
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	hostname := os.Getenv("HOSTNAME")
	if hostname == "" {
		hostname = "localhost"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Default connection string for local development
		server := getEnvOrDefault("DB_SERVER", "localhost")
		dbPort := getEnvOrDefault("DB_PORT", "1433")
		user := getEnvOrDefault("DB_USER", "sa")
		password := getEnvOrDefault("DB_PASSWORD", "YourStrong@Passw0rd")
		database := getEnvOrDefault("DB_NAME", "InventoryDB")

		dbURL = fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s&encrypt=true&trustServerCertificate=false",
			user, password, server, dbPort, database)
	}

	return &Config{
		Port:        port,
		Environment: env,
		DatabaseURL: dbURL,
		Hostname:    hostname,
	}, nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
