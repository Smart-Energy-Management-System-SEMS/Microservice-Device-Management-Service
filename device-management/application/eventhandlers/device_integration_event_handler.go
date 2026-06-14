// Package eventhandlers builds and sends "integration events". An integration
// event is a message we publish so OTHER microservices know something happened
// here (a device was registered, its status changed, etc.). This handler sits
// between the application services and the messaging infrastructure: services
// call it with simple arguments, and it takes care of building the event and
// handing it to the publisher.
package eventhandlers

import (
	"context"
	"log"
	"time"

	"device-management-service/device-management/application/outboundservices"
	"device-management-service/device-management/infrastructure/configuration"
	"github.com/google/uuid"
)

// These constants are the event "type" labels that travel inside each message.
// Consumers read this field to decide how to react, so the exact strings are
// part of our contract with other services and must not change carelessly.
const (
	EventTypeDeviceRegistered         = "device.registered"
	EventTypeDeviceStatusUpdated      = "device.status.updated"
	EventTypeDeviceLinked             = "device.linked"
	EventTypeDeviceUnlinked           = "device.unlinked"
	EventTypeDeviceConfigurationSaved = "device.configuration.updated"
	EventTypeDeviceEventRecorded      = "device.event.recorded"
)

// DeviceIntegrationEventHandler depends on a publisher interface (so the actual
// transport, Kafka, can be swapped) and on the configured topic names.
type DeviceIntegrationEventHandler struct {
	publisher outboundservices.DeviceEventPublisher // where messages are sent
	topics    configuration.KafkaTopics             // the topic name for each event
}

// NewDeviceIntegrationEventHandler injects the publisher and topic config.
func NewDeviceIntegrationEventHandler(publisher outboundservices.DeviceEventPublisher, topics configuration.KafkaTopics) *DeviceIntegrationEventHandler {
	return &DeviceIntegrationEventHandler{publisher: publisher, topics: topics}
}

// Publish assembles an IntegrationEvent and sends it. Two design choices worth
// noting:
//   - It is defensive: the nil checks let callers use the handler even when no
//     publisher was configured (for example, when Kafka is disabled), in which
//     case the method simply does nothing instead of crashing.
//   - It does not return an error. Publishing is best-effort: a failure is
//     logged but does not break the main use case that triggered the event.
func (h *DeviceIntegrationEventHandler) Publish(ctx context.Context, topic string, eventType string, deviceID uuid.UUID, userID uuid.UUID, payload map[string]interface{}) {
	if h == nil || h.publisher == nil {
		return
	}
	// We generate a unique EventID and stamp the time in UTC. The payload is a
	// free-form map so each event type can carry whatever fields it needs.
	event := outboundservices.IntegrationEvent{
		EventID:    uuid.New(),
		EventType:  eventType,
		DeviceID:   deviceID,
		UserID:     userID,
		OccurredAt: time.Now().UTC(),
		Payload:    payload,
	}
	// If sending fails we only log it. The device was already saved, so we do
	// not want to fail the whole request just because the message bus hiccuped.
	if err := h.publisher.Publish(ctx, topic, event); err != nil {
		log.Printf("kafka publish failed: %v", err)
	}
}

// Topics exposes the configured topic names to the callers (the command
// services use it to pick the right topic per event). The nil check returns a
// zero-value struct so this is always safe to call.
func (h *DeviceIntegrationEventHandler) Topics() configuration.KafkaTopics {
	if h == nil {
		return configuration.KafkaTopics{}
	}
	return h.topics
}
