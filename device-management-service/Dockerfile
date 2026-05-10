# =============================================
# Stage 1: Build
# =============================================
FROM maven:3.9.6-eclipse-temurin-21 AS build

WORKDIR /app

# Cache dependencies first
COPY pom.xml .
RUN mvn dependency:go-offline -B

# Copy source and build
COPY src ./src
RUN mvn clean package -DskipTests

# =============================================
# Stage 2: Runtime
# =============================================
FROM eclipse-temurin:21-jre-alpine

WORKDIR /app

# Copy JAR from build stage
COPY --from=build /app/target/device-management-service-0.0.1-SNAPSHOT.jar app.jar

# Expose port (Railway uses PORT env variable)
EXPOSE 8081

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD wget -qO- http://localhost:${PORT:-8081}/actuator/health || exit 1

# Run
ENTRYPOINT ["java", "-jar", "-Dserver.port=${PORT:-8081}", "app.jar"]
