package configuration

import (
	"os"
	"strings"
)

type AppConfig struct {
	Port                    string
	AppEnv                  string
	DatabaseURL             string
	DBDriver                string
	KafkaBrokers            []string
	KafkaClientID           string
	KafkaConsumerGroup      string
	APIGatewayAllowedOrigin string
	CORSAllowedOrigins      []string
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		Port:                    getEnv("PORT", "8083"),
		AppEnv:                  getEnv("APP_ENV", "local"),
		DatabaseURL:             getEnv("DATABASE_URL", ""),
		DBDriver:                getEnv("DB_DRIVER", "postgres"),
		KafkaBrokers:            splitCSV(getEnv("KAFKA_BROKERS", "localhost:9092")),
		KafkaClientID:           getEnv("KAFKA_CLIENT_ID", "device-management-service"),
		KafkaConsumerGroup:      getEnv("KAFKA_CONSUMER_GROUP", "device-management-group"),
		APIGatewayAllowedOrigin: getEnv("API_GATEWAY_ALLOWED_ORIGIN", "http://localhost:8080"),
		CORSAllowedOrigins:      splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:4200,http://localhost:5173,http://localhost:8080")),
	}
}

func getEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func splitCSV(value string) []string {
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
