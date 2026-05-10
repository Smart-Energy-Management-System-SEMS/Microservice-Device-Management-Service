# SEMS - Device Management Service

Microservicio del bounded context **Device Management** del sistema SEMS (Smart Energy Management System).

## Responsabilidades

- Vincular y desvincular dispositivos IoT al hogar del usuario
- Consultar estado y listado de dispositivos por usuario
- Gestionar preferencias de monitoreo energético

## Stack

- **Java 21** + **Spring Boot 3.2**
- **PostgreSQL** (Supabase)
- **JWT** (mismo secret que el IAM Service)
- **SpringDoc OpenAPI** (Swagger UI)

---

## Despliegue en Railway + Supabase

### 1. Crear base de datos en Supabase

1. Ir a [supabase.com](https://supabase.com) → New Project
2. Copiar la **Connection String** desde: `Settings → Database → Connection string → URI`
3. El formato será: `postgresql://postgres:[PASSWORD]@db.[REF].supabase.co:5432/postgres`
4. Convertirlo a JDBC: `jdbc:postgresql://db.[REF].supabase.co:5432/postgres`

### 2. Desplegar en Railway

1. Ir a [railway.app](https://railway.app) → New Project → Deploy from GitHub repo
2. Seleccionar este repositorio
3. Railway detecta el `Dockerfile` automáticamente
4. Configurar las variables de entorno:

```
DATABASE_URL=jdbc:postgresql://db.TUREF.supabase.co:5432/postgres
DATABASE_USERNAME=postgres
DATABASE_PASSWORD=TU_PASSWORD
JWT_SECRET=SEMSSecretKeyForJWTTokenGenerationAndValidation2025
CORS_ORIGINS=http://localhost:4200
```

5. Railway asignará la variable `PORT` automáticamente

### 3. Verificar despliegue

Una vez desplegado, acceder a:
- **Swagger UI:** `https://tu-app.railway.app/swagger-ui.html`
- **API Docs:** `https://tu-app.railway.app/v3/api-docs`
- **Health:** `https://tu-app.railway.app/actuator/health`

---

## Desarrollo local

```bash
# 1. Copiar variables de entorno
cp .env.example .env
# Editar .env con tus credenciales

# 2. Compilar
./mvnw clean package -DskipTests

# 3. Ejecutar
java -jar target/device-management-service-0.0.1-SNAPSHOT.jar

# O con Docker
docker build -t sems-device-management .
docker run -p 8081:8081 --env-file .env sems-device-management
```

Swagger disponible en: http://localhost:8081/swagger-ui.html

---

## Endpoints principales

| Método | Endpoint | Descripción |
|--------|----------|-------------|
| `POST` | `/api/v1/devices` | Crear/vincular dispositivo |
| `GET` | `/api/v1/devices` | Listar todos los dispositivos |
| `GET` | `/api/v1/devices/{id}` | Obtener dispositivo por ID |
| `GET` | `/api/v1/devices/my-devices` | Mis dispositivos (usuario autenticado) |
| `GET` | `/api/v1/devices/users/{userId}` | Dispositivos por usuario |
| `DELETE` | `/api/v1/devices/{id}` | Eliminar dispositivo |
| `GET` | `/api/v1/devices/users/{userId}/preferences` | Obtener preferencias |
| `PUT` | `/api/v1/devices/users/{userId}/preferences` | Actualizar preferencias |

## Autenticación

Todos los endpoints requieren JWT. Obtener el token desde el IAM Service (monolito) y usarlo como:

```
Authorization: Bearer <token>
```
