// Package outboundservices defines the "ports" through which our application
// talks to the outside world. In hexagonal (ports & adapters) architecture, a
// port is an interface owned by the inner layers; the concrete "adapter" that
// implements it lives in the infrastructure layer. Here we declare WHAT we need
// (a way to publish events) without saying HOW it is done.
package outboundservices

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// IntegrationEvent is the message shape we send to other microservices. The
// json tags control how each field is named once the struct is serialised to
// JSON (e.g. EventID becomes "eventId"). Payload is a map so different event
// types can carry different extra fields without needing a new struct each time.
type IntegrationEvent struct {
	EventID    uuid.UUID              `json:"eventId"`
	EventType  string                 `json:"eventType"`
	DeviceID   uuid.UUID              `json:"deviceId"`
	UserID     uuid.UUID              `json:"userId"`
	OccurredAt time.Time              `json:"occurredAt"`
	Payload    map[string]interface{} `json:"payload"`
}

// DeviceEventPublisher is the port (interface). The application depends only on
// this small contract: "give me a topic and an event, and publish it". The real
// Kafka producer in the infrastructure layer satisfies this interface, but for
// tests we could provide a fake that records calls instead of sending them.
type DeviceEventPublisher interface {
	Publish(ctx context.Context, topic string, event IntegrationEvent) error
}
