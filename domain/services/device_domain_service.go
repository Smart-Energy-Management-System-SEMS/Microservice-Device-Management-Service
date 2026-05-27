package services

import (
	"device-management-service/domain/model/aggregates"
	"device-management-service/domain/model/entities"
)

type DeviceDomainService struct{}

func NewDeviceDomainService() *DeviceDomainService {
	return &DeviceDomainService{}
}

func (s *DeviceDomainService) ValidateBinding(device *aggregates.Device) error {
	return device.CanBeBound()
}

func (s *DeviceDomainService) ValidateConfigurationUpdate(device *aggregates.Device) error {
	return device.CanUpdateConfiguration()
}

func (s *DeviceDomainService) UnlinkBinding(binding *entities.DeviceBinding) error {
	return binding.Unlink()
}
