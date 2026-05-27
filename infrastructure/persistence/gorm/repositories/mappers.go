package repositories

import (
	"device-management-service/domain/model/aggregates"
	"device-management-service/domain/model/entities"
	"device-management-service/domain/model/valueobjects"
	persistencemodel "device-management-service/infrastructure/persistence/gorm/model"
)

func toDeviceModel(device *aggregates.Device) persistencemodel.DeviceModel {
	return persistencemodel.DeviceModel{
		DeviceID:           device.DeviceID,
		ExternalDeviceCode: device.ExternalDeviceCode,
		UserID:             device.UserID,
		DeviceName:         device.DeviceName,
		DeviceType:         device.DeviceType,
		Brand:              device.Brand,
		Model:              device.Model,
		ConnectionProtocol: string(device.ConnectionProtocol),
		Status:             string(device.Status),
		RegisteredAt:       device.RegisteredAt,
		UpdatedAt:          device.UpdatedAt,
	}
}

func toDeviceDomain(device persistencemodel.DeviceModel) *aggregates.Device {
	return aggregates.RehydrateDevice(
		device.DeviceID,
		device.ExternalDeviceCode,
		device.UserID,
		device.DeviceName,
		device.DeviceType,
		device.Brand,
		device.Model,
		valueobjects.ConnectionProtocol(device.ConnectionProtocol),
		valueobjects.DeviceStatus(device.Status),
		device.RegisteredAt,
		device.UpdatedAt,
	)
}

func toBindingModel(binding *entities.DeviceBinding) persistencemodel.DeviceBindingModel {
	return persistencemodel.DeviceBindingModel{
		BindingID:     binding.BindingID,
		DeviceID:      binding.DeviceID,
		UserID:        binding.UserID,
		HomeID:        binding.HomeID,
		BindingStatus: string(binding.BindingStatus),
		LinkedAt:      binding.LinkedAt,
		UnlinkedAt:    binding.UnlinkedAt,
		UpdatedAt:     binding.UpdatedAt,
	}
}

func toBindingDomain(binding persistencemodel.DeviceBindingModel) *entities.DeviceBinding {
	return &entities.DeviceBinding{
		BindingID:     binding.BindingID,
		DeviceID:      binding.DeviceID,
		UserID:        binding.UserID,
		HomeID:        binding.HomeID,
		BindingStatus: valueobjects.BindingStatus(binding.BindingStatus),
		LinkedAt:      binding.LinkedAt,
		UnlinkedAt:    binding.UnlinkedAt,
		UpdatedAt:     binding.UpdatedAt,
	}
}

func toConfigurationModel(configuration *entities.DeviceConfiguration) persistencemodel.DeviceConfigurationModel {
	return persistencemodel.DeviceConfigurationModel{
		ConfigurationID: configuration.ConfigurationID,
		DeviceID:        configuration.DeviceID,
		ConfigKey:       configuration.ConfigKey,
		ConfigValue:     configuration.ConfigValue,
		UpdatedAt:       configuration.UpdatedAt,
	}
}

func toConfigurationDomain(configuration persistencemodel.DeviceConfigurationModel) *entities.DeviceConfiguration {
	return &entities.DeviceConfiguration{
		ConfigurationID: configuration.ConfigurationID,
		DeviceID:        configuration.DeviceID,
		ConfigKey:       configuration.ConfigKey,
		ConfigValue:     configuration.ConfigValue,
		UpdatedAt:       configuration.UpdatedAt,
	}
}

func toEventModel(event *entities.DeviceEvent) persistencemodel.DeviceEventModel {
	return persistencemodel.DeviceEventModel{
		EventID:     event.EventID,
		DeviceID:    event.DeviceID,
		EventType:   event.EventType,
		Description: event.Description,
		OccurredAt:  event.OccurredAt,
	}
}

func toEventDomain(event persistencemodel.DeviceEventModel) *entities.DeviceEvent {
	return &entities.DeviceEvent{
		EventID:     event.EventID,
		DeviceID:    event.DeviceID,
		EventType:   event.EventType,
		Description: event.Description,
		OccurredAt:  event.OccurredAt,
	}
}
