package valueobjects

import shared "device-management-service/shared/domain"

type DeviceStatus string

const (
	DeviceStatusActive       DeviceStatus = "ACTIVE"
	DeviceStatusInactive     DeviceStatus = "INACTIVE"
	DeviceStatusDisconnected DeviceStatus = "DISCONNECTED"
	DeviceStatusRemoved      DeviceStatus = "REMOVED"
)

func NewDeviceStatus(value string) (DeviceStatus, error) {
	status := DeviceStatus(value)
	if !status.IsValid() {
		return "", shared.NewValidationError("invalid device status")
	}
	return status, nil
}

func (s DeviceStatus) IsValid() bool {
	switch s {
	case DeviceStatusActive, DeviceStatusInactive, DeviceStatusDisconnected, DeviceStatusRemoved:
		return true
	default:
		return false
	}
}

func (s DeviceStatus) CanTransitionTo(next DeviceStatus) bool {
	if s == DeviceStatusRemoved {
		return false
	}
	if next == DeviceStatusRemoved {
		return true
	}
	return next == DeviceStatusActive || next == DeviceStatusInactive || next == DeviceStatusDisconnected
}
