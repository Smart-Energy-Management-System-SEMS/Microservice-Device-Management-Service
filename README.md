# Device Management Service

Microservicio Go para SEMS que gestiona dispositivos, vinculaciones, configuraciones y eventos con arquitectura DDD.

## Integración local con Gateway + Config Service

Entorno objetivo local:
- Config Service: `http://localhost:8090`
- API Gateway: `http://localhost:8081`
- Microservicio: `http://localhost:8083`

`base_url_local` final:
- `http://localhost:8083`

`route_prefix`:
- `/api/v1/device-management`

## Variables de entorno locales mínimas

```env
PORT=8083
SERVICE_NAME=device-management-service
CONFIG_SERVICE_URL=http://localhost:8090
DATABASE_URL=postgresql://USER:PASSWORD@HOST:PORT/DB_NAME?sslmode=require
API_GATEWAY_ALLOWED_ORIGIN=http://localhost:8081
CORS_ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,http://localhost:8081
```

Sensibles:
- `DATABASE_URL`

No sensibles:
- `PORT`
- `SERVICE_NAME`
- `CONFIG_SERVICE_URL`
- `API_GATEWAY_ALLOWED_ORIGIN`
- `CORS_ALLOWED_ORIGINS`

## Configuración remota desde Config Service

Este servicio consulta:

```http
GET {CONFIG_SERVICE_URL}/api/v1/config/{service-name}
```

Si Config Service no responde, usa fallback local sin romper el arranque.

## Health check

Endpoint público sin autenticación:

```http
GET /api/v1/device-management/health
```

## Endpoints reales (verificados en código)

```http
POST   /api/v1/device-management/devices
GET    /api/v1/device-management/devices
GET    /api/v1/device-management/devices/:deviceId
GET    /api/v1/device-management/users/:userId/devices
PUT    /api/v1/device-management/devices/:deviceId
PATCH  /api/v1/device-management/devices/:deviceId/status
DELETE /api/v1/device-management/devices/:deviceId
POST   /api/v1/device-management/devices/:deviceId/bindings
GET    /api/v1/device-management/devices/:deviceId/bindings
GET    /api/v1/device-management/users/:userId/bindings
PATCH  /api/v1/device-management/bindings/:bindingId/unlink
POST   /api/v1/device-management/devices/:deviceId/configurations
GET    /api/v1/device-management/devices/:deviceId/configurations
PUT    /api/v1/device-management/configurations/:configurationId
POST   /api/v1/device-management/devices/:deviceId/events
GET    /api/v1/device-management/devices/:deviceId/events
```

## Auth/JWT

Este microservicio no aplica middleware JWT propio. La autenticación/autorización se delega al API Gateway.
- Si `API_GATEWAY_AUTH_REQUIRED=false` en Gateway: se puede probar sin token.
- Endpoint público recomendado siempre: `GET /api/v1/device-management/health`.

## Dependencias locales

- PostgreSQL accesible con `DATABASE_URL`.
- Kafka opcional para publicación de eventos.

Kafka local (ejemplo):
```bash
docker compose up -d
```

## Pruebas mínimas

Health del microservicio:
```bash
curl -i http://localhost:8083/api/v1/device-management/health
```

Endpoint principal del microservicio:
```bash
curl -i http://localhost:8083/api/v1/device-management/devices
```

Endpoint vía Gateway (proxied):
```bash
curl -i http://localhost:8081/api/v1/device-management/health
```

## Registro para Config Service

Objeto listo para `GET/POST` de servicios (sin secretos):

```json
{
  "name": "device-management-service",
  "base_url_local": "http://localhost:8083",
  "base_url_deploy": "https://device-management-service.<tu-dominio>",
  "route_prefix": "/api/v1/device-management",
  "main_endpoints": [
    "GET /api/v1/device-management/health",
    "POST /api/v1/device-management/devices",
    "GET /api/v1/device-management/devices",
    "PATCH /api/v1/device-management/devices/:deviceId/status",
    "POST /api/v1/device-management/devices/:deviceId/events"
  ]
}
```
