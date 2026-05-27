package rest

import (
	"device-management-service/interfaces/rest/controllers"
	sharedinfrastructure "device-management-service/shared/infrastructure"
	sharedinterfaces "device-management-service/shared/interfaces"
	"github.com/gin-gonic/gin"
)

type RouterDependencies struct {
	DeviceController        *controllers.DeviceController
	BindingController       *controllers.BindingController
	ConfigurationController *controllers.ConfigurationController
	EventController         *controllers.EventController
	CORSAllowedOrigins      []string
}

func NewRouter(dependencies RouterDependencies) *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(sharedinfrastructure.CORSMiddleware(dependencies.CORSAllowedOrigins))

	api := router.Group("/api/v1/device-management")
	api.GET("/health", sharedinterfaces.Health)

	dependencies.DeviceController.RegisterRoutes(api)
	dependencies.BindingController.RegisterRoutes(api)
	dependencies.ConfigurationController.RegisterRoutes(api)
	dependencies.EventController.RegisterRoutes(api)

	return router
}
