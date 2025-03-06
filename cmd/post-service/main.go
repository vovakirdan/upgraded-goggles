package main

import (
	"log"
	"os"

	"upgraded-goggles/pkg/database"
	"upgraded-goggles/internal/post"
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

	// Создаем репозиторий и сервис для постов
	postRepo := post.NewRepository(database.DB)
	postService := post.NewService(postRepo)

	// Создаем gRPC сервер для постов
	postGRPCServer := grpcserver.NewPostServer(postService)

	// Читаем порт для gRPC сервера из переменной окружения (по умолчанию :50052)
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = ":50052"
	}

	// Запускаем gRPC сервер для Post Service
	grpcserver.StartPostGRPCServer(port, postGRPCServer)
}
