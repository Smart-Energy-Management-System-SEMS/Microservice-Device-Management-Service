// Package controllers is the entry point of the interfaces layer: it exposes the
// REST API using the Gin web framework. A controller's only job is to handle the
// HTTP details — read the request, call the right application service, and turn
// the result (or error) into an HTTP response. It should contain no business
// logic; that lives in the domain and application layers.
package controllers

import (
	"net/http"

	"device-management-service/device-management/application/commandservices"
	"device-management-service/device-management/application/queryservices"
	"device-management-service/device-management/domain/model/commands"
	"device-management-service/device-management/domain/model/queries"
	"device-management-service/device-management/interfaces/rest/resources"
	"device-management-service/device-management/interfaces/rest/transform"
	"github.com/gin-gonic/gin"
)

// DeviceController depends on the two application services that match the CQRS
// split: one for writes (commandService) and one for reads (queryService).
type DeviceController struct {
	commandService *commandservices.DeviceCommandService
	queryService   *queryservices.DeviceQueryService
}

// NewDeviceController injects both services.
func NewDeviceController(commandService *commandservices.DeviceCommandService, queryService *queryservices.DeviceQueryService) *DeviceController {
	return &DeviceController{commandService: commandService, queryService: queryService}
}

// RegisterRoutes maps each URL + HTTP verb to the method that handles it. The
// choice of verb follows REST conventions: POST creates, GET reads, PUT
// replaces, PATCH partially updates, and DELETE removes. ":deviceId" is a path
// parameter that Gin fills in from the URL.
func (c *DeviceController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/devices", c.CreateDevice)
	router.GET("/devices", c.GetDevices)
	router.GET("/devices/:deviceId", c.GetDeviceByID)
	router.GET("/users/:userId/devices", c.GetDevicesByUserID)
	router.PUT("/devices/:deviceId", c.UpdateDevice)
	router.PATCH("/devices/:deviceId/status", c.UpdateDeviceStatus)
	router.DELETE("/devices/:deviceId", c.DeleteDevice)
}

// CreateDevice handles "POST /devices". It shows the typical controller flow
// that the other handlers repeat:
//   1. Bind the JSON body into a request resource (and validate its shape).
//   2. Convert/parse raw inputs such as the UUID string.
//   3. Build a command and hand it to the application service.
//   4. Translate the result into a response resource, or the error into a
//      proper HTTP status via the shared RespondError helper.
// Every step returns early on failure so the success path stays flat and clear.
func (c *DeviceController) CreateDevice(ctx *gin.Context) {
	var resource resources.CreateDeviceResource
	// ShouldBindJSON parses the request body into the struct; an error here
	// means the client sent malformed or invalid JSON.
	if err := ctx.ShouldBindJSON(&resource); err != nil {
		RespondValidation(ctx, err.Error())
		return
	}
	// The UUID arrives as text, so we parse it into a real uuid.UUID value.
	userID, err := parseUUID(resource.UserID, "userId")
	if err != nil {
		RespondError(ctx, err)
		return
	}
	// ctx.Request.Context() forwards the HTTP request's context down into the
	// service, so a cancelled/timed-out request stops work all the way down.
	device, err := c.commandService.RegisterDevice(ctx.Request.Context(), commands.RegisterDeviceCommand{
		ExternalDeviceCode: resource.ExternalDeviceCode,
		UserID:             userID,
		DeviceName:         resource.DeviceName,
		DeviceType:         resource.DeviceType,
		Brand:              resource.Brand,
		Model:              resource.Model,
		ConnectionProtocol: resource.ConnectionProtocol,
	})
	if err != nil {
		RespondError(ctx, err)
		return
	}
	// 201 Created is the correct status for a successful POST that created a
	// resource. We never return the domain object directly; transform turns it
	// into a response DTO so we control exactly which fields are exposed.
	ctx.JSON(http.StatusCreated, transform.ToDeviceResource(device))
}

// GetDevices handles "GET /devices" and returns the whole list. Read handlers
// like this one are short because they just call the query service and respond.
func (c *DeviceController) GetDevices(ctx *gin.Context) {
	devices, err := c.queryService.GetAll(ctx.Request.Context())
	if err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToDeviceResources(devices))
}

// GetDeviceByID handles "GET /devices/:deviceId". ctx.Param reads the path
// parameter from the URL. If the device is missing, the query service returns a
// NOT_FOUND error which RespondError maps to HTTP 404.
func (c *DeviceController) GetDeviceByID(ctx *gin.Context) {
	deviceID, err := parseUUID(ctx.Param("deviceId"), "deviceId")
	if err != nil {
		RespondError(ctx, err)
		return
	}
	device, err := c.queryService.GetByID(ctx.Request.Context(), queries.GetDeviceByIDQuery{DeviceID: deviceID})
	if err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToDeviceResource(device))
}

// GetDevicesByUserID handles "GET /users/:userId/devices" — all devices that
// belong to one user. The nested URL expresses the "devices of a user"
// relationship in a RESTful way.
func (c *DeviceController) GetDevicesByUserID(ctx *gin.Context) {
	userID, err := parseUUID(ctx.Param("userId"), "userId")
	if err != nil {
		RespondError(ctx, err)
		return
	}
	devices, err := c.queryService.GetByUserID(ctx.Request.Context(), queries.GetDevicesByUserQuery{UserID: userID})
	if err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToDeviceResources(devices))
}

// UpdateDevice handles "PUT /devices/:deviceId". It needs BOTH a path parameter
// (which device) and a JSON body (the new values), so it parses the ID first and
// then binds the body.
func (c *DeviceController) UpdateDevice(ctx *gin.Context) {
	deviceID, err := parseUUID(ctx.Param("deviceId"), "deviceId")
	if err != nil {
		RespondError(ctx, err)
		return
	}
	var resource resources.UpdateDeviceResource
	if err := ctx.ShouldBindJSON(&resource); err != nil {
		RespondValidation(ctx, err.Error())
		return
	}
	device, err := c.commandService.UpdateDevice(ctx.Request.Context(), commands.UpdateDeviceCommand{
		DeviceID:           deviceID,
		DeviceName:         resource.DeviceName,
		DeviceType:         resource.DeviceType,
		Brand:              resource.Brand,
		Model:              resource.Model,
		ConnectionProtocol: resource.ConnectionProtocol,
	})
	if err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToDeviceResource(device))
}

// UpdateDeviceStatus handles "PATCH /devices/:deviceId/status". PATCH is used
// because we are changing only one part of the device (its status), not
// replacing the whole resource.
func (c *DeviceController) UpdateDeviceStatus(ctx *gin.Context) {
	deviceID, err := parseUUID(ctx.Param("deviceId"), "deviceId")
	if err != nil {
		RespondError(ctx, err)
		return
	}
	var resource resources.UpdateDeviceStatusResource
	if err := ctx.ShouldBindJSON(&resource); err != nil {
		RespondValidation(ctx, err.Error())
		return
	}
	device, err := c.commandService.UpdateDeviceStatus(ctx.Request.Context(), commands.UpdateDeviceStatusCommand{DeviceID: deviceID, Status: resource.Status})
	if err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToDeviceResource(device))
}

// DeleteDevice handles "DELETE /devices/:deviceId". Remember this is a soft
// delete in the domain, so we still return the (now REMOVED) device with 200 OK
// rather than an empty 204 response.
func (c *DeviceController) DeleteDevice(ctx *gin.Context) {
	deviceID, err := parseUUID(ctx.Param("deviceId"), "deviceId")
	if err != nil {
		RespondError(ctx, err)
		return
	}
	device, err := c.commandService.DeleteDevice(ctx.Request.Context(), commands.DeleteDeviceCommand{DeviceID: deviceID})
	if err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToDeviceResource(device))
}
