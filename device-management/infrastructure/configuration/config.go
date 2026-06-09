// Package configuration is responsible for loading all the settings the service
// needs to run. It follows the Twelve-Factor App idea of "config in the
// environment": defaults live in code, environment variables override them, and
// on top of that a central config microservice can override them again. This
// keeps secrets and per-environment values out of the source code.
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

// AppConfig holds the few values we need very early, before anything else is
// wired up (the port to listen on, where the config service is, etc.).
type AppConfig struct {
	ServiceName      string
	Port             string
	ConfigServiceURL string
	DatabaseURL      string
}

// RuntimeConfig is the full set of settings used while the service is running:
// database driver, Kafka options, CORS rules and the event topic names.
type RuntimeConfig struct {
	AppEnv                  string
	DBDriver                string
	AutoMigrate             bool
	KafkaEnabled            bool
	KafkaBrokers            []string
	KafkaSecurityProtocol   string
	KafkaSASLMechanism      string
	KafkaUsername           string
	KafkaPassword           string
	KafkaClientID           string
	KafkaConsumerGroup      string
	KafkaWriteTimeoutMS     int
	APIGatewayAllowedOrigin string
	CORSAllowedOrigins      []string
	KafkaTopics             KafkaTopics
}

// KafkaTopics groups the topic name used for each kind of event. Storing them
// in config (instead of hard-coding the strings) lets each environment use its
// own topic names without recompiling.
type KafkaTopics struct {
	DeviceRegistered           string
	DeviceStatusUpdated        string
	DeviceLinked               string
	DeviceUnlinked             string
	DeviceConfigurationUpdated string
	DeviceEventRecorded        string
}

// serviceConfigResponse mirrors the JSON returned by the central config
// service. Some fields are pointers (*bool) on purpose: a nil pointer means
// "the remote config did not mention this setting", which we must distinguish
// from an explicit "false". With a plain bool we could not tell those apart.
type serviceConfigResponse struct {
	AppEnv                  string      `json:"appEnv"`
	DBDriver                string      `json:"dbDriver"`
	AutoMigrate             *bool       `json:"autoMigrate"`
	KafkaEnabled            *bool       `json:"kafkaEnabled"`
	KafkaBrokers            []string    `json:"kafkaBrokers"`
	KafkaSecurityProtocol   string      `json:"kafkaSecurityProtocol"`
	KafkaSASLMechanism      string      `json:"kafkaSaslMechanism"`
	KafkaUsername           string      `json:"kafkaUsername"`
	KafkaPassword           string      `json:"kafkaPassword"`
	KafkaClientID           string      `json:"kafkaClientId"`
	KafkaConsumerGroup      string      `json:"kafkaConsumerGroup"`
	KafkaWriteTimeoutMS     int         `json:"kafkaWriteTimeoutMs"`
	APIGatewayAllowedOrigin string      `json:"apiGatewayAllowedOrigin"`
	CORSAllowedOrigins      []string    `json:"corsAllowedOrigins"`
	KafkaTopics             KafkaTopics `json:"kafkaTopics"`
}

// LoadAppConfig reads the bootstrap settings from environment variables, each
// with a sensible default so the service can still start locally with no setup.
func LoadAppConfig() AppConfig {
	return AppConfig{
		ServiceName:      getEnv("SERVICE_NAME", "device-management-service"),
		Port:             getEnv("PORT", "8083"),
		ConfigServiceURL: strings.TrimRight(getEnv("CONFIG_SERVICE_URL", "http://localhost:8090"), "/"),
		DatabaseURL:      getEnv("DATABASE_URL", ""),
	}
}

// ResolveRuntimeConfig builds the final configuration in two steps:
//  1. Start from local defaults/env vars (this always works).
//  2. Try to fetch overrides from the central config service. If that call
//     fails (service down, timeout...) we gracefully fall back to the local
//     config instead of crashing. This pattern keeps the service resilient.
func ResolveRuntimeConfig(ctx context.Context, appConfig AppConfig) RuntimeConfig {
	localFallback := defaultRuntimeConfig(appConfig.ServiceName)
	serviceConfig, err := fetchServiceConfig(ctx, appConfig.ConfigServiceURL, appConfig.ServiceName)
	if err != nil {
		return localFallback
	}
	return mergeRuntimeConfig(localFallback, serviceConfig)
}

// fetchServiceConfig performs the HTTP GET to the central config service and
// decodes the JSON answer. A short 2-second timeout is set on the client so a
// slow config service cannot block our start-up forever.
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
	// "defer" schedules Body.Close() to run when the function returns, no matter
	// which path we take. Forgetting to close the body would leak connections.
	defer response.Body.Close()
	// Any status outside the 2xx range means the request did not succeed.
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("config service returned status %d", response.StatusCode)
	}

	var payload serviceConfigResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		return nil, err
	}
	return &payload, nil
}

// defaultRuntimeConfig builds the baseline configuration purely from env vars
// and hard-coded defaults. This is what we use when the central config service
// is unreachable, and also the base that remote values are merged on top of.
func defaultRuntimeConfig(serviceName string) RuntimeConfig {
	apiGatewayOrigin := getEnv("API_GATEWAY_ALLOWED_ORIGIN", "http://localhost:8081")
	corsAllowedOrigins := loadAllowedOrigins("http://localhost:3000,http://localhost:5173,http://localhost:8081")
	corsAllowedOrigins = appendIfMissing(corsAllowedOrigins, apiGatewayOrigin)

	return RuntimeConfig{
		AppEnv:                  getEnv("APP_ENV", "local"),
		DBDriver:                getEnv("DB_DRIVER", "postgres"),
		AutoMigrate:             getBoolEnv("AUTO_MIGRATE", false),
		KafkaEnabled:            getBoolEnv("KAFKA_ENABLED", false),
		KafkaBrokers:            splitCSV(getEnv("KAFKA_BROKERS", "localhost:9092")),
		KafkaSecurityProtocol:   getEnv("KAFKA_SECURITY_PROTOCOL", ""),
		KafkaSASLMechanism:      getEnv("KAFKA_SASL_MECHANISM", ""),
		KafkaUsername:           getEnv("KAFKA_USERNAME", ""),
		KafkaPassword:           getEnv("KAFKA_PASSWORD", ""),
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

// mergeRuntimeConfig overlays the remote values on top of the local base. The
// rule throughout is "only override when the remote actually provided a value":
// for strings that means non-empty after trimming, for slices a non-zero
// length, and for the *bool pointers a non-nil value. This way a partial remote
// config never accidentally wipes a good local default.
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
	if strings.TrimSpace(remote.KafkaSecurityProtocol) != "" {
		base.KafkaSecurityProtocol = strings.TrimSpace(remote.KafkaSecurityProtocol)
	}
	if strings.TrimSpace(remote.KafkaSASLMechanism) != "" {
		base.KafkaSASLMechanism = strings.TrimSpace(remote.KafkaSASLMechanism)
	}
	if strings.TrimSpace(remote.KafkaUsername) != "" {
		base.KafkaUsername = strings.TrimSpace(remote.KafkaUsername)
	}
	if strings.TrimSpace(remote.KafkaPassword) != "" {
		base.KafkaPassword = strings.TrimSpace(remote.KafkaPassword)
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

// The helpers below are small wrappers around os.Getenv. They all share the
// same idea: read the variable, and if it is missing or blank, return a
// fallback. They keep the configuration code above clean and repetition-free.

// getEnv reads a string env var or returns the fallback.
func getEnv(key string, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

// getBoolEnv reads a boolean env var. It is lenient about what counts as true
// ("true", "1" or "yes"), which is friendlier for people setting the variable.
func getBoolEnv(key string, fallback bool) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	if value == "" {
		return fallback
	}
	return value == "true" || value == "1" || value == "yes"
}

// getIntEnv reads an integer env var. If the text is not a valid number,
// strconv.Atoi returns an error and we fall back instead of crashing.
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

// splitCSV turns a comma-separated string like "a, b ,c" into a clean slice
// ["a", "b", "c"], dropping empty entries. We pre-size the slice with make(...,
// 0, len(parts)) as a small optimisation to avoid repeated re-allocations.
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

func loadAllowedOrigins(fallback string) []string {
	value := strings.TrimSpace(os.Getenv("CORS_ALLOWED_ORIGINS"))
	if value == "" {
		value = strings.TrimSpace(os.Getenv("ALLOWED_ORIGINS"))
	}
	if value == "" {
		value = fallback
	}
	return splitCSV(value)
}

// appendIfMissing adds a value to a slice only if it is not already there,
// acting like a tiny "set". We use it to make sure the API gateway origin is
// always part of the allowed CORS origins without creating duplicates.
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
