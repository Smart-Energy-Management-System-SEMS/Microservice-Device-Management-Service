package commandservices

import (
	"context"

	"device-management-service/device-management/application/eventhandlers"
	"device-management-service/device-management/domain/model/commands"
	"device-management-service/device-management/domain/model/entities"
	"device-management-service/device-management/domain/repositories"
	"device-management-service/device-management/domain/services"
)

type DeviceBindingCommandService struct {
	deviceRepository  repositories.DeviceRepository
	bindingRepository repositories.DeviceBindingRepository
	domainService     *services.DeviceDomainService
	eventHandler      *eventhandlers.DeviceIntegrationEventHandler
}

func NewDeviceBindingCommandService(deviceRepository repositories.DeviceRepository, bindingRepository repositories.DeviceBindingRepository, domainService *services.DeviceDomainService, eventHandler *eventhandlers.DeviceIntegrationEventHandler) *DeviceBindingCommandService {
	return &DeviceBindingCommandService{
		deviceRepository:  deviceRepository,
		bindingRepository: bindingRepository,
		domainService:     domainService,
		eventHandler:      eventHandler,
	}
}

func (s *DeviceBindingCommandService) CreateBinding(ctx context.Context, command commands.CreateDeviceBindingCommand) (*entities.DeviceBinding, error) {
	device, err := s.deviceRepository.FindByID(ctx, command.DeviceID)
	if err != nil {
		return nil, err
	}
	if err := s.domainService.ValidateBinding(device); err != nil {
		return nil, err
	}
	binding, err := entities.NewDeviceBinding(command.DeviceID, command.UserID, command.HomeID)
	if err != nil {
		return nil, err
	}
	if err := s.bindingRepository.Save(ctx, binding); err != nil {
		return nil, err
	}
	s.eventHandler.Publish(ctx, eventhandlers.TopicDeviceLinked, eventhandlers.EventTypeDeviceLinked, binding.DeviceID, binding.UserID, map[string]interface{}{
		"bindingId": binding.BindingID,
		"homeId":    binding.HomeID,
	})
	return binding, nil
}

func (s *DeviceBindingCommandService) UnlinkBinding(ctx context.Context, command commands.UnlinkDeviceBindingCommand) (*entities.DeviceBinding, error) {
	binding, err := s.bindingRepository.FindByID(ctx, command.BindingID)
	if err != nil {
		return nil, err
	}
	if err := s.domainService.UnlinkBinding(binding); err != nil {
		return nil, err
	}
	if err := s.bindingRepository.Update(ctx, binding); err != nil {
		return nil, err
	}
	s.eventHandler.Publish(ctx, eventhandlers.TopicDeviceUnlinked, eventhandlers.EventTypeDeviceUnlinked, binding.DeviceID, binding.UserID, map[string]interface{}{
		"bindingId":  binding.BindingID,
		"unlinkedAt": binding.UnlinkedAt,
	})
	return binding, nil
}
