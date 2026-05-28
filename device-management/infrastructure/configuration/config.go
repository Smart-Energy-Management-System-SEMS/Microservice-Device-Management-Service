package configuration

import (
	"os"
	"strconv"
	"strings"
)

type AppConfig struct {
	Port                    string
	AppEnv                  string
	DatabaseURL             string
	DBDriver                string
	AutoMigrate             bool
	KafkaEnabled            bool
	KafkaBrokers            []string
	KafkaClientID           string
	KafkaConsumerGroup      string
	KafkaWriteTimeoutMS     int
	APIGatewayAllowedOrigin string
	CORSAllowedOrigins      []string
}

func LoadAppConfig() AppConfig {
	apiGatewayOrigin := getEnv("API_GATEWAY_ALLOWED_ORIGIN", "http://localhost:8080")
	corsAllowedOrigins := splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:4200,http://localhost:5173,http://localhost:8080"))
	corsAllowedOrigins = appendIfMissing(corsAllowedOrigins, apiGatewayOrigin)

	return AppConfig{
		Port:                    getEnv("PORT", "8083"),
		AppEnv:                  getEnv("APP_ENV", "local"),
		DatabaseURL:             getEnv("DATABASE_URL", ""),
		DBDriver:                getEnv("DB_DRIVER", "postgres"),
		AutoMigrate:             getBoolEnv("AUTO_MIGRATE", false),
		KafkaEnabled:            getBoolEnv("KAFKA_ENABLED", false),
		KafkaBrokers:            splitCSV(getEnv("KAFKA_BROKERS", "localhost:9092")),
		KafkaClientID:           getEnv("KAFKA_CLIENT_ID", "device-management-service"),
		KafkaConsumerGroup:      getEnv("KAFKA_CONSUMER_GROUP", "device-management-group"),
		KafkaWriteTimeoutMS:     getIntEnv("KAFKA_WRITE_TIMEOUT_MS", 2000),
		APIGatewayAllowedOrigin: apiGatewayOrigin,
		CORSAllowedOrigins:      corsAllowedOrigins,
	}
}

func getEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func getBoolEnv(key string, fallback bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if value == "" {
		return fallback
	}
	return value == "true" || value == "1" || value == "yes"
}

func getIntEnv(key string, fallback int) int {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		return fallback
	}
	return parsed
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

func appendIfMissing(values []string, value string) []string {
	for _, item := range values {
		if item == value {
			return values
		}
	}
	if value == "" {
		return values
	}
	return append(values, value)
}
