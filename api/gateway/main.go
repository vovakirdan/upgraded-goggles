package main

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"user-post-system/api/gateway/config"
	"user-post-system/api/gateway/middleware"
	"user-post-system/api/gateway/routes"
)

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Загрузка конфигурации
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Создание gRPC Gateway mux
	mux := runtime.NewServeMux()

	// Опции для подключения к gRPC сервисам
	opts := []grpc.DialOption{grpc.WithInsecure()} // Для разработки, без TLS

	// Регистрация маршрутов gRPC сервисов
	err = routes.RegisterRoutes(ctx, mux, cfg, opts)
	if err != nil {
		log.Fatalf("failed to register routes: %v", err)
	}

	// Обслуживание статических файлов Swagger (доступно по /swagger/)
	http.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger"))))

	// Оборачиваем mux авторизационным middleware
	handler := middleware.AuthMiddleware(mux)

	log.Printf("Starting API Gateway at %s", cfg.HTTPPort)
	if err := http.ListenAndServe(cfg.HTTPPort, handler); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
