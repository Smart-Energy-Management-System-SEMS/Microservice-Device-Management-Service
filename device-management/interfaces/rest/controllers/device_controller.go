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

type DeviceController struct {
	commandService *commandservices.DeviceCommandService
	queryService   *queryservices.DeviceQueryService
}

func NewDeviceController(commandService *commandservices.DeviceCommandService, queryService *queryservices.DeviceQueryService) *DeviceController {
	return &DeviceController{commandService: commandService, queryService: queryService}
}

func (c *DeviceController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/devices", c.CreateDevice)
	router.GET("/devices", c.GetDevices)
	router.GET("/devices/:deviceId", c.GetDeviceByID)
	router.GET("/users/:userId/devices", c.GetDevicesByUserID)
	router.PUT("/devices/:deviceId", c.UpdateDevice)
	router.PATCH("/devices/:deviceId/status", c.UpdateDeviceStatus)
	router.DELETE("/devices/:deviceId", c.DeleteDevice)
}

func (c *DeviceController) CreateDevice(ctx *gin.Context) {
	var resource resources.CreateDeviceResource
	if err := ctx.ShouldBindJSON(&resource); err != nil {
		RespondValidation(ctx, err.Error())
		return
	}
	userID, err := parseUUID(resource.UserID, "userId")
	if err != nil {
		RespondError(ctx, err)
		return
	}
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
	ctx.JSON(http.StatusCreated, transform.ToDeviceResource(device))
}

func (c *DeviceController) GetDevices(ctx *gin.Context) {
	devices, err := c.queryService.GetAll(ctx.Request.Context())
	if err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToDeviceResources(devices))
}

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
