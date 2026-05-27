package model

import (
	"time"

	"github.com/google/uuid"
)

type DeviceBindingModel struct {
	BindingID     uuid.UUID  `gorm:"type:uuid;primaryKey;column:binding_id"`
	DeviceID      uuid.UUID  `gorm:"type:uuid;not null;index;column:device_id"`
	UserID        uuid.UUID  `gorm:"type:uuid;not null;index;column:user_id"`
	HomeID        *uuid.UUID `gorm:"type:uuid;column:home_id"`
	BindingStatus string     `gorm:"type:varchar(30);not null;column:binding_status"`
	LinkedAt      time.Time  `gorm:"not null;column:linked_at"`
	UnlinkedAt    *time.Time `gorm:"column:unlinked_at"`
	UpdatedAt     time.Time  `gorm:"not null;column:updated_at"`
}

func (DeviceBindingModel) TableName() string {
	return "device_bindings"
}
