# Device Management Service

Microservicio Go de SEMS para gestión de dispositivos con arquitectura DDD.

## Health checks

- `GET /api/v1/health`
- `GET /api/v1/device-management/health` (compatibilidad existente)

## Variables requeridas

Usa `.env.example` como base.

```env
PORT=8080
SERVICE_NAME=device-management-service
CONFIG_SERVICE_URL=
DATABASE_URL=
KAFKA_BROKERS=
KAFKA_SECURITY_PROTOCOL=
KAFKA_SASL_MECHANISM=
KAFKA_USERNAME=
KAFKA_PASSWORD=
GIN_MODE=release
```

Variables adicionales soportadas:

- `KAFKA_ENABLED` (default `false`)
- `KAFKA_CLIENT_ID`
- `KAFKA_CONSUMER_GROUP`
- `KAFKA_WRITE_TIMEOUT_MS`
- `API_GATEWAY_ALLOWED_ORIGIN`
- `CORS_ALLOWED_ORIGINS`
- `AUTO_MIGRATE`

## Compatibilidad local

Para desarrollo local con Kafka en Docker:

```env
KAFKA_BROKERS=localhost:9092
CONFIG_SERVICE_URL=http://localhost:8090
PORT=8083
```

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

1. Levanta Kafka local:

```bash
docker compose up -d
```

2. Configura `.env` con:

```env
PORT=8083
CONFIG_SERVICE_URL=http://localhost:8090
KAFKA_BROKERS=localhost:9092
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
curl -i http://localhost:8083/api/v1/health
```

## Ejemplo Azure Container Apps

Configura variables en Container App (sin `localhost`):

```text
PORT=8080
SERVICE_NAME=device-management-service
CONFIG_SERVICE_URL=https://<config-service-domain>
DATABASE_URL=<postgresql-connection-string>
KAFKA_BROKERS=<broker1:9092,broker2:9092>
KAFKA_SECURITY_PROTOCOL=SASL_SSL
KAFKA_SASL_MECHANISM=PLAIN
KAFKA_USERNAME=<kafka-username>
KAFKA_PASSWORD=<kafka-password>
GIN_MODE=release
```

Mapea el puerto de ingreso de la app a `8080`.

## Endpoints de negocio

Se mantienen sin cambios bajo el prefijo:

- `/api/v1/device-management`
