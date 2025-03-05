package config

import "os"

// Config хранит конфигурационные параметры для API Gateway
type Config struct {
	HTTPPort           string // Порт для HTTP сервера, например ":8080"
	UserServiceAddress string // Адрес gRPC-сервиса пользователей, например "localhost:50051"
	PostServiceAddress string // Адрес gRPC-сервиса постов, например "localhost:50052"
	GatewayAddress     string // Адрес gRPC-сервиса API Gateway, например "localhost:50053"
}

// LoadConfig загружает конфигурацию из переменных окружения или возвращает значения по умолчанию
func LoadConfig() (*Config, error) {
	cfg := &Config{
		HTTPPort:           getEnv("HTTP_PORT", ":8080"),
		UserServiceAddress: getEnv("USER_SERVICE_ADDRESS", "localhost:50051"),
		PostServiceAddress: getEnv("POST_SERVICE_ADDRESS", "localhost:50052"),
		GatewayAddress:     getEnv("GATEWAY_ADDRESS", "localhost:50053"),
	}
	return cfg, nil
}

// getEnv возвращает значение переменной окружения или значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
