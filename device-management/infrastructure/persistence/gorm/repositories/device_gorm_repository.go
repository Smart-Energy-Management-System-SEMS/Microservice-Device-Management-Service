// Package repositories (infrastructure side) contains the concrete database
// implementations of the repository interfaces declared in the domain. This one
// uses GORM, a popular Go ORM, to talk to a SQL database. It is another
// "adapter": the domain says WHAT operations exist, this file says HOW they run
// against the database.
package repositories

import (
	"context"
	"errors"

	dmerrors "device-management-service/device-management/domain"
	"device-management-service/device-management/domain/model/aggregates"
	domainrepositories "device-management-service/device-management/domain/repositories"
	persistencemodel "device-management-service/device-management/infrastructure/persistence/gorm/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// DeviceGormRepository wraps a *gorm.DB handle, which represents the database
// connection/session used to run queries.
type DeviceGormRepository struct {
	db *gorm.DB
}

// NewDeviceGormRepository returns the domain interface type
// (DeviceRepository), not the concrete struct. Returning the interface means
// callers can only use the methods the domain defined, which keeps the
// dependency pointing the right way (infrastructure depends on domain, never
// the reverse).
func NewDeviceGormRepository(db *gorm.DB) domainrepositories.DeviceRepository {
	return &DeviceGormRepository{db: db}
}

// Save inserts a brand new device row. A few important translations happen here:
//   - toDeviceModel converts the domain object into the DB-shaped struct.
//   - WithContext(ctx) ties the query to the request context (for cancellation).
//   - We hide GORM's raw error behind our own AppError, so the layers above
//     never need to know which database library we use.
func (r *DeviceGormRepository) Save(ctx context.Context, device *aggregates.Device) error {
	model := toDeviceModel(device)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		return dmerrors.NewConflictError("device could not be saved")
	}
	return nil
}

// Update writes changes to an existing device. GORM's Save updates all columns
// of the row identified by its primary key.
func (r *DeviceGormRepository) Update(ctx context.Context, device *aggregates.Device) error {
	model := toDeviceModel(device)
	if err := r.db.WithContext(ctx).Save(&model).Error; err != nil {
		return dmerrors.NewInternalError("device could not be updated")
	}
	return nil
}

// FindByID loads one device by its primary key. Notice how it carefully tells
// two failures apart:
//   - gorm.ErrRecordNotFound means "no such device" -> a NOT_FOUND error.
//   - any other error means the database itself failed -> an INTERNAL error.
// errors.Is is the idiomatic way to compare against a known sentinel error.
// The "?" placeholder is a parameterised query, which prevents SQL injection.
func (r *DeviceGormRepository) FindByID(ctx context.Context, deviceID uuid.UUID) (*aggregates.Device, error) {
	var model persistencemodel.DeviceModel
	err := r.db.WithContext(ctx).First(&model, "device_id = ?", deviceID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, dmerrors.NewNotFoundError("device not found")
	}
	if err != nil {
		return nil, dmerrors.NewInternalError("device could not be retrieved")
	}
	// On success we map the DB model back into a domain aggregate before
	// returning it, so the caller only ever deals with domain objects.
	return toDeviceDomain(model), nil
}

// FindByExternalDeviceCode looks a device up by its external code instead of its
// internal ID. It follows the exact same not-found vs. internal-error handling
// as FindByID. (This is the method the command service uses to enforce that
// external codes are unique.)
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

// FindAll returns every device, newest first. Order("registered_at DESC") adds
// an ORDER BY clause; Find loads many rows into the slice. If the query fails we
// return an internal error.
func (r *DeviceGormRepository) FindAll(ctx context.Context) ([]aggregates.Device, error) {
	var models []persistencemodel.DeviceModel
	if err := r.db.WithContext(ctx).Order("registered_at DESC").Find(&models).Error; err != nil {
		return nil, dmerrors.NewInternalError("devices could not be retrieved")
	}
	return toDeviceDomains(models), nil
}

// FindByUserID is like FindAll but adds a Where filter so only the devices that
// belong to the given user are returned.
func (r *DeviceGormRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]aggregates.Device, error) {
	var models []persistencemodel.DeviceModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("registered_at DESC").Find(&models).Error; err != nil {
		return nil, dmerrors.NewInternalError("devices could not be retrieved")
	}
	return toDeviceDomains(models), nil
}

// toDeviceDomains converts a slice of DB models into a slice of domain devices
// by mapping each element. We dereference (*toDeviceDomain(...)) because the
// mapper returns a pointer but we want values in the slice.
func toDeviceDomains(models []persistencemodel.DeviceModel) []aggregates.Device {
	devices := make([]aggregates.Device, 0, len(models))
	for _, model := range models {
		devices = append(devices, *toDeviceDomain(model))
	}
	return devices
}
