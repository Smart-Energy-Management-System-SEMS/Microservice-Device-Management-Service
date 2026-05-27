package controllers

import (
	"net/http"

	"device-management-service/application/commandservices"
	"device-management-service/application/queryservices"
	"device-management-service/domain/model/commands"
	"device-management-service/domain/model/queries"
	"device-management-service/interfaces/rest/resources"
	"device-management-service/interfaces/rest/transform"
	sharedinterfaces "device-management-service/shared/interfaces"
	"github.com/gin-gonic/gin"
)

type ConfigurationController struct {
	commandService *commandservices.DeviceConfigurationCommandService
	queryService   *queryservices.DeviceConfigurationQueryService
}

func NewConfigurationController(commandService *commandservices.DeviceConfigurationCommandService, queryService *queryservices.DeviceConfigurationQueryService) *ConfigurationController {
	return &ConfigurationController{commandService: commandService, queryService: queryService}
}

func (c *ConfigurationController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/devices/:deviceId/configurations", c.CreateConfiguration)
	router.GET("/devices/:deviceId/configurations", c.GetConfigurationsByDeviceID)
	router.PUT("/configurations/:configurationId", c.UpdateConfiguration)
}

func (c *ConfigurationController) CreateConfiguration(ctx *gin.Context) {
	deviceID, err := parseUUID(ctx.Param("deviceId"), "deviceId")
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	var resource resources.CreateDeviceConfigurationResource
	if err := ctx.ShouldBindJSON(&resource); err != nil {
		sharedinterfaces.RespondValidation(ctx, err.Error())
		return
	}
	configuration, err := c.commandService.CreateConfiguration(ctx.Request.Context(), commands.CreateDeviceConfigurationCommand{
		DeviceID:    deviceID,
		ConfigKey:   resource.ConfigKey,
		ConfigValue: resource.ConfigValue,
	})
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, transform.ToConfigurationResource(configuration))
}

func (c *ConfigurationController) GetConfigurationsByDeviceID(ctx *gin.Context) {
	deviceID, err := parseUUID(ctx.Param("deviceId"), "deviceId")
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	configurations, err := c.queryService.GetByDeviceID(ctx.Request.Context(), queries.GetDeviceConfigurationsQuery{DeviceID: deviceID})
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToConfigurationResources(configurations))
}

func (c *ConfigurationController) UpdateConfiguration(ctx *gin.Context) {
	configurationID, err := parseUUID(ctx.Param("configurationId"), "configurationId")
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	var resource resources.UpdateDeviceConfigurationResource
	if err := ctx.ShouldBindJSON(&resource); err != nil {
		sharedinterfaces.RespondValidation(ctx, err.Error())
		return
	}
	configuration, err := c.commandService.UpdateConfiguration(ctx.Request.Context(), commands.UpdateDeviceConfigurationCommand{
		ConfigurationID: configurationID,
		ConfigValue:     resource.ConfigValue,
	})
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToConfigurationResource(configuration))
}
