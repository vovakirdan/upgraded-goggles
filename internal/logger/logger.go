package logger

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Logger является глобальным логгером
var Logger *log.Logger

// InitLogger инициализирует глобальный логгер, записывающий логи в указанный файл.
// logFilePath - путь к файлу логов, например "logs/app.log"
func InitLogger(logFilePath string) error {
	// Определяем директорию из пути файла
	dir := filepath.Dir(logFilePath)
	// Создаем директорию, если её нет
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	// Открываем или создаем файл для логов
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	// Инициализируем глобальный логгер с метками времени и указанием файла
	Logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Logger.Println("Logger initialized")
	return nil
}

// LoggingMiddleware логирует входящие HTTP-запросы и время их обработки
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		// Логируем входящий запрос
		Logger.Printf("Incoming request: method=%s, path=%s, remote=%s", r.Method, r.URL.Path, r.RemoteAddr)

		// Обработка запроса
		next.ServeHTTP(w, r)

		// Логируем завершение запроса
		duration := time.Since(start)
		Logger.Printf("Completed request: method=%s, path=%s, duration=%s", r.Method, r.URL.Path, duration)
	})
}
