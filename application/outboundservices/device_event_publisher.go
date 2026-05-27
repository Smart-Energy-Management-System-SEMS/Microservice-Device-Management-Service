package outboundservices

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type IntegrationEvent struct {
	EventID    uuid.UUID              `json:"eventId"`
	EventType  string                 `json:"eventType"`
	DeviceID   uuid.UUID              `json:"deviceId"`
	UserID     uuid.UUID              `json:"userId"`
	OccurredAt time.Time              `json:"occurredAt"`
	Payload    map[string]interface{} `json:"payload"`
}

type DeviceEventPublisher interface {
	Publish(ctx context.Context, topic string, event IntegrationEvent) error
}
