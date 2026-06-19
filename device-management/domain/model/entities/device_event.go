package entities

import (
	"strings"
	"time"

	dmerrors "device-management-service/device-management/domain"
	"github.com/google/uuid"
)

var allowedDeviceEventTypes = map[string]struct{}{
	"CONNECTED":    {},
	"DISCONNECTED": {},
	"ERROR":        {},
	"UPDATED":      {},
	"REMOVED":      {},
}

type DeviceEvent struct {
	EventID     uuid.UUID
	DeviceID    uuid.UUID
	EventType   string
	Description *string
	OccurredAt  time.Time
}

func NewDeviceEvent(deviceID uuid.UUID, eventType string, description *string, occurredAt *time.Time) (*DeviceEvent, error) {
	if deviceID == uuid.Nil {
		return nil, dmerrors.NewValidationError("device_id is required")
	}
	normalizedEventType := strings.ToUpper(strings.TrimSpace(eventType))
	if normalizedEventType == "" {
		return nil, dmerrors.NewValidationError("event_type is required")
	}
	if _, ok := allowedDeviceEventTypes[normalizedEventType]; !ok {
		return nil, dmerrors.NewValidationError("event_type must be one of CONNECTED, DISCONNECTED, ERROR, UPDATED or REMOVED")
	}
	eventTime := time.Now().UTC()
	if occurredAt != nil {
		eventTime = occurredAt.UTC()
	}
	return &DeviceEvent{
		EventID:     uuid.New(),
		DeviceID:    deviceID,
		EventType:   normalizedEventType,
		Description: description,
		OccurredAt:  eventTime,
	}, nil
}
