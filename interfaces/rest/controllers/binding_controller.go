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
	"github.com/google/uuid"
)

type BindingController struct {
	commandService *commandservices.DeviceBindingCommandService
	queryService   *queryservices.DeviceBindingQueryService
}

func NewBindingController(commandService *commandservices.DeviceBindingCommandService, queryService *queryservices.DeviceBindingQueryService) *BindingController {
	return &BindingController{commandService: commandService, queryService: queryService}
}

func (c *BindingController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/devices/:deviceId/bindings", c.CreateBinding)
	router.GET("/devices/:deviceId/bindings", c.GetBindingsByDeviceID)
	router.GET("/users/:userId/bindings", c.GetBindingsByUserID)
	router.PATCH("/bindings/:bindingId/unlink", c.UnlinkBinding)
}

func (c *BindingController) CreateBinding(ctx *gin.Context) {
	deviceID, err := parseUUID(ctx.Param("deviceId"), "deviceId")
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	var resource resources.CreateDeviceBindingResource
	if err := ctx.ShouldBindJSON(&resource); err != nil {
		sharedinterfaces.RespondValidation(ctx, err.Error())
		return
	}
	userID, err := parseUUID(resource.UserID, "userId")
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	var homeID *uuid.UUID
	if resource.HomeID != nil {
		parsedHomeID, err := parseUUID(*resource.HomeID, "homeId")
		if err != nil {
			sharedinterfaces.RespondError(ctx, err)
			return
		}
		homeID = &parsedHomeID
	}
	binding, err := c.commandService.CreateBinding(ctx.Request.Context(), commands.CreateDeviceBindingCommand{
		DeviceID: deviceID,
		UserID:   userID,
		HomeID:   homeID,
	})
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, transform.ToBindingResource(binding))
}

func (c *BindingController) GetBindingsByDeviceID(ctx *gin.Context) {
	deviceID, err := parseUUID(ctx.Param("deviceId"), "deviceId")
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	bindings, err := c.queryService.GetByDeviceID(ctx.Request.Context(), queries.GetDeviceBindingsQuery{DeviceID: deviceID})
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToBindingResources(bindings))
}

func (c *BindingController) GetBindingsByUserID(ctx *gin.Context) {
	userID, err := parseUUID(ctx.Param("userId"), "userId")
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	bindings, err := c.queryService.GetByUserID(ctx.Request.Context(), queries.GetUserBindingsQuery{UserID: userID})
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToBindingResources(bindings))
}

func (c *BindingController) UnlinkBinding(ctx *gin.Context) {
	bindingID, err := parseUUID(ctx.Param("bindingId"), "bindingId")
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	binding, err := c.commandService.UnlinkBinding(ctx.Request.Context(), commands.UnlinkDeviceBindingCommand{BindingID: bindingID})
	if err != nil {
		sharedinterfaces.RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToBindingResource(binding))
}
