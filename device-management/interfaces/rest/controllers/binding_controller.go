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
	"github.com/google/uuid"
)
// BindingController gestiona las operaciones REST relacionadas con la vinculación de dispositivos a usuarios.
// Implementa un patrón de controlador que coordina entre servicios de comandos (escritura) y servicios
// de consultas (lectura) para mantener una separación clara entre operaciones de lectura y escritura (CQRS).

type BindingController struct {
	commandService *commandservices.DeviceBindingCommandService
	queryService   *queryservices.DeviceBindingQueryService
}

func NewBindingController(commandService *commandservices.DeviceBindingCommandService, queryService *queryservices.DeviceBindingQueryService) *BindingController {
	return &BindingController{commandService: commandService, queryService: queryService}
}
// RegisterRoutes registra todas las rutas HTTP asociadas con la vinculación de dispositivos.
// Mapea los endpoints REST a sus respectivos manejadores en el controlador.
// Parámetros:
//   - router: grupo de rutas Gin donde se registrarán los endpoints

func (c *BindingController) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/devices/:deviceId/bindings", c.CreateBinding)
	router.GET("/devices/:deviceId/bindings", c.GetBindingsByDeviceID)
	router.GET("/users/:userId/bindings", c.GetBindingsByUserID)
	router.PATCH("/bindings/:bindingId/unlink", c.UnlinkBinding)
}
// CreateBinding maneja la creación de una nueva vinculación entre un dispositivo y un usuario.
// Realiza validaciones de entrada, convierte el recurso REST en un comando de dominio,
// lo ejecuta a través del servicio de comandos, y retorna la vinculación creada.
// Parámetros:
//   - ctx: contexto de Gin que contiene los parámetros y cuerpo de la solicitud

func (c *BindingController) CreateBinding(ctx *gin.Context) {
	deviceID, err := parseUUID(ctx.Param("deviceId"), "deviceId")
	if err != nil {
		RespondError(ctx, err)
		return
	}
	var resource resources.CreateDeviceBindingResource
	if err := ctx.ShouldBindJSON(&resource); err != nil {
		RespondValidation(ctx, err.Error())
		return
	}
	userID, err := parseUUID(resource.UserID, "userId")
	if err != nil {
		RespondError(ctx, err)
		return
	}
	var homeID *uuid.UUID
	if resource.HomeID != nil {
		parsedHomeID, err := parseUUID(*resource.HomeID, "homeId")
		if err != nil {
			RespondError(ctx, err)
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
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, transform.ToBindingResource(binding))
}
// GetBindingsByDeviceID obtiene todas las vinculaciones asociadas a un dispositivo específico.
// Realiza una consulta de solo lectura para recuperar todas las asociaciones usuario-dispositivo
// para un dispositivo determinado.
// Parámetros:
//   - ctx: contexto de Gin que contiene el ID del dispositivo en los parámetros de la URL

func (c *BindingController) GetBindingsByDeviceID(ctx *gin.Context) {
	deviceID, err := parseUUID(ctx.Param("deviceId"), "deviceId")
	if err != nil {
		RespondError(ctx, err)
		return
	}
	bindings, err := c.queryService.GetByDeviceID(ctx.Request.Context(), queries.GetDeviceBindingsQuery{DeviceID: deviceID})
	if err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToBindingResources(bindings))
}
// GetBindingsByUserID obtiene todas las vinculaciones asociadas a un usuario específico.
// Realiza una consulta de solo lectura para recuperar todos los dispositivos vinculados
// a un usuario determinado.
// Parámetros:
//   - ctx: contexto de Gin que contiene el ID del usuario en los parámetros de la URL

func (c *BindingController) GetBindingsByUserID(ctx *gin.Context) {
	userID, err := parseUUID(ctx.Param("userId"), "userId")
	if err != nil {
		RespondError(ctx, err)
		return
	}
	bindings, err := c.queryService.GetByUserID(ctx.Request.Context(), queries.GetUserBindingsQuery{UserID: userID})
	if err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToBindingResources(bindings))
}

func (c *BindingController) UnlinkBinding(ctx *gin.Context) {
	bindingID, err := parseUUID(ctx.Param("bindingId"), "bindingId")
	if err != nil {
		RespondError(ctx, err)
		return
	}
	binding, err := c.commandService.UnlinkBinding(ctx.Request.Context(), commands.UnlinkDeviceBindingCommand{BindingID: bindingID})
	if err != nil {
		RespondError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, transform.ToBindingResource(binding))
}
