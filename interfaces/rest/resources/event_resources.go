package resources

import "time"

type CreateDeviceEventResource struct {
	EventType   string     `json:"eventType" binding:"required"`
	Description *string    `json:"description"`
	OccurredAt  *time.Time `json:"occurredAt"`
}

type DeviceEventResource struct {
	EventID     string    `json:"eventId"`
	DeviceID    string    `json:"deviceId"`
	EventType   string    `json:"eventType"`
	Description *string   `json:"description,omitempty"`
	OccurredAt  time.Time `json:"occurredAt"`
}
