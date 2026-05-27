package resources

import "time"

type CreateDeviceConfigurationResource struct {
	ConfigKey   string  `json:"configKey" binding:"required"`
	ConfigValue *string `json:"configValue"`
}

type UpdateDeviceConfigurationResource struct {
	ConfigValue *string `json:"configValue"`
}

type DeviceConfigurationResource struct {
	ConfigurationID string    `json:"configurationId"`
	DeviceID        string    `json:"deviceId"`
	ConfigKey       string    `json:"configKey"`
	ConfigValue     *string   `json:"configValue,omitempty"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
