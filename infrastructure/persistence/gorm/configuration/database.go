package configuration

import (
	"errors"

	persistencemodel "device-management-service/infrastructure/persistence/gorm/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(databaseURL string) (*gorm.DB, error) {
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL is required")
	}
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&persistencemodel.DeviceModel{},
		&persistencemodel.DeviceBindingModel{},
		&persistencemodel.DeviceEventModel{},
		&persistencemodel.DeviceConfigurationModel{},
	)
}
