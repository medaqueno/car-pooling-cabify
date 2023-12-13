package config

import "os"

type AppConfig struct {
	ServerPort string
	LogLevel   string
}

// LoadConfig loads the application configuration from environment variables
func LoadConfig() (*AppConfig, error) {
	config := AppConfig{
		ServerPort: getEnv("SERVER_PORT", "9091"),
		LogLevel:   getEnv("LOG_LEVEL", "info"),
	}

	return &config, nil
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
