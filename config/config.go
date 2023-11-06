package config

import (
	"os"
)

type Config struct {
	// Configurações de banco de dados
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Configurações de servidor
	ServerPort string
}

func NewConfig() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "username"),
		DBPassword: getEnv("DB_PASSWORD", "password"),
		DBName:     getEnv("DB_NAME", "database_name"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
