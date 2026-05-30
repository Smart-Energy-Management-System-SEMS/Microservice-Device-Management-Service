package rest

import (
	httpinfrastructure "device-management-service/device-management/infrastructure/http"
	"device-management-service/device-management/interfaces/rest/controllers"
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
	router.Use(httpinfrastructure.CORSMiddleware(dependencies.CORSAllowedOrigins))

	health := router.Group("/api/v1")
	health.GET("/health", controllers.Health)

	api := router.Group("/api/v1/device-management")
	api.GET("/health", controllers.Health)

	dependencies.DeviceController.RegisterRoutes(api)
	dependencies.BindingController.RegisterRoutes(api)
	dependencies.ConfigurationController.RegisterRoutes(api)
	dependencies.EventController.RegisterRoutes(api)

	return router
}
