// Package valueobjects contains "value objects": small immutable types that
// represent a concept by their value, not by an identity. A DeviceStatus has no
// ID of its own; two ACTIVE statuses are considered equal. Wrapping a plain
// string in its own type lets us attach validation and behaviour to it.
package valueobjects

import dmerrors "device-management-service/device-management/domain"

// DeviceStatus is defined as a named string type. Using a dedicated type
// instead of a raw string means the compiler can help us: a function expecting
// a DeviceStatus will not accidentally accept any random string.
type DeviceStatus string

// These constants are the only valid statuses. Listing them as typed constants
// gives us an "enum-like" set in Go (the language has no built-in enum).
const (
	DeviceStatusActive       DeviceStatus = "ACTIVE"
	DeviceStatusInactive     DeviceStatus = "INACTIVE"
	DeviceStatusDisconnected DeviceStatus = "DISCONNECTED"
	DeviceStatusRemoved      DeviceStatus = "REMOVED"
)

// NewDeviceStatus safely converts an untrusted string (for example, coming from
// an HTTP request) into a DeviceStatus. It validates first and returns an error
// if the text does not match one of the constants above.
func NewDeviceStatus(value string) (DeviceStatus, error) {
	status := DeviceStatus(value)
	if !status.IsValid() {
		return "", dmerrors.NewValidationError("invalid device status")
	}
	return status, nil
}

// IsValid reports whether the status is one of the known constants. The switch
// lists every allowed case and the default catches anything unexpected. Note
// the value receiver (s DeviceStatus): value objects are read-only, so we do
// not need a pointer here.
func (s DeviceStatus) IsValid() bool {
	switch s {
	case DeviceStatusActive, DeviceStatusInactive, DeviceStatusDisconnected, DeviceStatusRemoved:
		return true
	default:
		return false
	}
}

// CanTransitionTo encodes the rules of the status state machine: which status
// can become which. The order of the checks matters:
//  1. A REMOVED device is final, so it can never transition anywhere.
//  2. Any non-removed device is allowed to be removed.
//  3. Otherwise, only "live" statuses (active/inactive/disconnected) are valid.
func (s DeviceStatus) CanTransitionTo(next DeviceStatus) bool {
	if s == DeviceStatusRemoved {
		return false
	}
	if next == DeviceStatusRemoved {
		return true
	}
	return next == DeviceStatusActive || next == DeviceStatusInactive || next == DeviceStatusDisconnected
}
