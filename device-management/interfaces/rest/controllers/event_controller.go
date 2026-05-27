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

type EventController struct {
	commandService *commandservices.DeviceEventCommandService
	queryService   *queryservices.DeviceEventQueryService
}

func NewEventController(commandService *commandservices.DeviceEventCommandService, queryService *queryservices.DeviceEventQueryService) *EventController {
	return &EventController{commandService: commandService, queryService: queryService}
}

func (c *EventController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/devices/:deviceId/events", c.RecordEvent)
	router.GET("/devices/:deviceId/events", c.GetEventsByDeviceID)
}

func (c *EventController) RecordEvent(ctx *gin.Context) {
	deviceID, err := parseUUID(ctx.Param("deviceId"), "deviceId")
	if err != nil {
		RespondError(ctx, err)
		return
	}
	var resource resources.CreateDeviceEventResource
	if err := ctx.ShouldBindJSON(&resource); err != nil {
		RespondValidation(ctx, err.Error())
		return
	}
	event, err := c.commandService.RecordEvent(ctx.Request.Context(), commands.RecordDeviceEventCommand{
		DeviceID:    deviceID,
		EventType:   resource.EventType,
		Description: resource.Description,
		OccurredAt:  resource.OccurredAt,
	})
	if err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, transform.ToEventResource(event))
}

func (c *EventController) GetEventsByDeviceID(ctx *gin.Context) {
	deviceID, err := parseUUID(ctx.Param("deviceId"), "deviceId")
	if err != nil {
		RespondError(ctx, err)
		return
	}
	events, err := c.queryService.GetByDeviceID(ctx.Request.Context(), queries.GetDeviceEventsQuery{DeviceID: deviceID})
	if err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToEventResources(events))
}
