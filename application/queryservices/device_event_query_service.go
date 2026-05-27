package queryservices

import (
	"context"

	"device-management-service/domain/model/entities"
	"device-management-service/domain/model/queries"
	"device-management-service/domain/repositories"
)

type DeviceEventQueryService struct {
	eventRepository repositories.DeviceEventRepository
}

func NewDeviceEventQueryService(eventRepository repositories.DeviceEventRepository) *DeviceEventQueryService {
	return &DeviceEventQueryService{eventRepository: eventRepository}
}

func (s *DeviceEventQueryService) GetByDeviceID(ctx context.Context, query queries.GetDeviceEventsQuery) ([]entities.DeviceEvent, error) {
	return s.eventRepository.FindByDeviceID(ctx, query.DeviceID)
}
