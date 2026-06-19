package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Device Management Service Swagger</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  <style>
    html { box-sizing: border-box; overflow-y: scroll; }
    *, *:before, *:after { box-sizing: inherit; }
    body { margin: 0; background: #fafafa; }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
  <script>
    window.onload = function () {
      SwaggerUIBundle({
        url: "/swagger/doc.json",
        dom_id: "#swagger-ui",
        deepLinking: true,
        docExpansion: "list",
        defaultModelsExpandDepth: 1
      });
    };
  </script>
</body>
</html>
`

func registerSwaggerRoutes(router *gin.Engine) {
	router.GET("/swagger", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	router.GET("/swagger/", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	router.GET("/swagger/index.html", func(context *gin.Context) {
		context.Data(http.StatusOK, "text/html; charset=utf-8", []byte(swaggerUIHTML))
	})
	router.GET("/swagger/doc.json", func(context *gin.Context) {
		context.JSON(http.StatusOK, swaggerDocument())
	})
}

func swaggerDocument() gin.H {
	return gin.H{
		"openapi": "3.0.3",
		"info": gin.H{
			"title":       "Device Management Service API",
			"description": "HTTP API for device, binding, configuration and event management.",
			"version":     "1.0.0",
		},
		"servers": []gin.H{
			{"url": "/"},
		},
		"paths": gin.H{
			"/api/v1/health": gin.H{
				"get": swaggerOperation("Health check", "Basic service health probe."),
			},
			"/api/v1/device-management/health": gin.H{
				"get": swaggerOperation("Scoped health check", "Compatibility health probe under the device-management prefix."),
			},
			"/api/v1/device-management/devices": gin.H{
				"get": swaggerOperation("List devices", "Returns the registered devices."),
				"post": swaggerCreateOperation(
					"Create device",
					"Creates a new device.",
					gin.H{
						"externalDeviceCode": "MED-001",
						"userId":             "9f4f7aef-b4ae-4284-a28c-a2dc2dfe4a93",
						"deviceName":         "Smart Meter Kitchen",
						"deviceType":         "meter",
						"brand":              "Itron",
						"model":              "EM100",
						"connectionProtocol": "WIFI",
					},
				),
			},
			"/api/v1/device-management/devices/{deviceId}": gin.H{
				"get": swaggerOperationWithPathParam("Get device by ID", "Returns a device by its identifier.", "deviceId", "Device identifier."),
				"put": swaggerOperationWithPathParamAndBody(
					"Update device",
					"Updates an existing device.",
					"deviceId",
					"Device identifier.",
					gin.H{
						"deviceName":         "Smart Meter Kitchen v2",
						"deviceType":         "meter",
						"brand":              "Itron",
						"model":              "EM100-Pro",
						"connectionProtocol": "WIFI",
					},
				),
				"delete": swaggerOperationWithPathParam("Delete device", "Deletes an existing device.", "deviceId", "Device identifier."),
			},
			"/api/v1/device-management/users/{userId}/devices": gin.H{
				"get": swaggerOperationWithPathParam("List devices by user", "Returns the devices linked to a user.", "userId", "User identifier."),
			},
			"/api/v1/device-management/devices/{deviceId}/bindings": gin.H{
				"get": swaggerOperationWithPathParam("List bindings by device", "Returns bindings for a device.", "deviceId", "Device identifier."),
				"post": swaggerOperationWithPathParamAndBody(
					"Create binding",
					"Creates a binding for a device.",
					"deviceId",
					"Device identifier.",
					gin.H{
						"userId": "9f4f7aef-b4ae-4284-a28c-a2dc2dfe4a93",
						"homeId": "f60a6fb0-734c-4b3e-9396-6fd6d2d21b44",
					},
				),
			},
			"/api/v1/device-management/users/{userId}/bindings": gin.H{
				"get": swaggerOperationWithPathParam("List bindings by user", "Returns bindings for a user.", "userId", "User identifier."),
			},
			"/api/v1/device-management/devices/{deviceId}/configurations": gin.H{
				"get": swaggerOperationWithPathParam("List configurations by device", "Returns configurations for a device.", "deviceId", "Device identifier."),
				"post": swaggerOperationWithPathParamAndBody(
					"Create configuration",
					"Creates a configuration for a device.",
					"deviceId",
					"Device identifier.",
					gin.H{
						"configKey":   "samplingIntervalSeconds",
						"configValue": "300",
					},
				),
			},
			"/api/v1/device-management/configurations/{configurationId}": gin.H{
				"put": swaggerOperationWithPathParamAndBody(
					"Update configuration",
					"Updates an existing configuration.",
					"configurationId",
					"Configuration identifier.",
					gin.H{
						"configValue": "600",
					},
				),
			},
			"/api/v1/device-management/devices/{deviceId}/events": gin.H{
				"get": swaggerOperationWithPathParam("List events by device", "Returns domain events recorded for a device.", "deviceId", "Device identifier."),
				"post": swaggerOperationWithPathParamAndBody(
					"Record event",
					"Records an event and publishes it when Kafka is enabled.",
					"deviceId",
					"Device identifier.",
					gin.H{
						"eventType":   "device.event.recorded",
						"description": "Manual test event from Swagger",
						"occurredAt":  "2026-06-18T04:45:00Z",
					},
				),
			},
		},
	}
}

func swaggerOperation(summary string, description string) gin.H {
	return gin.H{
		"summary":     summary,
		"description": description,
		"responses": gin.H{
			"200": gin.H{"description": "Successful response."},
		},
	}
}

func swaggerCreateOperation(summary string, description string, example gin.H) gin.H {
	operation := swaggerOperation(summary, description)
	operation["responses"] = gin.H{
		"201": gin.H{"description": "Resource created successfully."},
		"400": gin.H{"description": "Validation error or malformed JSON."},
	}
	operation["requestBody"] = swaggerJSONRequestBody(example)
	return operation
}

func swaggerOperationWithPathParam(summary string, description string, paramName string, paramDescription string) gin.H {
	operation := swaggerOperation(summary, description)
	operation["parameters"] = []gin.H{
		{
			"name":        paramName,
			"in":          "path",
			"required":    true,
			"description": paramDescription,
			"schema": gin.H{
				"type": "string",
			},
		},
	}
	return operation
}

func swaggerOperationWithPathParamAndBody(summary string, description string, paramName string, paramDescription string, example gin.H) gin.H {
	operation := swaggerOperationWithPathParam(summary, description, paramName, paramDescription)
	operation["requestBody"] = swaggerJSONRequestBody(example)
	operation["responses"] = gin.H{
		"200": gin.H{"description": "Successful response."},
		"400": gin.H{"description": "Validation error or malformed JSON."},
	}
	return operation
}

func swaggerJSONRequestBody(example gin.H) gin.H {
	return gin.H{
		"required": true,
		"content": gin.H{
			"application/json": gin.H{
				"schema": gin.H{
					"type": "object",
				},
				"example": example,
			},
		},
	}
}
