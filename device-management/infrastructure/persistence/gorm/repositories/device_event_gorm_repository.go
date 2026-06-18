package repositories

import (
	"context"
	"log"

	dmerrors "device-management-service/device-management/domain"
	"device-management-service/device-management/domain/model/entities"
	domainrepositories "device-management-service/device-management/domain/repositories"
	persistencemodel "device-management-service/device-management/infrastructure/persistence/gorm/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DeviceEventGormRepository struct {
	db *gorm.DB
}

func NewDeviceEventGormRepository(db *gorm.DB) domainrepositories.DeviceEventRepository {
	return &DeviceEventGormRepository{db: db}
}

func (r *DeviceEventGormRepository) Save(ctx context.Context, event *entities.DeviceEvent) error {
	model := toEventModel(event)
	if err := r.db.WithContext(ctx).Create(&model).Error; err != nil {
		log.Printf(
			"device event save failed: event_id=%s device_id=%s event_type=%s occurred_at=%s err=%v",
			model.EventID,
			model.DeviceID,
			model.EventType,
			model.OccurredAt.Format("2006-01-02T15:04:05.999999999Z07:00"),
			err,
		)
		return dmerrors.NewInternalError("event could not be saved")
	}
	return nil
}

func (r *DeviceEventGormRepository) FindByDeviceID(ctx context.Context, deviceID uuid.UUID) ([]entities.DeviceEvent, error) {
	var models []persistencemodel.DeviceEventModel
	if err := r.db.WithContext(ctx).Where("device_id = ?", deviceID).Order("occurred_at DESC").Find(&models).Error; err != nil {
		log.Printf("device event query failed: device_id=%s err=%v", deviceID, err)
		return nil, dmerrors.NewInternalError("events could not be retrieved")
	}
	events := make([]entities.DeviceEvent, 0, len(models))
	for _, model := range models {
		events = append(events, *toEventDomain(model))
	}
	return events, nil
}
