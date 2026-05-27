package commandservices

import (
	"context"
	"errors"

	"device-management-service/device-management/application/eventhandlers"
	dmerrors "device-management-service/device-management/domain"
	"device-management-service/device-management/domain/model/aggregates"
	"device-management-service/device-management/domain/model/commands"
	"device-management-service/device-management/domain/model/valueobjects"
	"device-management-service/device-management/domain/repositories"
	"device-management-service/device-management/interfaces/acl"
)

type DeviceCommandService struct {
	deviceRepository repositories.DeviceRepository
	referenceService acl.ExternalReferenceService
	eventHandler     *eventhandlers.DeviceIntegrationEventHandler
}

func NewDeviceCommandService(deviceRepository repositories.DeviceRepository, referenceService acl.ExternalReferenceService, eventHandler *eventhandlers.DeviceIntegrationEventHandler) *DeviceCommandService {
	return &DeviceCommandService{deviceRepository: deviceRepository, referenceService: referenceService, eventHandler: eventHandler}
}

func (s *DeviceCommandService) RegisterDevice(ctx context.Context, command commands.RegisterDeviceCommand) (*aggregates.Device, error) {
	if err := s.referenceService.ValidateUserReference(ctx, command.UserID); err != nil {
		return nil, err
	}
	if _, err := s.deviceRepository.FindByExternalDeviceCode(ctx, command.ExternalDeviceCode); err == nil {
		return nil, dmerrors.NewConflictError("external_device_code already exists")
	} else {
		var appErr *dmerrors.AppError
		if !errors.As(err, &appErr) || appErr.Code != dmerrors.ErrNotFound {
			return nil, err
		}
	}
	protocol, err := valueobjects.NewConnectionProtocol(command.ConnectionProtocol)
	if err != nil {
		return nil, err
	}
	device, err := aggregates.RegisterDevice(command.ExternalDeviceCode, command.UserID, command.DeviceName, command.DeviceType, command.Brand, command.Model, protocol)
	if err != nil {
		return nil, err
	}
	if err := s.deviceRepository.Save(ctx, device); err != nil {
		return nil, err
	}
	s.eventHandler.Publish(ctx, eventhandlers.TopicDeviceRegistered, eventhandlers.EventTypeDeviceRegistered, device.DeviceID, device.UserID, map[string]interface{}{
		"externalDeviceCode": device.ExternalDeviceCode,
		"deviceType":         device.DeviceType,
		"status":             device.Status,
	})
	return device, nil
}

func (s *DeviceCommandService) UpdateDevice(ctx context.Context, command commands.UpdateDeviceCommand) (*aggregates.Device, error) {
	device, err := s.deviceRepository.FindByID(ctx, command.DeviceID)
	if err != nil {
		return nil, err
	}
	protocol, err := valueobjects.NewConnectionProtocol(command.ConnectionProtocol)
	if err != nil {
		return nil, err
	}
	if err := device.UpdateDetails(command.DeviceName, command.DeviceType, command.Brand, command.Model, protocol); err != nil {
		return nil, err
	}
	if err := s.deviceRepository.Update(ctx, device); err != nil {
		return nil, err
	}
	return device, nil
}

func (s *DeviceCommandService) UpdateDeviceStatus(ctx context.Context, command commands.UpdateDeviceStatusCommand) (*aggregates.Device, error) {
	device, err := s.deviceRepository.FindByID(ctx, command.DeviceID)
	if err != nil {
		return nil, err
	}
	previousStatus := device.Status
	nextStatus, err := valueobjects.NewDeviceStatus(command.Status)
	if err != nil {
		return nil, err
	}
	if err := device.ChangeStatus(nextStatus); err != nil {
		return nil, err
	}
	if err := s.deviceRepository.Update(ctx, device); err != nil {
		return nil, err
	}
	s.eventHandler.Publish(ctx, eventhandlers.TopicDeviceStatusUpdated, eventhandlers.EventTypeDeviceStatusUpdated, device.DeviceID, device.UserID, map[string]interface{}{
		"previousStatus": previousStatus,
		"newStatus":      device.Status,
	})
	return device, nil
}

func (s *DeviceCommandService) DeleteDevice(ctx context.Context, command commands.DeleteDeviceCommand) (*aggregates.Device, error) {
	device, err := s.deviceRepository.FindByID(ctx, command.DeviceID)
	if err != nil {
		return nil, err
	}
	if err := device.Remove(); err != nil {
		return nil, err
	}
	if err := s.deviceRepository.Update(ctx, device); err != nil {
		return nil, err
	}
	s.eventHandler.Publish(ctx, eventhandlers.TopicDeviceStatusUpdated, eventhandlers.EventTypeDeviceStatusUpdated, device.DeviceID, device.UserID, map[string]interface{}{
		"newStatus": device.Status,
	})
	return device, nil
}
