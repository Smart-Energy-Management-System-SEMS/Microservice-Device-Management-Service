package repositories

import (
	"device-management-service/device-management/domain/model/aggregates"
	"device-management-service/device-management/domain/model/entities"
	"device-management-service/device-management/domain/model/valueobjects"
	persistencemodel "device-management-service/device-management/infrastructure/persistence/gorm/model"
)

// This file holds the "mappers": small functions that translate between the two
// worlds the repository lives in. We deliberately keep TWO separate shapes for
// the same data:
//   - domain objects (aggregates/entities): rich, with behaviour and value
//     objects, and unaware of any database.
//   - persistence models: plain structs with GORM tags, shaped for SQL columns.
// Keeping them apart means the database schema can change without forcing the
// domain to change, and vice versa. Each concept therefore has a pair of
// functions: toXModel (domain -> DB) and toXDomain (DB -> domain).

// toDeviceModel converts a domain Device into the row struct GORM persists.
// Note how value objects are cast back to plain strings (e.g. string(status))
// because the database column is just text.
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

// toDeviceDomain does the reverse: it rebuilds a domain Device from a DB row.
// It goes through RehydrateDevice (not RegisterDevice) on purpose, because the
// stored data is already valid and should be reconstructed as-is. The plain
// strings from the DB are wrapped back into their value object types.
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

// The remaining functions repeat the same domain<->DB mapping pattern for the
// other concepts in the bounded context: bindings, configurations and events.
// They are intentionally repetitive and boring — that predictability makes the
// data flow easy to follow and to test.

// toBindingModel maps a domain DeviceBinding to its persistence model.
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
