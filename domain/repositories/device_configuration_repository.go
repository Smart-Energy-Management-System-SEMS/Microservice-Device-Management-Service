package repositories

import (
	"context"

	"device-management-service/domain/model/entities"
	"github.com/google/uuid"
)

type DeviceConfigurationRepository interface {
	Save(ctx context.Context, configuration *entities.DeviceConfiguration) error
	Update(ctx context.Context, configuration *entities.DeviceConfiguration) error
	FindByID(ctx context.Context, configurationID uuid.UUID) (*entities.DeviceConfiguration, error)
	FindByDeviceID(ctx context.Context, deviceID uuid.UUID) ([]entities.DeviceConfiguration, error)
}
