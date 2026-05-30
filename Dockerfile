# syntax=docker/dockerfile:1

FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/device-management-service .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates curl tzdata \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /out/device-management-service /app/device-management-service

ENV APP_ENV=production \
    GIN_MODE=release \
    PORT=8080 \
    AUTO_MIGRATE=false

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=20s --retries=3 \
    CMD curl -fsS "http://localhost:${PORT}/api/v1/health" || exit 1

USER appuser

ENTRYPOINT ["/app/device-management-service"]
