package main

import (
	"log"
	"os"

	"upgraded-goggles/pkg/database"
	"upgraded-goggles/internal/user"
	grpcserver "upgraded-goggles/api/grpc"
)

func main() {
	// Получаем строку подключения к БД из переменной окружения
	dataSource := os.Getenv("DATABASE_URL")
	if dataSource == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Инициализируем подключение к базе данных
	if err := database.InitDB(dataSource); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Создаем репозиторий и сервис для пользователей
	userRepo := user.NewRepository(database.DB)
	userService := user.NewService(userRepo)

	// Создаем gRPC сервер для пользователей
	userGRPCServer := grpcserver.NewUserServer(userService)

	// Читаем порт для gRPC сервера из переменной окружения (по умолчанию :50051)
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = ":50051"
	}

	// Запускаем gRPC сервер для User Service
	grpcserver.StartUserGRPCServer(port, userGRPCServer)
}
