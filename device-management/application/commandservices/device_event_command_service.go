package commandservices

import (
	"context"

	"device-management-service/device-management/application/eventhandlers"
	"device-management-service/device-management/domain/model/commands"
	"device-management-service/device-management/domain/model/entities"
	"device-management-service/device-management/domain/repositories"
)

type DeviceEventCommandService struct {
	deviceRepository repositories.DeviceRepository
	eventRepository  repositories.DeviceEventRepository
	eventHandler     *eventhandlers.DeviceIntegrationEventHandler
}

func NewDeviceEventCommandService(deviceRepository repositories.DeviceRepository, eventRepository repositories.DeviceEventRepository, eventHandler *eventhandlers.DeviceIntegrationEventHandler) *DeviceEventCommandService {
	return &DeviceEventCommandService{deviceRepository: deviceRepository, eventRepository: eventRepository, eventHandler: eventHandler}
}

func (s *DeviceEventCommandService) RecordEvent(ctx context.Context, command commands.RecordDeviceEventCommand) (*entities.DeviceEvent, error) {
	device, err := s.deviceRepository.FindByID(ctx, command.DeviceID)
	if err != nil {
		return nil, err
	}
	event, err := entities.NewDeviceEvent(command.DeviceID, command.EventType, command.Description, command.OccurredAt)
	if err != nil {
		return nil, err
	}
	if err := s.eventRepository.Save(ctx, event); err != nil {
		return nil, err
	}
	s.eventHandler.Publish(ctx, eventhandlers.TopicDeviceEventRecorded, eventhandlers.EventTypeDeviceEventRecorded, event.DeviceID, device.UserID, map[string]interface{}{
		"eventId":   event.EventID,
		"eventType": event.EventType,
	})
	return event, nil
}
