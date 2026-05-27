package repositories

import (
	"context"

	"device-management-service/domain/model/entities"
	"github.com/google/uuid"
)

type DeviceBindingRepository interface {
	Save(ctx context.Context, binding *entities.DeviceBinding) error
	Update(ctx context.Context, binding *entities.DeviceBinding) error
	FindByID(ctx context.Context, bindingID uuid.UUID) (*entities.DeviceBinding, error)
	FindByDeviceID(ctx context.Context, deviceID uuid.UUID) ([]entities.DeviceBinding, error)
	FindByUserID(ctx context.Context, userID uuid.UUID) ([]entities.DeviceBinding, error)
}
