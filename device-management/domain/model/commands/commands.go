package commands

import (
	"time"

	"github.com/google/uuid"
)

type RegisterDeviceCommand struct {
	ExternalDeviceCode string
	UserID             uuid.UUID
	DeviceName         string
	DeviceType         string
	Brand              *string
	Model              *string
	ConnectionProtocol string
}

type UpdateDeviceCommand struct {
	DeviceID           uuid.UUID
	DeviceName         string
	DeviceType         string
	Brand              *string
	Model              *string
	ConnectionProtocol string
}

type UpdateDeviceStatusCommand struct {
	DeviceID uuid.UUID
	Status   string
}

type DeleteDeviceCommand struct {
	DeviceID uuid.UUID
}

type CreateDeviceBindingCommand struct {
	DeviceID uuid.UUID
	UserID   uuid.UUID
	HomeID   *uuid.UUID
}

type UnlinkDeviceBindingCommand struct {
	BindingID uuid.UUID
}

type CreateDeviceConfigurationCommand struct {
	DeviceID    uuid.UUID
	ConfigKey   string
	ConfigValue *string
}

type UpdateDeviceConfigurationCommand struct {
	ConfigurationID uuid.UUID
	ConfigValue     *string
}

type RecordDeviceEventCommand struct {
	DeviceID    uuid.UUID
	EventType   string
	Description *string
	OccurredAt  *time.Time
}
