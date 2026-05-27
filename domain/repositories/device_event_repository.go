package repositories

import (
	"context"

	"device-management-service/domain/model/entities"
	"github.com/google/uuid"
)

type DeviceEventRepository interface {
	Save(ctx context.Context, event *entities.DeviceEvent) error
	FindByDeviceID(ctx context.Context, deviceID uuid.UUID) ([]entities.DeviceEvent, error)
}
