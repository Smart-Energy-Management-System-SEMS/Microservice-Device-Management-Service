package repositories

import (
	"context"
	"errors"

	"device-management-service/device-management/domain/model/aggregates"
	domainrepositories "device-management-service/device-management/domain/repositories"
	persistencemodel "device-management-service/device-management/infrastructure/persistence/gorm/model"
	dmerrors "device-management-service/device-management/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceGormRepository struct {
	db *gorm.DB
}

func NewDeviceGormRepository(db *gorm.DB) domainrepositories.DeviceRepository {
	return &DeviceGormRepository{db: db}
}

func (r *DeviceGormRepository) Save(ctx context.Context, device *aggregates.Device) error {
	model := toDeviceModel(device)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return dmerrors.NewConflictError("device could not be saved")
	}
	return nil
}

func (r *DeviceGormRepository) Update(ctx context.Context, device *aggregates.Device) error {
	model := toDeviceModel(device)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return dmerrors.NewInternalError("device could not be updated")
	}
	return nil
}

func (r *DeviceGormRepository) FindByID(ctx context.Context, deviceID uuid.UUID) (*aggregates.Device, error) {
	var model persistencemodel.DeviceModel
	err := r.db.WithContext(ctx).First(&model, "device_id = ?", deviceID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, dmerrors.NewNotFoundError("device not found")
	}
	if err != nil {
		return nil, dmerrors.NewInternalError("device could not be retrieved")
	}
	return toDeviceDomain(model), nil
}

func (r *DeviceGormRepository) FindByExternalDeviceCode(ctx context.Context, code string) (*aggregates.Device, error) {
	var model persistencemodel.DeviceModel
	err := r.db.WithContext(ctx).First(&model, "external_device_code = ?", code).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, dmerrors.NewNotFoundError("device not found")
	}
	if err != nil {
		return nil, dmerrors.NewInternalError("device could not be retrieved")
	}
	return toDeviceDomain(model), nil
}

func (r *DeviceGormRepository) FindAll(ctx context.Context) ([]aggregates.Device, error) {
	var models []persistencemodel.DeviceModel
	if err := r.db.WithContext(ctx).Order("registered_at DESC").Find(&models).Error; err != nil {
		return nil, dmerrors.NewInternalError("devices could not be retrieved")
	}
	return toDeviceDomains(models), nil
}

func (r *DeviceGormRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]aggregates.Device, error) {
	var models []persistencemodel.DeviceModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("registered_at DESC").Find(&models).Error; err != nil {
		return nil, dmerrors.NewInternalError("devices could not be retrieved")
	}
	return toDeviceDomains(models), nil
}

func toDeviceDomains(models []persistencemodel.DeviceModel) []aggregates.Device {
	devices := make([]aggregates.Device, 0, len(models))
	for _, model := range models {
		devices = append(devices, *toDeviceDomain(model))
	}
	return devices
}
