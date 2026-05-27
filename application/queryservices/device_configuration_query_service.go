package queryservices

import (
	"context"

	"device-management-service/domain/model/entities"
	"device-management-service/domain/model/queries"
	"device-management-service/domain/repositories"
)

type DeviceConfigurationQueryService struct {
	configurationRepository repositories.DeviceConfigurationRepository
}

func NewDeviceConfigurationQueryService(configurationRepository repositories.DeviceConfigurationRepository) *DeviceConfigurationQueryService {
	return &DeviceConfigurationQueryService{configurationRepository: configurationRepository}
}

func (s *DeviceConfigurationQueryService) GetByDeviceID(ctx context.Context, query queries.GetDeviceConfigurationsQuery) ([]entities.DeviceConfiguration, error) {
	return s.configurationRepository.FindByDeviceID(ctx, query.DeviceID)
}
