package queryservices

import (
	"context"

	"device-management-service/domain/model/entities"
	"device-management-service/domain/model/queries"
	"device-management-service/domain/repositories"
)

type DeviceBindingQueryService struct {
	bindingRepository repositories.DeviceBindingRepository
}

func NewDeviceBindingQueryService(bindingRepository repositories.DeviceBindingRepository) *DeviceBindingQueryService {
	return &DeviceBindingQueryService{bindingRepository: bindingRepository}
}

func (s *DeviceBindingQueryService) GetByDeviceID(ctx context.Context, query queries.GetDeviceBindingsQuery) ([]entities.DeviceBinding, error) {
	return s.bindingRepository.FindByDeviceID(ctx, query.DeviceID)
}

func (s *DeviceBindingQueryService) GetByUserID(ctx context.Context, query queries.GetUserBindingsQuery) ([]entities.DeviceBinding, error) {
	return s.bindingRepository.FindByUserID(ctx, query.UserID)
}
