// Package rest wires the whole HTTP layer together. This file builds the Gin
// engine, plugs in the middleware, and asks every controller to register its
// routes. It is the single place that defines the "map" of the API.
package rest

import (
	httpinfrastructure "device-management-service/device-management/infrastructure/http"
	"device-management-service/device-management/interfaces/rest/controllers"
	"github.com/gin-gonic/gin"
)

// RouterDependencies bundles everything NewRouter needs. Passing one struct
// instead of many parameters keeps the function signature stable and readable,
// and makes it obvious at the call site what each value is.
type RouterDependencies struct {
	DeviceController        *controllers.DeviceController
	BindingController       *controllers.BindingController
	ConfigurationController *controllers.ConfigurationController
	EventController         *controllers.EventController
	CORSAllowedOrigins      []string
}

// NewRouter assembles and returns the configured HTTP engine.
func NewRouter(dependencies RouterDependencies) *gin.Engine {
	// gin.New() gives a bare engine (no default middleware) so we add exactly
	// the middleware we want. Middleware runs on every request in order:
	//   - Logger writes a line per request.
	//   - Recovery catches panics so one bad request cannot crash the server.
	//   - CORSMiddleware controls which web origins may call this API.
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.Use(httpinfrastructure.CORSMiddleware(dependencies.CORSAllowedOrigins))
	registerSwaggerRoutes(router)

	// Health-check endpoints used by load balancers / orchestrators (e.g.
	// Kubernetes, Docker) to know the service is alive. We expose it both at the
	// generic path and under the service prefix.
	health := router.Group("/api/v1")
	health.GET("/health", controllers.Health)

	// "router.Group" creates a shared URL prefix so every controller's routes
	// live under /api/v1/device-management without repeating that string.
	api := router.Group("/api/v1/device-management")
	api.GET("/health", controllers.Health)

	// Each controller knows its own routes, so we delegate registration to keep
	// this function from turning into one giant list.
	dependencies.DeviceController.RegisterRoutes(api)
	dependencies.BindingController.RegisterRoutes(api)
	dependencies.ConfigurationController.RegisterRoutes(api)
	dependencies.EventController.RegisterRoutes(api)

	return router
}
