# Device Management Service

Microservicio Go de SEMS para gestion de dispositivos con arquitectura DDD.

## Health checks

- `GET /api/v1/health`
- `GET /api/v1/device-management/health` (compatibilidad existente)

## Swagger

Swagger queda habilitado sin autenticacion en:

- `GET /swagger/index.html`
- `GET /swagger/doc.json`

Para pruebas locales, abre `http://localhost:8083/swagger/index.html` o cambia el puerto segun tu `.env` activo.

## Variables requeridas

Usa `.env.example` como base.

```env
PORT=8080
SERVICE_NAME=device-management-service
APP_ENV=local
CONFIG_SERVICE_URL=
DATABASE_URL=
KAFKA_BOOTSTRAP_SERVERS=
KAFKA_BROKERS=
KAFKA_SECURITY_PROTOCOL=
KAFKA_SASL_MECHANISM=
KAFKA_USERNAME=
KAFKA_PASSWORD=
KAFKA_SASL_USERNAME=
KAFKA_SASL_PASSWORD=
KAFKA_TOPIC_DEVICE_EVENTS=device.events
TOPIC_DEVICE_EVENTS=device.events
GIN_MODE=release
```

Variables adicionales soportadas:

- `KAFKA_ENABLED` (default `false`)
- `KAFKA_AUTO_CREATE_TOPICS` (default `false`, recomendado `false` para Azure Event Hubs)
- `KAFKA_CLIENT_ID`
- `KAFKA_CONSUMER_GROUP`
- `KAFKA_WRITE_TIMEOUT_MS`
- `API_GATEWAY_ALLOWED_ORIGIN`
- `CORS_ALLOWED_ORIGINS`
- `AUTO_MIGRATE`

Tambien acepta aliases utiles para homologar otros micros:

- `KAFKA_BOOTSTRAP_SERVERS` como fallback de `KAFKA_BROKERS`
- `KAFKA_SASL_USERNAME` como fallback de `KAFKA_USERNAME`
- `KAFKA_SASL_PASSWORD` como fallback de `KAFKA_PASSWORD`
- `KAFKA_TOPIC_DEVICE_EVENTS` como fallback de `TOPIC_DEVICE_EVENTS`

Si `POST /api/v1/device-management/devices/{id}/events` falla al persistir antes de Kafka, ejecuta al menos una vez con `AUTO_MIGRATE=true` para que GORM cree `device_events` si falta en la base local. Luego puedes volverlo a `false` si prefieres manejar el esquema manualmente.

## Kafka por dominio

Este microservicio publica todos los eventos del dominio Device en un solo topic fisico:

- `device.events`

Los valores como `device.registered`, `device.status.updated` y los demas de Device son `eventType`, no topics fisicos.
El tipo real del evento viaja en el payload JSON bajo `eventType`, por ejemplo:

```json
{
  "eventId": "0a3ab694-8382-43f0-847e-8c1ae87fce75",
  "eventType": "device.registered",
  "deviceId": "9ccfa2e6-52a8-4cc3-98af-1d2fc46ba0d8",
  "userId": "9f4f7aef-b4ae-4284-a28c-a2dc2dfe4a93",
  "occurredAt": "2026-06-12T22:30:00Z",
  "payload": {
    "externalDeviceCode": "MED-001",
    "deviceType": "meter",
    "status": "ACTIVE"
  }
}
```

Eventos publicados por este micro:

- `device.registered`
- `device.linked`
- `device.unlinked`
- `device.status.updated`
- `device.configuration.updated`
- `device.event.recorded`

## Docker build

```bash
docker build -t device-management-service:local .
```

## Docker run

Ejemplo usando archivo `.env`:

```bash
docker run --name device-management-service \
  --env-file .env \
  -p 8083:8083 \
  device-management-service:local
```

Si quieres usar `PORT=8080`:

```bash
docker run --name device-management-service \
  --env-file .env \
  -e PORT=8080 \
  -p 8080:8080 \
  device-management-service:local
```

## Ejemplo local completo

Archivos locales recomendados:

- [`.env.local-kafka`](</c:/Users/ASUS/Desktop/UPC/UPC-Ciclo Vll/Fundamentos de Arquitectura de Software/Sems/Microservice-Device-Management-Service/.env.local-kafka>)
- [`.env.azure-eventhubs`](</c:/Users/ASUS/Desktop/UPC/UPC-Ciclo Vll/Fundamentos de Arquitectura de Software/Sems/Microservice-Device-Management-Service/.env.azure-eventhubs>)
- [`.env`](</c:/Users/ASUS/Desktop/UPC/UPC-Ciclo Vll/Fundamentos de Arquitectura de Software/Sems/Microservice-Device-Management-Service/.env>) como archivo activo

Cuando quieras cambiar de modo, copia el contenido del perfil deseado sobre `.env`.

1. Levanta la infraestructura:

```bash
docker compose up -d
```

2. Configura `.env` con:

```env
PORT=8083
CONFIG_SERVICE_URL=http://config-service:8090
KAFKA_BROKERS=kafka:9092
KAFKA_ENABLED=true
KAFKA_AUTO_CREATE_TOPICS=true
KAFKA_TOPIC_DEVICE_EVENTS=device.events
TOPIC_DEVICE_EVENTS=device.events
DATABASE_URL=postgresql://USER:PASSWORD@HOST:PORT/DB_NAME?sslmode=require
```

3. Ejecuta contenedor:

```bash
docker run --name device-management-service \
  --env-file .env \
  -p 8083:8083 \
  device-management-service:local
```

4. Verifica health:

```bash
curl -i http://127.0.0.1:8083/api/v1/health
```

5. Abre Swagger:

```bash
http://127.0.0.1:8083/swagger/index.html
```

## Ejemplo Azure Container Apps

Usa [`.env.azure.example`](</c:/Users/ASUS/Desktop/UPC/UPC-Ciclo Vll/Fundamentos de Arquitectura de Software/Sems/Microservice-Device-Management-Service/.env.azure.example>) como referencia.

Configura variables en Container App:

```text
PORT=8080
SERVICE_NAME=device-management-service
APP_ENV=azure
CONFIG_SERVICE_URL=https://<config-service-domain>
DATABASE_URL=<postgresql-connection-string>
KAFKA_BOOTSTRAP_SERVERS=<namespace>.servicebus.windows.net:9093
KAFKA_BROKERS=<namespace>.servicebus.windows.net:9093
KAFKA_SECURITY_PROTOCOL=SASL_SSL
KAFKA_SASL_MECHANISM=PLAIN
KAFKA_USERNAME=$ConnectionString
KAFKA_PASSWORD=Endpoint=sb://<namespace>.servicebus.windows.net/;SharedAccessKeyName=<policy>;SharedAccessKey=<key>;EntityPath=device.events
KAFKA_SASL_USERNAME=$ConnectionString
KAFKA_SASL_PASSWORD=Endpoint=sb://<namespace>.servicebus.windows.net/;SharedAccessKeyName=<policy>;SharedAccessKey=<key>;EntityPath=device.events
KAFKA_ENABLED=true
KAFKA_AUTO_CREATE_TOPICS=false
KAFKA_CONSUMER_GROUP=device-management-group
KAFKA_TOPIC_DEVICE_EVENTS=device.events
TOPIC_DEVICE_EVENTS=device.events
GIN_MODE=release
```

Mapea el puerto de ingreso de la app a `8080`.

## Endpoints de negocio

Se mantienen sin cambios bajo el prefijo:

- `/api/v1/device-management`
