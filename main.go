package main

import (
	"log"

	"device-management-service/device-management/application/commandservices"
	"device-management-service/device-management/application/eventhandlers"
	"device-management-service/device-management/application/queryservices"
	"device-management-service/device-management/domain/services"
	appconfiguration "device-management-service/device-management/infrastructure/configuration"
	kafkamessaging "device-management-service/device-management/infrastructure/messaging/kafka"
	gormconfiguration "device-management-service/device-management/infrastructure/persistence/gorm/configuration"
	gormrepositories "device-management-service/device-management/infrastructure/persistence/gorm/repositories"
	"device-management-service/device-management/interfaces/rest"
	"device-management-service/device-management/interfaces/rest/controllers"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println(".env file not found, using environment variables")
	}

	config := appconfiguration.LoadAppConfig()

	db, err := gormconfiguration.ConnectDatabase(config.DatabaseURL)
	if err != nil {
		log.Fatalf("database connection failed: %v", err)
	}
	if err := gormconfiguration.AutoMigrate(db); err != nil {
		log.Fatalf("database migration failed: %v", err)
	}

	deviceRepository := gormrepositories.NewDeviceGormRepository(db)
	bindingRepository := gormrepositories.NewDeviceBindingGormRepository(db)
	configurationRepository := gormrepositories.NewDeviceConfigurationGormRepository(db)
	eventRepository := gormrepositories.NewDeviceEventGormRepository(db)

	kafkaProducer := kafkamessaging.NewProducer(config.KafkaBrokers, config.KafkaClientID)
	integrationEventHandler := eventhandlers.NewDeviceIntegrationEventHandler(kafkaProducer)
	deviceDomainService := services.NewDeviceDomainService()

	deviceCommandService := commandservices.NewDeviceCommandService(deviceRepository, integrationEventHandler)
	bindingCommandService := commandservices.NewDeviceBindingCommandService(deviceRepository, bindingRepository, deviceDomainService, integrationEventHandler)
	configurationCommandService := commandservices.NewDeviceConfigurationCommandService(deviceRepository, configurationRepository, deviceDomainService, integrationEventHandler)
	eventCommandService := commandservices.NewDeviceEventCommandService(deviceRepository, eventRepository, integrationEventHandler)

	deviceQueryService := queryservices.NewDeviceQueryService(deviceRepository)
	bindingQueryService := queryservices.NewDeviceBindingQueryService(bindingRepository)
	configurationQueryService := queryservices.NewDeviceConfigurationQueryService(configurationRepository)
	eventQueryService := queryservices.NewDeviceEventQueryService(eventRepository)

	router := rest.NewRouter(rest.RouterDependencies{
		DeviceController:        controllers.NewDeviceController(deviceCommandService, deviceQueryService),
		BindingController:       controllers.NewBindingController(bindingCommandService, bindingQueryService),
		ConfigurationController: controllers.NewConfigurationController(configurationCommandService, configurationQueryService),
		EventController:         controllers.NewEventController(eventCommandService, eventQueryService),
		CORSAllowedOrigins:      config.CORSAllowedOrigins,
	})

	log.Printf("device-management-service running on port %s", config.Port)
	if err := router.Run(":" + config.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
