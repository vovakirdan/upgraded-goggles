package main

import (
	"context"
	"log"
	"net/http"

	"upgraded-goggles/api/gateway"
	"upgraded-goggles/internal/logger"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Run запускает API Gateway
func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Загружаем конфигурацию
	cfg, err := gateway.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Инициализируем логгер, записывающий логи в файл "logs/app.log".
	if err := logger.InitLogger("logs/app.log"); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	// Создаем gRPC Gateway mux
	mux := runtime.NewServeMux()

	// Опции для подключения к gRPC-сервисам
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Регистрируем маршруты для gRPC-сервисов
	err = gateway.RegisterRoutes(ctx, mux, cfg, opts)
	if err != nil {
		log.Fatalf("failed to register routes: %v", err)
	}

	// Обслуживание статических файлов Swagger (доступно по /swagger/)
	http.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("api/gateway/swagger"))))

	// Оборачиваем mux авторизационным middleware
	handler := gateway.AuthMiddleware(mux)

	// Оборачиваем handler middleware логирования.
	handler = logger.LoggingMiddleware(handler)

	log.Printf("Starting API Gateway at %s", cfg.HTTPPort)
	if err := http.ListenAndServe(cfg.HTTPPort, handler); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
