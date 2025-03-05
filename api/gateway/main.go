package gateway

import (
	"context"
	"log"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// Run запускает API Gateway
func Run() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Загружаем конфигурацию
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Создаем gRPC Gateway mux
	mux := runtime.NewServeMux()

	// Опции для подключения к gRPC-сервисам
	opts := []grpc.DialOption{grpc.WithInsecure()}

	// Регистрируем маршруты для gRPC-сервисов
	err = RegisterRoutes(ctx, mux, cfg, opts)
	if err != nil {
		log.Fatalf("failed to register routes: %v", err)
	}

	// Обслуживание статических файлов Swagger (доступно по /swagger/)
	http.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir("./swagger"))))

	// Оборачиваем mux авторизационным middleware
	handler := AuthMiddleware(mux)

	log.Printf("Starting API Gateway at %s", cfg.HTTPPort)
	if err := http.ListenAndServe(cfg.HTTPPort, handler); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
