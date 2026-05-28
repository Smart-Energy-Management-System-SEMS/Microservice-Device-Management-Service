package eventhandlers

import (
	"context"
	"log"
	"time"

	"device-management-service/device-management/application/outboundservices"
	"device-management-service/device-management/infrastructure/configuration"
	"github.com/google/uuid"
)

const (
	EventTypeDeviceRegistered         = "DEVICE_REGISTERED"
	EventTypeDeviceStatusUpdated      = "DEVICE_STATUS_UPDATED"
	EventTypeDeviceLinked             = "DEVICE_LINKED"
	EventTypeDeviceUnlinked           = "DEVICE_UNLINKED"
	EventTypeDeviceConfigurationSaved = "DEVICE_CONFIGURATION_UPDATED"
	EventTypeDeviceEventRecorded      = "DEVICE_EVENT_RECORDED"
)

type DeviceIntegrationEventHandler struct {
	publisher outboundservices.DeviceEventPublisher
	topics    configuration.KafkaTopics
}

func NewDeviceIntegrationEventHandler(publisher outboundservices.DeviceEventPublisher, topics configuration.KafkaTopics) *DeviceIntegrationEventHandler {
	return &DeviceIntegrationEventHandler{publisher: publisher, topics: topics}
}

func (h *DeviceIntegrationEventHandler) Publish(ctx context.Context, topic string, eventType string, deviceID uuid.UUID, userID uuid.UUID, payload map[string]interface{}) {
	if h == nil || h.publisher == nil {
		return
	}
	event := outboundservices.IntegrationEvent{
		EventID:    uuid.New(),
		EventType:  eventType,
		DeviceID:   deviceID,
		UserID:     userID,
		OccurredAt: time.Now().UTC(),
		Payload:    payload,
	}
	if err := h.publisher.Publish(ctx, topic, event); err != nil {
		log.Printf("kafka publish failed: %v", err)
	}
}

func (h *DeviceIntegrationEventHandler) Topics() configuration.KafkaTopics {
	if h == nil {
		return configuration.KafkaTopics{}
	}
	return h.topics
}
