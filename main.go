package main

import (
	"log"
	"time"

	"device-management-service/device-management/application/commandservices"
	"device-management-service/device-management/application/eventhandlers"
	"device-management-service/device-management/application/outboundservices"
	"device-management-service/device-management/application/queryservices"
	"device-management-service/device-management/domain/services"
	appconfiguration "device-management-service/device-management/infrastructure/configuration"
	kafkamessaging "device-management-service/device-management/infrastructure/messaging/kafka"
	gormconfiguration "device-management-service/device-management/infrastructure/persistence/gorm/configuration"
	gormrepositories "device-management-service/device-management/infrastructure/persistence/gorm/repositories"
	"device-management-service/device-management/interfaces/acl"
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
	if config.AutoMigrate {
		if err := gormconfiguration.AutoMigrate(db); err != nil {
			log.Fatalf("database migration failed: %v", err)
		}
	} else {
		log.Println("database auto migration disabled")
	}

	deviceRepository := gormrepositories.NewDeviceGormRepository(db)
	bindingRepository := gormrepositories.NewDeviceBindingGormRepository(db)
	configurationRepository := gormrepositories.NewDeviceConfigurationGormRepository(db)
	eventRepository := gormrepositories.NewDeviceEventGormRepository(db)

	var deviceEventPublisher outboundservices.DeviceEventPublisher = kafkamessaging.NewNoopPublisher()
	if config.KafkaEnabled {
		kafkaProducer := kafkamessaging.NewProducer(config.KafkaBrokers, config.KafkaClientID, time.Duration(config.KafkaWriteTimeoutMS)*time.Millisecond)
		defer kafkaProducer.Close()
		deviceEventPublisher = kafkaProducer
		log.Printf("kafka publishing enabled with brokers=%v timeout_ms=%d", config.KafkaBrokers, config.KafkaWriteTimeoutMS)
	} else {
		log.Println("kafka publishing disabled (KAFKA_ENABLED=false)")
	}
	integrationEventHandler := eventhandlers.NewDeviceIntegrationEventHandler(deviceEventPublisher)
	deviceDomainService := services.NewDeviceDomainService()
	externalReferenceService := acl.NewLocalExternalReferenceService()

	deviceCommandService := commandservices.NewDeviceCommandService(deviceRepository, externalReferenceService, integrationEventHandler)
	bindingCommandService := commandservices.NewDeviceBindingCommandService(deviceRepository, bindingRepository, deviceDomainService, externalReferenceService, integrationEventHandler)
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
