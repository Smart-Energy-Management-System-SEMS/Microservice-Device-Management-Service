package entities

import (
	"time"

	dmerrors "device-management-service/device-management/domain"
	"device-management-service/device-management/domain/model/valueobjects"
	"github.com/google/uuid"
)

type DeviceBinding struct {
	BindingID     uuid.UUID
	DeviceID      uuid.UUID
	UserID        uuid.UUID
	HomeID        *uuid.UUID
	BindingStatus valueobjects.BindingStatus
	LinkedAt      time.Time
	UnlinkedAt    *time.Time
	UpdatedAt     time.Time
}

func NewDeviceBinding(deviceID, userID uuid.UUID, homeID *uuid.UUID) (*DeviceBinding, error) {
	if deviceID == uuid.Nil {
		return nil, dmerrors.NewValidationError("device_id is required")
	}
	if userID == uuid.Nil {
		return nil, dmerrors.NewValidationError("user_id is required")
	}
	now := time.Now().UTC()
	return &DeviceBinding{
		BindingID:     uuid.New(),
		DeviceID:      deviceID,
		UserID:        userID,
		HomeID:        homeID,
		BindingStatus: valueobjects.BindingStatusLinked,
		LinkedAt:      now,
		UpdatedAt:     now,
	}, nil
}

func (b *DeviceBinding) Unlink() error {
	if b.BindingStatus == valueobjects.BindingStatusUnlinked {
		return dmerrors.NewConflictError("binding is already unlinked")
	}
	now := time.Now().UTC()
	b.BindingStatus = valueobjects.BindingStatusUnlinked
	b.UnlinkedAt = &now
	b.UpdatedAt = now
	return nil
}
