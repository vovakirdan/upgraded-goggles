# syntax=docker/dockerfile:1

##########################
# Stage 1: Builder
FROM golang:1.23-alpine AS builder
WORKDIR /app

# Копируем файлы модулей и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код проекта
COPY . .

# Собираем бинарники для каждого сервиса
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service ./cmd/user-service/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o post-service ./cmd/post-service/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway ./cmd/api-gateway/main.go

##########################
# Stage 2: User Service Image
FROM alpine:latest AS user-service
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/user-service .
EXPOSE 50051
ENTRYPOINT ["/app/user-service"]

##########################
# Stage 3: Post Service Image
FROM alpine:latest AS post-service
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=builder /app/post-service .
EXPOSE 50052
ENTRYPOINT ["/app/post-service"]

##########################
# Stage 4: API Gateway Image
FROM alpine:latest AS api-gateway
RUN apk add --no-cache ca-certificates
WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/api-gateway . 

# Копируем папку swagger внутрь /app/api/gateway/swagger
RUN mkdir -p /app/api/gateway
COPY --from=builder /app/api/gateway/swagger /app/api/gateway/swagger

EXPOSE 8080
ENTRYPOINT ["/app/api-gateway"]
