// Package commandservices contains the "command" side of a CQRS-style design.
// CQRS (Command Query Responsibility Segregation) splits operations that CHANGE
// data (commands: create, update, delete) from operations that only READ data
// (queries). This file handles the write operations for devices and is part of
// the application layer, which orchestrates the domain but contains no business
// rules of its own.
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

// DeviceCommandService wires together the collaborators it needs. Note that the
// fields are interfaces (DeviceRepository, ExternalReferenceService), not
// concrete types. This is "dependency injection": the service depends on
// abstractions, so we can swap a real database for a fake one in tests.
type DeviceCommandService struct {
	deviceRepository repositories.DeviceRepository                // saves/loads devices
	referenceService acl.ExternalReferenceService                 // validates data owned by other services
	eventHandler     *eventhandlers.DeviceIntegrationEventHandler // publishes events to Kafka
}

// NewDeviceCommandService is the constructor. The dependencies are passed in
// from the outside (usually wired up once at application start-up) instead of
// being created here.
func NewDeviceCommandService(deviceRepository repositories.DeviceRepository, referenceService acl.ExternalReferenceService, eventHandler *eventhandlers.DeviceIntegrationEventHandler) *DeviceCommandService {
	return &DeviceCommandService{deviceRepository: deviceRepository, referenceService: referenceService, eventHandler: eventHandler}
}

// RegisterDevice is the use case for creating a new device. It reads top to
// bottom as a recipe, and every step returns early on error (Go's typical
// "if err != nil { return }" style). The ctx (context.Context) is passed all
// the way down so the operation can be cancelled or timed out as a whole.
func (s *DeviceCommandService) RegisterDevice(ctx context.Context, command commands.RegisterDeviceCommand) (*aggregates.Device, error) {
	// Step 1: make sure the referenced user really exists. The user lives in
	// another microservice, so we ask through the ACL reference service.
	if err := s.referenceService.ValidateUserReference(ctx, command.UserID); err != nil {
		return nil, err
	}
	// Step 2: enforce that the external code is unique. The logic is a little
	// tricky: if FindByExternalDeviceCode succeeds (err == nil) the code is
	// already taken, so we report a conflict. If it fails, we only continue
	// when the failure is specifically a NOT_FOUND; any other error (e.g. a DB
	// outage) is a real problem and is returned. errors.As inspects the error
	// to see whether it is one of our AppError values.
	if _, err := s.deviceRepository.FindByExternalDeviceCode(ctx, command.ExternalDeviceCode); err == nil {
		return nil, dmerrors.NewConflictError("external_device_code already exists")
	} else {
		var appErr *dmerrors.AppError
		if !errors.As(err, &appErr) || appErr.Code != dmerrors.ErrNotFound {
			return nil, err
		}
	}
	// Step 3: turn the raw protocol string into a validated value object.
	protocol, err := valueobjects.NewConnectionProtocol(command.ConnectionProtocol)
	if err != nil {
		return nil, err
	}
	// Step 4: let the domain build the device (this is where the creation rules
	// run). The application layer does not re-implement those rules; it trusts
	// the aggregate.
	device, err := aggregates.RegisterDevice(command.ExternalDeviceCode, command.UserID, command.DeviceName, command.DeviceType, command.Brand, command.Model, protocol)
	if err != nil {
		return nil, err
	}
	// Step 5: persist the new device through the repository.
	if err := s.deviceRepository.Save(ctx, device); err != nil {
		return nil, err
	}
	// Step 6: announce what happened by publishing an integration event so other
	// microservices can react. This is "fire and forget": a failure to publish
	// is logged inside the handler but does not roll back the save.
	s.eventHandler.Publish(ctx, s.eventHandler.Topics().DeviceEvents, eventhandlers.EventTypeDeviceRegistered, device.DeviceID, device.UserID, map[string]interface{}{
		"externalDeviceCode": device.ExternalDeviceCode,
		"deviceType":         device.DeviceType,
		"status":             device.Status,
	})
	return device, nil
}

// UpdateDevice follows the classic "load, change, save" pattern: fetch the
// existing device, ask the domain to apply the change (which re-validates the
// rules), and persist the result.
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

// UpdateDeviceStatus changes only the device status. We remember the previous
// status before changing it so we can include both the old and the new value in
// the event payload, which is useful for other services that audit changes.
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
	s.eventHandler.Publish(ctx, s.eventHandler.Topics().DeviceEvents, eventhandlers.EventTypeDeviceStatusUpdated, device.DeviceID, device.UserID, map[string]interface{}{
		"previousStatus": previousStatus,
		"newStatus":      device.Status,
	})
	return device, nil
}

// DeleteDevice performs a "soft delete": it does not erase the row, it asks the
// device to move to the REMOVED status and saves it. We still return the device
// (now marked as removed) and publish an event so others learn about it.
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
	s.eventHandler.Publish(ctx, s.eventHandler.Topics().DeviceEvents, eventhandlers.EventTypeDeviceStatusUpdated, device.DeviceID, device.UserID, map[string]interface{}{
		"newStatus": device.Status,
	})
	return device, nil
}
