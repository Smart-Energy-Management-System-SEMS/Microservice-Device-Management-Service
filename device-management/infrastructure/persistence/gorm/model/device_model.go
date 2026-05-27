package model

import (
	"time"

	"github.com/google/uuid"
)

type DeviceModel struct {
	DeviceID           uuid.UUID `gorm:"type:uuid;primaryKey;column:device_id"`
	ExternalDeviceCode string    `gorm:"type:varchar(120);uniqueIndex;not null;column:external_device_code"`
	UserID             uuid.UUID `gorm:"type:uuid;not null;column:user_id"`
	DeviceName         string    `gorm:"type:varchar(120);not null;column:device_name"`
	DeviceType         string    `gorm:"type:varchar(80);not null;column:device_type"`
	Brand              *string   `gorm:"type:varchar(80);column:brand"`
	Model              *string   `gorm:"type:varchar(80);column:model"`
	ConnectionProtocol string    `gorm:"type:varchar(30);not null;column:connection_protocol"`
	Status             string    `gorm:"type:varchar(30);not null;column:status"`
	RegisteredAt       time.Time `gorm:"not null;column:registered_at"`
	UpdatedAt          time.Time `gorm:"not null;column:updated_at"`
}

func (DeviceModel) TableName() string {
	return "devices"
}
