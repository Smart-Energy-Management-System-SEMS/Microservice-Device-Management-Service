package commandservices

import (
	"context"

	"device-management-service/device-management/application/eventhandlers"
	"device-management-service/device-management/domain/model/commands"
	"device-management-service/device-management/domain/model/entities"
	"device-management-service/device-management/domain/repositories"
	"device-management-service/device-management/domain/services"
)

type DeviceConfigurationCommandService struct {
	deviceRepository        repositories.DeviceRepository
	configurationRepository repositories.DeviceConfigurationRepository
	domainService           *services.DeviceDomainService
	eventHandler            *eventhandlers.DeviceIntegrationEventHandler
}

func NewDeviceConfigurationCommandService(deviceRepository repositories.DeviceRepository, configurationRepository repositories.DeviceConfigurationRepository, domainService *services.DeviceDomainService, eventHandler *eventhandlers.DeviceIntegrationEventHandler) *DeviceConfigurationCommandService {
	return &DeviceConfigurationCommandService{
		deviceRepository:        deviceRepository,
		configurationRepository: configurationRepository,
		domainService:           domainService,
		eventHandler:            eventHandler,
	}
}

func (s *DeviceConfigurationCommandService) CreateConfiguration(ctx context.Context, command commands.CreateDeviceConfigurationCommand) (*entities.DeviceConfiguration, error) {
	device, err := s.deviceRepository.FindByID(ctx, command.DeviceID)
	if err != nil {
		return nil, err
	}
	if err := s.domainService.ValidateConfigurationUpdate(device); err != nil {
		return nil, err
	}
	configuration, err := entities.NewDeviceConfiguration(command.DeviceID, command.ConfigKey, command.ConfigValue)
	if err != nil {
		return nil, err
	}
	if err := s.configurationRepository.Save(ctx, configuration); err != nil {
		return nil, err
	}
	s.eventHandler.Publish(ctx, eventhandlers.TopicDeviceConfigurationUpdated, eventhandlers.EventTypeDeviceConfigurationSaved, configuration.DeviceID, device.UserID, map[string]interface{}{
		"configurationId": configuration.ConfigurationID,
		"configKey":       configuration.ConfigKey,
	})
	return configuration, nil
}

func (s *DeviceConfigurationCommandService) UpdateConfiguration(ctx context.Context, command commands.UpdateDeviceConfigurationCommand) (*entities.DeviceConfiguration, error) {
	configuration, err := s.configurationRepository.FindByID(ctx, command.ConfigurationID)
	if err != nil {
		return nil, err
	}
	device, err := s.deviceRepository.FindByID(ctx, configuration.DeviceID)
	if err != nil {
		return nil, err
	}
	if err := s.domainService.ValidateConfigurationUpdate(device); err != nil {
		return nil, err
	}
	configuration.Update(command.ConfigValue)
	if err := s.configurationRepository.Update(ctx, configuration); err != nil {
		return nil, err
	}
	s.eventHandler.Publish(ctx, eventhandlers.TopicDeviceConfigurationUpdated, eventhandlers.EventTypeDeviceConfigurationSaved, configuration.DeviceID, device.UserID, map[string]interface{}{
		"configurationId": configuration.ConfigurationID,
		"configKey":       configuration.ConfigKey,
	})
	return configuration, nil
}
