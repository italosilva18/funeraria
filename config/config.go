package config

import (
	"os"
)

// Config armazena as configurações da aplicação
type Config struct {
	// Adicione aqui as configurações necessárias
	MongoDBURI string
	Port       string
	// Outras configurações...
}

// LoadConfig carrega as configurações da aplicação a partir de variáveis de ambiente
func LoadConfig() *Config {
	return &Config{
		MongoDBURI: getEnv("MONGODB_URI", "localhost:27017"),
		Port:       getEnv("PORT", "8080"),
		// Configure outras variáveis de ambiente conforme necessário
	}
}

// getEnv retorna o valor de uma variável de ambiente ou um valor padrão caso não esteja definido
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
