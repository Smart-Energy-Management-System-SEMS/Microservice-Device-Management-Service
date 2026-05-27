package entities

import (
	"strings"
	"time"

	shared "device-management-service/shared/domain"
	"github.com/google/uuid"
)

type DeviceEvent struct {
	EventID     uuid.UUID
	DeviceID    uuid.UUID
	EventType   string
	Description *string
	OccurredAt  time.Time
}

func NewDeviceEvent(deviceID uuid.UUID, eventType string, description *string, occurredAt *time.Time) (*DeviceEvent, error) {
	if deviceID == uuid.Nil {
		return nil, shared.NewValidationError("device_id is required")
	}
	if strings.TrimSpace(eventType) == "" {
		return nil, shared.NewValidationError("event_type is required")
	}
	eventTime := time.Now().UTC()
	if occurredAt != nil {
		eventTime = occurredAt.UTC()
	}
	return &DeviceEvent{
		EventID:     uuid.New(),
		DeviceID:    deviceID,
		EventType:   strings.TrimSpace(eventType),
		Description: description,
		OccurredAt:  eventTime,
	}, nil
}
