package entities

import (
	"time"

	"device-management-service/domain/model/valueobjects"
	shared "device-management-service/shared/domain"
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
		return nil, shared.NewValidationError("device_id is required")
	}
	if userID == uuid.Nil {
		return nil, shared.NewValidationError("user_id is required")
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
		return shared.NewConflictError("binding is already unlinked")
	}
	now := time.Now().UTC()
	b.BindingStatus = valueobjects.BindingStatusUnlinked
	b.UnlinkedAt = &now
	b.UpdatedAt = now
	return nil
}
