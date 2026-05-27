# Device Management Service

Microservicio Go para el Smart Energy Management System (SEMS). Gestiona dispositivos, vinculaciones, configuraciones y eventos usando DDD, Gin, GORM, PostgreSQL Neon y Apache Kafka.

## Stack

- Go 1.26
- Gin REST framework
- GORM ORM
- PostgreSQL en Neon mediante `DATABASE_URL`
- Apache Kafka mediante `github.com/segmentio/kafka-go`
- Arquitectura DDD con capas `device-management/domain`, `device-management/application`, `device-management/infrastructure` y `device-management/interfaces`

## Variables de entorno

Copia `.env.example` a `.env` y coloca tu cadena real de Neon. No ejecutes el servicio con el placeholder `USER:PASSWORD@HOST:PORT/DB_NAME`, porque `PORT` debe ser un numero real como `5432`.

```env
PORT=8083
APP_ENV=local
DATABASE_URL=postgresql://neondb_owner:YOUR_PASSWORD@ep-example-123456.us-east-2.aws.neon.tech/neondb?sslmode=require
DB_DRIVER=postgres
AUTO_MIGRATE=false
KAFKA_BROKERS=localhost:9092
KAFKA_CLIENT_ID=device-management-service
KAFKA_CONSUMER_GROUP=device-management-group
API_GATEWAY_ALLOWED_ORIGIN=http://localhost:8080
CORS_ALLOWED_ORIGINS=http://localhost:4200,http://localhost:5173,http://localhost:8080
```

## Ejecutar localmente

```bash
go mod tidy
go run .
```

El servicio arranca por defecto en:

`AUTO_MIGRATE=false` es el valor recomendado cuando Neon ya tiene las tablas creadas. Usa `AUTO_MIGRATE=true` solo para una base vacia donde quieres que GORM cree o ajuste el esquema.

```text
http://localhost:8083/api/v1/device-management
```

Health check:

```http
GET /api/v1/device-management/health
```

## Endpoints REST

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

## Valores permitidos

`connectionProtocol`:

```text
WIFI, BLUETOOTH
```

`status`:

```text
ACTIVE, INACTIVE, DISCONNECTED, REMOVED
```

`bindingStatus`:

```text
LINKED, UNLINKED, PENDING
```

## Kafka

El servicio publica eventos JSON en estos tópicos:

```text
device.registered
device.status.updated
device.linked
device.unlinked
device.configuration.updated
device.event.recorded
```

Formato base:

```json
{
  "eventId": "2a28039d-df5b-4a53-8f09-4e1243939d25",
  "eventType": "DEVICE_REGISTERED",
  "deviceId": "b9f9a832-4d2a-4a48-95e1-8c14687d16d5",
  "userId": "0c389ba8-99ca-492b-8b7d-86c5056613f6",
  "occurredAt": "2026-05-26T22:00:00Z",
  "payload": {}
}
```

## Ejemplos JSON

### Registrar dispositivo

```json
{
  "externalDeviceCode": "METER-UPC-001",
  "userId": "0c389ba8-99ca-492b-8b7d-86c5056613f6",
  "deviceName": "Smart Meter Sala",
  "deviceType": "SMART_METER",
  "brand": "Shelly",
  "model": "Pro 3EM",
  "connectionProtocol": "WIFI"
}
```

### Actualizar dispositivo

```json
{
  "deviceName": "Smart Meter Sala Principal",
  "deviceType": "SMART_METER",
  "brand": "Shelly",
  "model": "Pro 3EM",
  "connectionProtocol": "WIFI"
}
```

### Actualizar estado

```json
{
  "status": "DISCONNECTED"
}
```

### Vincular dispositivo

```json
{
  "userId": "0c389ba8-99ca-492b-8b7d-86c5056613f6",
  "homeId": "243d53c0-c023-4e15-9d06-8f8fd618e377"
}
```

### Crear configuración

```json
{
  "configKey": "sampling_interval_seconds",
  "configValue": "60"
}
```

### Actualizar configuración

```json
{
  "configValue": "30"
}
```

### Registrar evento del dispositivo

```json
{
  "eventType": "ENERGY_READING_REPORTED",
  "description": "Lectura de consumo enviada por el medidor",
  "occurredAt": "2026-05-26T22:00:00Z"
}
```

Si `occurredAt` no se envía, el dominio asigna automáticamente la fecha actual en UTC.

## Reglas de dominio implementadas

- No se registra un dispositivo sin `externalDeviceCode`, `userId`, `deviceName`, `deviceType` ni `connectionProtocol`.
- `user_id` y `home_id` se guardan como UUID de referencia externa, sin foreign keys a otros microservicios.
- No se vincula ni se actualiza configuración de un dispositivo `REMOVED`.
- El delete es lógico: cambia el estado a `REMOVED`.
- `REMOVED` es estado final.
- Al desvincular un binding se marca `UNLINKED` y se llena `unlinkedAt`.
- Los controllers no contienen lógica de negocio; delegan a command/query services.
