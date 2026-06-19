// Package queryservices is the "query" side of CQRS — the read-only half. While
// the command services change state and publish events, these services only
// fetch data. Notice how short they are: reads usually have no business rules,
// so the service just forwards the call to the repository. Keeping reads
// separate makes each side easier to understand and optimise on its own.
package queryservices

import (
	"context"

	"device-management-service/device-management/domain/model/aggregates"
	"device-management-service/device-management/domain/model/queries"
	"device-management-service/device-management/domain/repositories"
)

// DeviceQueryService only needs the repository to read devices. It depends on
// the repository interface, not a concrete database, just like the command side.
type DeviceQueryService struct {
	deviceRepository repositories.DeviceRepository
}

// NewDeviceQueryService injects the repository dependency.
func NewDeviceQueryService(deviceRepository repositories.DeviceRepository) *DeviceQueryService {
	return &DeviceQueryService{deviceRepository: deviceRepository}
}

// GetAll returns every device. It returns a slice ([]aggregates.Device); an
// empty slice (not an error) is the normal result when there are no devices.
func (s *DeviceQueryService) GetAll(ctx context.Context) ([]aggregates.Device, error) {
	return s.deviceRepository.FindAll(ctx)
}

// GetByID looks up a single device. The query object (GetDeviceByIDQuery) wraps
// the parameter so the method signature stays stable even if we add more search
// fields later.
func (s *DeviceQueryService) GetByID(ctx context.Context, query queries.GetDeviceByIDQuery) (*aggregates.Device, error) {
	return s.deviceRepository.FindByID(ctx, query.DeviceID)
}

// GetByUserID returns all devices that belong to a given user.
func (s *DeviceQueryService) GetByUserID(ctx context.Context, query queries.GetDevicesByUserQuery) ([]aggregates.Device, error) {
	return s.deviceRepository.FindByUserID(ctx, query.UserID)
}
