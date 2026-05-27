package entities

import (
	"strings"
	"time"

	shared "device-management-service/shared/domain"
	"github.com/google/uuid"
)

type DeviceConfiguration struct {
	ConfigurationID uuid.UUID
	DeviceID        uuid.UUID
	ConfigKey       string
	ConfigValue     *string
	UpdatedAt       time.Time
}

func NewDeviceConfiguration(deviceID uuid.UUID, key string, value *string) (*DeviceConfiguration, error) {
	if deviceID == uuid.Nil {
		return nil, shared.NewValidationError("device_id is required")
	}
	if strings.TrimSpace(key) == "" {
		return nil, shared.NewValidationError("config_key is required")
	}
	return &DeviceConfiguration{
		ConfigurationID: uuid.New(),
		DeviceID:        deviceID,
		ConfigKey:       strings.TrimSpace(key),
		ConfigValue:     value,
		UpdatedAt:       time.Now().UTC(),
	}, nil
}

func (c *DeviceConfiguration) Update(value *string) {
	c.ConfigValue = value
	c.UpdatedAt = time.Now().UTC()
}
