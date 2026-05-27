package eventhandlers

import (
	"context"
	"log"
	"time"

	"device-management-service/device-management/application/outboundservices"
	"github.com/google/uuid"
)

const (
	TopicDeviceRegistered             = "device.registered"
	TopicDeviceStatusUpdated          = "device.status.updated"
	TopicDeviceLinked                 = "device.linked"
	TopicDeviceUnlinked               = "device.unlinked"
	TopicDeviceConfigurationUpdated   = "device.configuration.updated"
	TopicDeviceEventRecorded          = "device.event.recorded"
	EventTypeDeviceRegistered         = "DEVICE_REGISTERED"
	EventTypeDeviceStatusUpdated      = "DEVICE_STATUS_UPDATED"
	EventTypeDeviceLinked             = "DEVICE_LINKED"
	EventTypeDeviceUnlinked           = "DEVICE_UNLINKED"
	EventTypeDeviceConfigurationSaved = "DEVICE_CONFIGURATION_UPDATED"
	EventTypeDeviceEventRecorded      = "DEVICE_EVENT_RECORDED"
)

type DeviceIntegrationEventHandler struct {
	publisher outboundservices.DeviceEventPublisher
}

func NewDeviceIntegrationEventHandler(publisher outboundservices.DeviceEventPublisher) *DeviceIntegrationEventHandler {
	return &DeviceIntegrationEventHandler{publisher: publisher}
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
