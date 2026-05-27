package model

import (
	"time"

	"github.com/google/uuid"
)

type DeviceEventModel struct {
	EventID     uuid.UUID `gorm:"type:uuid;primaryKey;column:event_id"`
	DeviceID    uuid.UUID `gorm:"type:uuid;not null;index;column:device_id"`
	EventType   string    `gorm:"type:varchar(80);not null;column:event_type"`
	Description *string   `gorm:"type:text;column:description"`
	OccurredAt  time.Time `gorm:"not null;column:occurred_at"`
}

func (DeviceEventModel) TableName() string {
	return "device_events"
}
