// Package services holds "domain services". A domain service is where we put
// business logic that does not naturally belong to a single entity. In DDD, if
// a rule fits inside one object we put it there; if it coordinates several
// objects (or just reads nicely as its own action) it goes into a service.
package services

import (
	"device-management-service/device-management/domain/model/aggregates"
	"device-management-service/device-management/domain/model/entities"
)

// DeviceDomainService is an empty struct: it holds no data, it only groups
// related domain operations. An empty struct{} in Go uses zero memory, so this
// is a cheap, stateless service.
type DeviceDomainService struct{}

// NewDeviceDomainService is the constructor. Even though the struct is empty, we
// keep a constructor for consistency with the rest of the codebase and so we
// can add dependencies later without changing how callers create it.
func NewDeviceDomainService() *DeviceDomainService {
	return &DeviceDomainService{}
}

// ValidateBinding checks whether a device is allowed to be linked. Right now it
// simply delegates to the device itself (device.CanBeBound). This thin wrapper
// gives us a single place to expand the rule later without touching callers.
func (s *DeviceDomainService) ValidateBinding(device *aggregates.Device) error {
	return device.CanBeBound()
}

// ValidateConfigurationUpdate checks whether a device may receive new
// configuration, again delegating to the aggregate that owns the rule.
func (s *DeviceDomainService) ValidateConfigurationUpdate(device *aggregates.Device) error {
	return device.CanUpdateConfiguration()
}

// UnlinkBinding asks a binding entity to mark itself as unlinked. The service
// stays simple and forwards the work to the entity that knows the rule.
func (s *DeviceDomainService) UnlinkBinding(binding *entities.DeviceBinding) error {
	return binding.Unlink()
}
