package aggregates

import (
	"strings"
	"time"

	"device-management-service/device-management/domain/model/valueobjects"
	dmerrors "device-management-service/device-management/domain"
	"github.com/google/uuid"
)

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

func RegisterDevice(externalCode string, userID uuid.UUID, name string, deviceType string, brand *string, model *string, protocol valueobjects.ConnectionProtocol) (*Device, error) {
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

func (d *Device) Remove() error {
	return d.ChangeStatus(valueobjects.DeviceStatusRemoved)
}

func (d *Device) CanBeBound() error {
	if d.IsRemoved() {
		return dmerrors.NewConflictError("removed devices cannot be linked")
	}
	return nil
}

func (d *Device) CanUpdateConfiguration() error {
	if d.IsRemoved() {
		return dmerrors.NewConflictError("removed devices cannot update configuration")
	}
	return nil
}

func (d *Device) IsRemoved() bool {
	return d.Status == valueobjects.DeviceStatusRemoved
}

func normalizeOptional(value *string) *string {
	if value == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}
