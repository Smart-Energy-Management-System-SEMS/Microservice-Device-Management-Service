package configuration

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	persistencemodel "device-management-service/device-management/infrastructure/persistence/gorm/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDatabase(databaseURL string) (*gorm.DB, error) {
	trimmedURL := strings.TrimSpace(databaseURL)
	if trimmedURL == "" {
		return nil, errors.New("DATABASE_URL is required. Set it in .env with your real Neon PostgreSQL connection string")
	}
	if containsPlaceholder(trimmedURL) {
		return nil, errors.New("DATABASE_URL still contains placeholder values. Replace USER, PASSWORD, HOST, PORT and DB_NAME with the real Neon connection string")
	}
	parsedURL, err := url.Parse(trimmedURL)
	if err != nil {
		return nil, fmt.Errorf("DATABASE_URL must be a valid PostgreSQL URL: %w", err)
	}
	if parsedURL.Scheme != "postgresql" && parsedURL.Scheme != "postgres" {
		return nil, errors.New("DATABASE_URL must start with postgresql:// or postgres://")
	}
	if parsedURL.Host == "" || strings.Trim(parsedURL.Path, "/") == "" {
		return nil, errors.New("DATABASE_URL must include host and database name")
	}

	db, err := gorm.Open(postgres.Open(trimmedURL), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func containsPlaceholder(databaseURL string) bool {
	placeholders := []string{"USER", "PASSWORD", "HOST", "PORT", "DB_NAME"}
	for _, placeholder := range placeholders {
		if strings.Contains(databaseURL, placeholder) {
			return true
		}
	}
	return false
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&persistencemodel.DeviceModel{},
		&persistencemodel.DeviceBindingModel{},
		&persistencemodel.DeviceEventModel{},
		&persistencemodel.DeviceConfigurationModel{},
	)
}
