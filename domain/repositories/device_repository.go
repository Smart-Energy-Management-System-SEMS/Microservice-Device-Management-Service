package repositories

import (
	"context"

	"device-management-service/domain/model/aggregates"
	"github.com/google/uuid"
)

type DeviceRepository interface {
	Save(ctx context.Context, device *aggregates.Device) error
	Update(ctx context.Context, device *aggregates.Device) error
	FindByID(ctx context.Context, deviceID uuid.UUID) (*aggregates.Device, error)
	FindByExternalDeviceCode(ctx context.Context, code string) (*aggregates.Device, error)
	FindAll(ctx context.Context) ([]aggregates.Device, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]aggregates.Device, error)
}
