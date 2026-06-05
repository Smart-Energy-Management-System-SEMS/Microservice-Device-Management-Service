package commandservices

import (
	"context"

	"device-management-service/device-management/application/eventhandlers"
	"device-management-service/device-management/domain/model/commands"
	"device-management-service/device-management/domain/model/entities"
	"device-management-service/device-management/domain/repositories"
	"device-management-service/device-management/domain/services"
	"device-management-service/device-management/interfaces/acl"
)
// DeviceBindingCommandService es el servicio de aplicación encargado de manejar
// los comandos relacionados con la vinculación de dispositivos a usuarios y hogares.
// Implementa la lógica de negocio para crear y desenlazar dispositivos.
type DeviceBindingCommandService struct {
	deviceRepository  repositories.DeviceRepository
	bindingRepository repositories.DeviceBindingRepository
	domainService     *services.DeviceDomainService
	referenceService  acl.ExternalReferenceService
	eventHandler      *eventhandlers.DeviceIntegrationEventHandler
}
// NewDeviceBindingCommandService crea una nueva instancia del servicio de comando de vinculación de dispositivos
// Recibe todas las dependencias necesarias para el funcionamiento del servicio
// Retorna un puntero al servicio inicializado
func NewDeviceBindingCommandService(deviceRepository repositories.DeviceRepository, bindingRepository repositories.DeviceBindingRepository, domainService *services.DeviceDomainService, referenceService acl.ExternalReferenceService, eventHandler *eventhandlers.DeviceIntegrationEventHandler) *DeviceBindingCommandService {
	return &DeviceBindingCommandService{
		deviceRepository:  deviceRepository,
		bindingRepository: bindingRepository,
		domainService:     domainService,
		referenceService:  referenceService,
		eventHandler:      eventHandler,
	}
}
// CreateBinding crea una nueva vinculación entre un dispositivo, usuario y hogar.
// Valida que el usuario y hogar existan, verifica que el dispositivo sea válido,
// crea la vinculación y publica un evento de dispositivo vinculado.
//
// Parámetros:
//   - ctx: contexto para gestionar la cancelación y timeouts
//   - command: comando que contiene DeviceID, UserID y HomeID
//
// Retorna:
//   - *entities.DeviceBinding: la vinculación creada
//   - error: si algo falla durante el proceso
func (s *DeviceBindingCommandService) CreateBinding(ctx context.Context, command commands.CreateDeviceBindingCommand) (*entities.DeviceBinding, error) {
	if err := s.referenceService.ValidateUserReference(ctx, command.UserID); err != nil {
		return nil, err
	}
	if err := s.referenceService.ValidateHomeReference(ctx, command.HomeID); err != nil {
		return nil, err
	}
	device, err := s.deviceRepository.FindByID(ctx, command.DeviceID)
	if err != nil {
		return nil, err
	}
	if err := s.domainService.ValidateBinding(device); err != nil {
		return nil, err
	}
	binding, err := entities.NewDeviceBinding(command.DeviceID, command.UserID, command.HomeID)
	if err != nil {
		return nil, err
	}
	if err := s.bindingRepository.Save(ctx, binding); err != nil {
		return nil, err
	}
	s.eventHandler.Publish(ctx, s.eventHandler.Topics().DeviceLinked, eventhandlers.EventTypeDeviceLinked, binding.DeviceID, binding.UserID, map[string]interface{}{
		"bindingId": binding.BindingID,
		"homeId":    binding.HomeID,
	})
	return binding, nil
}
// UnlinkBinding desvincula un dispositivo de un usuario y hogar.
// Obtiene la vinculación, ejecuta la lógica de negocio de desenlace,
// actualiza el repositorio y publica un evento informativo.
//
// Parámetros:
//   - ctx: contexto para gestionar la cancelación y timeouts
//   - command: comando que contiene el BindingID a desenlazar
//
// Retorna:
//   - *entities.DeviceBinding: la vinculación actualizada
//   - error: si algo falla durante el proceso
func (s *DeviceBindingCommandService) UnlinkBinding(ctx context.Context, command commands.UnlinkDeviceBindingCommand) (*entities.DeviceBinding, error) {
	binding, err := s.bindingRepository.FindByID(ctx, command.BindingID)
	if err != nil {
		return nil, err
	}
	if err := s.domainService.UnlinkBinding(binding); err != nil {
		return nil, err
	}
	if err := s.bindingRepository.Update(ctx, binding); err != nil {
		return nil, err
	}
	s.eventHandler.Publish(ctx, s.eventHandler.Topics().DeviceUnlinked, eventhandlers.EventTypeDeviceUnlinked, binding.DeviceID, binding.UserID, map[string]interface{}{
		"bindingId":  binding.BindingID,
		"unlinkedAt": binding.UnlinkedAt,
	})
	return binding, nil
}
