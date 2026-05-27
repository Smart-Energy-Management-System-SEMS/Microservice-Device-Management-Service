package repositories

import (
	"context"
	"errors"

	"device-management-service/device-management/domain/model/entities"
	domainrepositories "device-management-service/device-management/domain/repositories"
	persistencemodel "device-management-service/device-management/infrastructure/persistence/gorm/model"
	shared "device-management-service/shared/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceBindingGormRepository struct {
	db *gorm.DB
}

func NewDeviceBindingGormRepository(db *gorm.DB) domainrepositories.DeviceBindingRepository {
	return &DeviceBindingGormRepository{db: db}
}

func (r *DeviceBindingGormRepository) Save(ctx context.Context, binding *entities.DeviceBinding) error {
	model := toBindingModel(binding)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return shared.NewInternalError("binding could not be saved")
	}
	return nil
}

func (r *DeviceBindingGormRepository) Update(ctx context.Context, binding *entities.DeviceBinding) error {
	model := toBindingModel(binding)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return shared.NewInternalError("binding could not be updated")
	}
	return nil
}

func (r *DeviceBindingGormRepository) FindByID(ctx context.Context, bindingID uuid.UUID) (*entities.DeviceBinding, error) {
	var model persistencemodel.DeviceBindingModel
	err := r.db.WithContext(ctx).First(&model, "binding_id = ?", bindingID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, shared.NewNotFoundError("binding not found")
	}
	if err != nil {
		return nil, shared.NewInternalError("binding could not be retrieved")
	}
	return toBindingDomain(model), nil
}

func (r *DeviceBindingGormRepository) FindByDeviceID(ctx context.Context, deviceID uuid.UUID) ([]entities.DeviceBinding, error) {
	var models []persistencemodel.DeviceBindingModel
	if err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).Order("linked_at DESC").Find(&models).Error; err != nil {
		return nil, shared.NewInternalError("bindings could not be retrieved")
	}
	return toBindingDomains(models), nil
}

func (r *DeviceBindingGormRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]entities.DeviceBinding, error) {
	var models []persistencemodel.DeviceBindingModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("linked_at DESC").Find(&models).Error; err != nil {
		return nil, shared.NewInternalError("bindings could not be retrieved")
	}
	return toBindingDomains(models), nil
}

func toBindingDomains(models []persistencemodel.DeviceBindingModel) []entities.DeviceBinding {
	bindings := make([]entities.DeviceBinding, 0, len(models))
	for _, model := range models {
		bindings = append(bindings, *toBindingDomain(model))
	}
	return bindings
}
