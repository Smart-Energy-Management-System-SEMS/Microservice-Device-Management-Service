# Device Management Service

Microservicio Go para SEMS que gestiona dispositivos, vinculaciones, configuraciones y eventos con arquitectura DDD.

## Configuración centralizada

Este servicio ahora prioriza configuración remota desde un Config Service vía `CONFIG_SERVICE_URL`.

Flujo de carga:
1. Carga variables locales mínimas (`.env`/entorno).
2. Consulta `GET {CONFIG_SERVICE_URL}/api/v1/config/{service-name}`.
3. Fusiona la respuesta remota sobre defaults/fallback locales.
4. Mantiene secretos solo localmente (por ejemplo `DATABASE_URL`).

Si el Config Service no responde, el servicio usa fallbacks locales para no romper ejecución.

## Variables de entorno locales (mínimas)

```env
PORT=8083
SERVICE_NAME=device-management-service
CONFIG_SERVICE_URL=http://localhost:8081
DATABASE_URL=postgresql://USER:PASSWORD@HOST:PORT/DB_NAME?sslmode=require
```

Notas:
- `DATABASE_URL` es secreto y no debe publicarse desde Config Service.
- El resto de configuración compartida (Kafka topics, brokers, CORS, flags de runtime) debe venir del Config Service.

## Contrato esperado del Config Service

Endpoint consumido por este microservicio:

```http
GET /api/v1/config/{service-name}
```

Campos esperados (JSON):

```json
{
  "appEnv": "production",
  "dbDriver": "postgres",
  "autoMigrate": false,
  "kafkaEnabled": true,
  "kafkaBrokers": ["broker-1:9092"],
  "kafkaClientId": "device-management-service",
  "kafkaConsumerGroup": "device-management-group",
  "kafkaWriteTimeoutMs": 2000,
  "apiGatewayAllowedOrigin": "https://api.example.com",
  "corsAllowedOrigins": ["https://api.example.com", "https://web.example.com"],
  "kafkaTopics": {
    "deviceRegistered": "device.registered",
    "deviceStatusUpdated": "device.status.updated",
    "deviceLinked": "device.linked",
    "deviceUnlinked": "device.unlinked",
    "deviceConfigurationUpdated": "device.configuration.updated",
    "deviceEventRecorded": "device.event.recorded"
  }
}
```

## Ejecución local

```bash
go mod tidy
go run .
```

Base URL local:

```text
http://localhost:8083/api/v1/device-management
```

Health:

```http
GET /api/v1/device-management/health
```

## Endpoints REST (sin cambios)

### Devices

```http
POST   /api/v1/device-management/devices
GET    /api/v1/device-management/devices
GET    /api/v1/device-management/devices/:deviceId
GET    /api/v1/device-management/users/:userId/devices
PUT    /api/v1/device-management/devices/:deviceId
PATCH  /api/v1/device-management/devices/:deviceId/status
DELETE /api/v1/device-management/devices/:deviceId
```

### Bindings

```http
POST  /api/v1/device-management/devices/:deviceId/bindings
GET   /api/v1/device-management/devices/:deviceId/bindings
GET   /api/v1/device-management/users/:userId/bindings
PATCH /api/v1/device-management/bindings/:bindingId/unlink
```

### Configurations

```http
POST /api/v1/device-management/devices/:deviceId/configurations
GET  /api/v1/device-management/devices/:deviceId/configurations
PUT  /api/v1/device-management/configurations/:configurationId
```

### Events

```http
POST /api/v1/device-management/devices/:deviceId/events
GET  /api/v1/device-management/devices/:deviceId/events
```

## Docker

```bash
docker compose up --build
```

## Azure Container Apps (recomendado)

Configurar en Container App:
- `PORT`
- `SERVICE_NAME`
- `CONFIG_SERVICE_URL`
- `DATABASE_URL` (como secreto)

Recomendaciones:
- Inyectar `DATABASE_URL` desde Azure Key Vault o secretos de Container Apps.
- Exponer Config Service por DNS interno (ingress interno) para llamadas privadas.
- Versionar configuración en Config Service por entorno (`dev`, `qa`, `prod`).
- Definir liveness/readiness probe sobre `/api/v1/device-management/health`.

## Arquitectura

Se mantiene DDD con capas:
- `device-management/domain`
- `device-management/application`
- `device-management/infrastructure`
- `device-management/interfaces`
