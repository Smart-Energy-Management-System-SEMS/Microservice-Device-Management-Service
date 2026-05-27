package queryservices

import (
	"context"

	"device-management-service/device-management/domain/model/aggregates"
	"device-management-service/device-management/domain/model/queries"
	"device-management-service/device-management/domain/repositories"
)

type DeviceQueryService struct {
	deviceRepository repositories.DeviceRepository
}

func NewDeviceQueryService(deviceRepository repositories.DeviceRepository) *DeviceQueryService {
	return &DeviceQueryService{deviceRepository: deviceRepository}
}

func (s *DeviceQueryService) GetAll(ctx context.Context) ([]aggregates.Device, error) {
	return s.deviceRepository.FindAll(ctx)
}

func (s *DeviceQueryService) GetByID(ctx context.Context, query queries.GetDeviceByIDQuery) (*aggregates.Device, error) {
	return s.deviceRepository.FindByID(ctx, query.DeviceID)
}

func (s *DeviceQueryService) GetByUserID(ctx context.Context, query queries.GetDevicesByUserQuery) ([]aggregates.Device, error) {
	return s.deviceRepository.FindByUserID(ctx, query.UserID)
}
