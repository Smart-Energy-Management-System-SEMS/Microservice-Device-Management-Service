package configuration

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type AppConfig struct {
	ServiceName      string
	Port             string
	ConfigServiceURL string
	DatabaseURL      string
}

type RuntimeConfig struct {
	AppEnv                  string
	DBDriver                string
	AutoMigrate             bool
	KafkaEnabled            bool
	KafkaBrokers            []string
	KafkaClientID           string
	KafkaConsumerGroup      string
	KafkaWriteTimeoutMS     int
	APIGatewayAllowedOrigin string
	CORSAllowedOrigins      []string
	KafkaTopics             KafkaTopics
}

type KafkaTopics struct {
	DeviceRegistered           string
	DeviceStatusUpdated        string
	DeviceLinked               string
	DeviceUnlinked             string
	DeviceConfigurationUpdated string
	DeviceEventRecorded        string
}

type serviceConfigResponse struct {
	AppEnv                  string      `json:"appEnv"`
	DBDriver                string      `json:"dbDriver"`
	AutoMigrate             *bool       `json:"autoMigrate"`
	KafkaEnabled            *bool       `json:"kafkaEnabled"`
	KafkaBrokers            []string    `json:"kafkaBrokers"`
	KafkaClientID           string      `json:"kafkaClientId"`
	KafkaConsumerGroup      string      `json:"kafkaConsumerGroup"`
	KafkaWriteTimeoutMS     int         `json:"kafkaWriteTimeoutMs"`
	APIGatewayAllowedOrigin string      `json:"apiGatewayAllowedOrigin"`
	CORSAllowedOrigins      []string    `json:"corsAllowedOrigins"`
	KafkaTopics             KafkaTopics `json:"kafkaTopics"`
}

func LoadAppConfig() AppConfig {
	return AppConfig{
		ServiceName:      getEnv("SERVICE_NAME", "device-management-service"),
		Port:             getEnv("PORT", "8083"),
		ConfigServiceURL: strings.TrimRight(getEnv("CONFIG_SERVICE_URL", "http://localhost:8081"), "/"),
		DatabaseURL:      getEnv("DATABASE_URL", ""),
	}
}

func ResolveRuntimeConfig(ctx context.Context, appConfig AppConfig) RuntimeConfig {
	localFallback := defaultRuntimeConfig(appConfig.ServiceName)
	serviceConfig, err := fetchServiceConfig(ctx, appConfig.ConfigServiceURL, appConfig.ServiceName)
	if err != nil {
		return localFallback
	}
	return mergeRuntimeConfig(localFallback, serviceConfig)
}

func fetchServiceConfig(ctx context.Context, baseURL string, serviceName string) (*serviceConfigResponse, error) {
	if strings.TrimSpace(baseURL) == "" {
		return nil, fmt.Errorf("config service URL is empty")
	}
	url := fmt.Sprintf("%s/api/v1/config/%s", baseURL, serviceName)
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{Timeout: 2 * time.Second}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("config service returned status %d", response.StatusCode)
	}

	var payload serviceConfigResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

func defaultRuntimeConfig(serviceName string) RuntimeConfig {
	apiGatewayOrigin := getEnv("API_GATEWAY_ALLOWED_ORIGIN", "http://localhost:8080")
	corsAllowedOrigins := splitCSV(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:4200,http://localhost:5173,http://localhost:8080"))
	corsAllowedOrigins = appendIfMissing(corsAllowedOrigins, apiGatewayOrigin)

	return RuntimeConfig{
		AppEnv:                  getEnv("APP_ENV", "local"),
		DBDriver:                getEnv("DB_DRIVER", "postgres"),
		AutoMigrate:             getBoolEnv("AUTO_MIGRATE", false),
		KafkaEnabled:            getBoolEnv("KAFKA_ENABLED", false),
		KafkaBrokers:            splitCSV(getEnv("KAFKA_BROKERS", "localhost:9092")),
		KafkaClientID:           getEnv("KAFKA_CLIENT_ID", serviceName),
		KafkaConsumerGroup:      getEnv("KAFKA_CONSUMER_GROUP", "device-management-group"),
		KafkaWriteTimeoutMS:     getIntEnv("KAFKA_WRITE_TIMEOUT_MS", 2000),
		APIGatewayAllowedOrigin: apiGatewayOrigin,
		CORSAllowedOrigins:      corsAllowedOrigins,
		KafkaTopics: KafkaTopics{
			DeviceRegistered:           getEnv("TOPIC_DEVICE_REGISTERED", "device.registered"),
			DeviceStatusUpdated:        getEnv("TOPIC_DEVICE_STATUS_UPDATED", "device.status.updated"),
			DeviceLinked:               getEnv("TOPIC_DEVICE_LINKED", "device.linked"),
			DeviceUnlinked:             getEnv("TOPIC_DEVICE_UNLINKED", "device.unlinked"),
			DeviceConfigurationUpdated: getEnv("TOPIC_DEVICE_CONFIGURATION_UPDATED", "device.configuration.updated"),
			DeviceEventRecorded:        getEnv("TOPIC_DEVICE_EVENT_RECORDED", "device.event.recorded"),
		},
	}
}

func mergeRuntimeConfig(base RuntimeConfig, remote *serviceConfigResponse) RuntimeConfig {
	if remote == nil {
		return base
	}
	if strings.TrimSpace(remote.AppEnv) != "" {
		base.AppEnv = strings.TrimSpace(remote.AppEnv)
	}
	if strings.TrimSpace(remote.DBDriver) != "" {
		base.DBDriver = strings.TrimSpace(remote.DBDriver)
	}
	if remote.AutoMigrate != nil {
		base.AutoMigrate = *remote.AutoMigrate
	}
	if remote.KafkaEnabled != nil {
		base.KafkaEnabled = *remote.KafkaEnabled
	}
	if len(remote.KafkaBrokers) > 0 {
		base.KafkaBrokers = remote.KafkaBrokers
	}
	if strings.TrimSpace(remote.KafkaClientID) != "" {
		base.KafkaClientID = strings.TrimSpace(remote.KafkaClientID)
	}
	if strings.TrimSpace(remote.KafkaConsumerGroup) != "" {
		base.KafkaConsumerGroup = strings.TrimSpace(remote.KafkaConsumerGroup)
	}
	if remote.KafkaWriteTimeoutMS > 0 {
		base.KafkaWriteTimeoutMS = remote.KafkaWriteTimeoutMS
	}
	if strings.TrimSpace(remote.APIGatewayAllowedOrigin) != "" {
		base.APIGatewayAllowedOrigin = strings.TrimSpace(remote.APIGatewayAllowedOrigin)
	}
	if len(remote.CORSAllowedOrigins) > 0 {
		base.CORSAllowedOrigins = remote.CORSAllowedOrigins
	}
	base.CORSAllowedOrigins = appendIfMissing(base.CORSAllowedOrigins, base.APIGatewayAllowedOrigin)

	if strings.TrimSpace(remote.KafkaTopics.DeviceRegistered) != "" {
		base.KafkaTopics.DeviceRegistered = strings.TrimSpace(remote.KafkaTopics.DeviceRegistered)
	}
	if strings.TrimSpace(remote.KafkaTopics.DeviceStatusUpdated) != "" {
		base.KafkaTopics.DeviceStatusUpdated = strings.TrimSpace(remote.KafkaTopics.DeviceStatusUpdated)
	}
	if strings.TrimSpace(remote.KafkaTopics.DeviceLinked) != "" {
		base.KafkaTopics.DeviceLinked = strings.TrimSpace(remote.KafkaTopics.DeviceLinked)
	}
	if strings.TrimSpace(remote.KafkaTopics.DeviceUnlinked) != "" {
		base.KafkaTopics.DeviceUnlinked = strings.TrimSpace(remote.KafkaTopics.DeviceUnlinked)
	}
	if strings.TrimSpace(remote.KafkaTopics.DeviceConfigurationUpdated) != "" {
		base.KafkaTopics.DeviceConfigurationUpdated = strings.TrimSpace(remote.KafkaTopics.DeviceConfigurationUpdated)
	}
	if strings.TrimSpace(remote.KafkaTopics.DeviceEventRecorded) != "" {
		base.KafkaTopics.DeviceEventRecorded = strings.TrimSpace(remote.KafkaTopics.DeviceEventRecorded)
	}

	return base
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
