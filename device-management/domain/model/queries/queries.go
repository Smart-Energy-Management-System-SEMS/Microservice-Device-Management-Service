package queries

import "github.com/google/uuid"

type GetDeviceByIDQuery struct {
	DeviceID uuid.UUID
}

type GetDevicesByUserQuery struct {
	UserID uuid.UUID
}

type GetDeviceBindingsQuery struct {
	DeviceID uuid.UUID
}

type GetUserBindingsQuery struct {
	UserID uuid.UUID
}

type GetDeviceConfigurationsQuery struct {
	DeviceID uuid.UUID
}

type GetDeviceEventsQuery struct {
	DeviceID uuid.UUID
}
