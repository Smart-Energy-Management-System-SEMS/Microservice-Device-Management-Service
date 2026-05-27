package model

import (
	"time"

	"github.com/google/uuid"
)

type DeviceConfigurationModel struct {
	ConfigurationID uuid.UUID `gorm:"type:uuid;primaryKey;column:configuration_id"`
	DeviceID        uuid.UUID `gorm:"type:uuid;not null;index;column:device_id"`
	ConfigKey       string    `gorm:"type:varchar(100);not null;column:config_key"`
	ConfigValue     *string   `gorm:"type:text;column:config_value"`
	UpdatedAt       time.Time `gorm:"not null;column:updated_at"`
}

func (DeviceConfigurationModel) TableName() string {
	return "device_configurations"
}
