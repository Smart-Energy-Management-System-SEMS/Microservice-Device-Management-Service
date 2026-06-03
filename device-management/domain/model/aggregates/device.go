// Package aggregates holds the "aggregate roots" of our domain model.
// In Domain-Driven Design (DDD) an aggregate root is the main object that
// guards a group of related data and makes sure it always stays in a valid
// state. Here the Device is the root: any change to a device must go through
// its methods so the business rules are never skipped.
package aggregates

import (
	"strings"
	"time"

	// We alias the domain package as "dmerrors" so it is obvious that we are
	// using our own custom error helpers (validation, conflict, etc.) and not
	// Go's standard errors package.
	dmerrors "device-management-service/device-management/domain"
	"device-management-service/device-management/domain/model/valueobjects"
	"github.com/google/uuid"
)

// Device represents a smart device registered in the system. It is the
// aggregate root, so all of its business behaviour lives on this struct.
// Brand and Model are pointers (*string) because they are optional: a nil
// pointer means "no value was provided", which is different from an empty
// string.
type Device struct {
	DeviceID           uuid.UUID
	ExternalDeviceCode string
	UserID             uuid.UUID
	DeviceName         string
	DeviceType         string
	Brand              *string
	Model              *string
	ConnectionProtocol valueobjects.ConnectionProtocol
	Status             valueobjects.DeviceStatus
	RegisteredAt       time.Time
	UpdatedAt          time.Time
}

// RegisterDevice is a "factory function": it is the only correct way to build
// a brand new Device. Instead of letting callers create a Device struct by
// hand (which could leave it half-filled or invalid), we validate every
// required field first and only then return the object. This guarantees that a
// Device cannot exist in an invalid state. It returns (nil, error) when a rule
// is broken, following Go's idiomatic "value, error" pattern.
func RegisterDevice(externalCode string, userID uuid.UUID, name string, deviceType string, brand *string, model *string, protocol valueobjects.ConnectionProtocol) (*Device, error) {
	// Each check below protects one business rule. We trim spaces so that a
	// value made only of blanks (like "   ") is treated as empty.
	if strings.TrimSpace(externalCode) == "" {
		return nil, dmerrors.NewValidationError("external_device_code is required")
	}
	if userID == uuid.Nil {
		return nil, dmerrors.NewValidationError("user_id is required")
	}
	if strings.TrimSpace(name) == "" {
		return nil, dmerrors.NewValidationError("device_name is required")
	}
	if strings.TrimSpace(deviceType) == "" {
		return nil, dmerrors.NewValidationError("device_type is required")
	}
	if !protocol.IsValid() {
		return nil, dmerrors.NewValidationError("connection_protocol is invalid")
	}
	// We store timestamps in UTC so the data is consistent no matter which
	// server timezone the service runs on. A new device starts ACTIVE and gets
	// a freshly generated UUID as its identity.
	now := time.Now().UTC()
	return &Device{
		DeviceID:           uuid.New(),
		ExternalDeviceCode: strings.TrimSpace(externalCode),
		UserID:             userID,
		DeviceName:         strings.TrimSpace(name),
		DeviceType:         strings.TrimSpace(deviceType),
		Brand:              normalizeOptional(brand),
		Model:              normalizeOptional(model),
		ConnectionProtocol: protocol,
		Status:             valueobjects.DeviceStatusActive,
		RegisteredAt:       now,
		UpdatedAt:          now,
	}, nil
}

// RehydrateDevice rebuilds a Device from data that already exists (for example,
// a row loaded from the database). The important difference with RegisterDevice
// is that it does NOT run the creation validations: the data was already valid
// when it was first saved, so we trust it and simply reconstruct the object.
// "Rehydrate" is a common name for turning stored data back into a live object.
func RehydrateDevice(deviceID uuid.UUID, externalCode string, userID uuid.UUID, name string, deviceType string, brand *string, model *string, protocol valueobjects.ConnectionProtocol, status valueobjects.DeviceStatus, registeredAt time.Time, updatedAt time.Time) *Device {
	return &Device{
		DeviceID:           deviceID,
		ExternalDeviceCode: externalCode,
		UserID:             userID,
		DeviceName:         name,
		DeviceType:         deviceType,
		Brand:              brand,
		Model:              model,
		ConnectionProtocol: protocol,
		Status:             status,
		RegisteredAt:       registeredAt,
		UpdatedAt:          updatedAt,
	}
}

// UpdateDetails changes the editable fields of a device. Notice it is a method
// with a pointer receiver (d *Device): the pointer lets us modify the original
// object instead of a copy. The first rule is that a removed device is "frozen"
// and can no longer be edited, so we reject the change early.
func (d *Device) UpdateDetails(name string, deviceType string, brand *string, model *string, protocol valueobjects.ConnectionProtocol) error {
	if d.IsRemoved() {
		return dmerrors.NewConflictError("removed devices cannot be updated")
	}
	if strings.TrimSpace(name) == "" {
		return dmerrors.NewValidationError("device_name is required")
	}
	if strings.TrimSpace(deviceType) == "" {
		return dmerrors.NewValidationError("device_type is required")
	}
	if !protocol.IsValid() {
		return dmerrors.NewValidationError("connection_protocol is invalid")
	}
	d.DeviceName = strings.TrimSpace(name)
	d.DeviceType = strings.TrimSpace(deviceType)
	d.Brand = normalizeOptional(brand)
	d.Model = normalizeOptional(model)
	d.ConnectionProtocol = protocol
	d.UpdatedAt = time.Now().UTC()
	return nil
}

// ChangeStatus moves the device to a new status, but only if the move is
// allowed. This implements a simple "state machine": the value object decides
// (via CanTransitionTo) which jumps between statuses make sense. Keeping this
// rule inside the domain prevents illegal transitions like reviving a device
// that was already removed.
func (d *Device) ChangeStatus(next valueobjects.DeviceStatus) error {
	if !next.IsValid() {
		return dmerrors.NewValidationError("device status is invalid")
	}
	if !d.Status.CanTransitionTo(next) {
		return dmerrors.NewConflictError("invalid device status transition")
	}
	d.Status = next
	d.UpdatedAt = time.Now().UTC()
	return nil
}

// Remove is a small convenience method: deleting a device is really just a
// status change to REMOVED, so we reuse ChangeStatus instead of duplicating the
// transition rules. (This is a "soft delete": the row stays, only its status
// changes.)
func (d *Device) Remove() error {
	return d.ChangeStatus(valueobjects.DeviceStatusRemoved)
}

// CanBeBound answers the question "is this device allowed to be linked to a
// home/user?". It returns nil when it is fine, or an error explaining why not.
// Returning an error instead of a bool lets the caller forward a clear message.
func (d *Device) CanBeBound() error {
	if d.IsRemoved() {
		return dmerrors.NewConflictError("removed devices cannot be linked")
	}
	return nil
}

// CanUpdateConfiguration works like CanBeBound but guards configuration
// changes: a removed device should not accept new settings.
func (d *Device) CanUpdateConfiguration() error {
	if d.IsRemoved() {
		return dmerrors.NewConflictError("removed devices cannot update configuration")
	}
	return nil
}

// IsRemoved is a tiny helper that keeps the rest of the code readable: instead
// of comparing the status everywhere, we ask the device directly.
func (d *Device) IsRemoved() bool {
	return d.Status == valueobjects.DeviceStatusRemoved
}

// normalizeOptional cleans up an optional string. It is lowercase, so it is
// private to this package (Go uses capitalization to control visibility). If
// the input is nil or becomes empty after trimming, we return nil so that
// "blank" and "not provided" are treated the same way.
func normalizeOptional(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	// We return the address of the local "trimmed" variable. This is safe in
	// Go: the compiler keeps the value alive as long as the pointer is used.
	return &trimmed
}
