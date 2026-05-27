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

type DeviceConfigurationGormRepository struct {
	db *gorm.DB
}

func NewDeviceConfigurationGormRepository(db *gorm.DB) domainrepositories.DeviceConfigurationRepository {
	return &DeviceConfigurationGormRepository{db: db}
}

func (r *DeviceConfigurationGormRepository) Save(ctx context.Context, configuration *entities.DeviceConfiguration) error {
	model := toConfigurationModel(configuration)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return shared.NewInternalError("configuration could not be saved")
	}
	return nil
}

func (r *DeviceConfigurationGormRepository) Update(ctx context.Context, configuration *entities.DeviceConfiguration) error {
	model := toConfigurationModel(configuration)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return shared.NewInternalError("configuration could not be updated")
	}
	return nil
}

func (r *DeviceConfigurationGormRepository) FindByID(ctx context.Context, configurationID uuid.UUID) (*entities.DeviceConfiguration, error) {
	var model persistencemodel.DeviceConfigurationModel
	err := r.db.WithContext(ctx).First(&model, "configuration_id = ?", configurationID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, shared.NewNotFoundError("configuration not found")
	}
	if err != nil {
		return nil, shared.NewInternalError("configuration could not be retrieved")
	}
	return toConfigurationDomain(model), nil
}

func (r *DeviceConfigurationGormRepository) FindByDeviceID(ctx context.Context, deviceID uuid.UUID) ([]entities.DeviceConfiguration, error) {
	var models []persistencemodel.DeviceConfigurationModel
	if err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).Order("updated_at DESC").Find(&models).Error; err != nil {
		return nil, shared.NewInternalError("configurations could not be retrieved")
	}
	configurations := make([]entities.DeviceConfiguration, 0, len(models))
	for _, model := range models {
		configurations = append(configurations, *toConfigurationDomain(model))
	}
	return configurations, nil
}
