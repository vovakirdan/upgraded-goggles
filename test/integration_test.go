package tests

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestSwaggerEndpoint(t *testing.T) {
	// Создаем временную директорию для swagger файлов
	dir, err := os.MkdirTemp("", "swagger")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)

	// Задаем контент для тестового файла index.html
	dummyContent := []byte("<html><body>Swagger UI</body></html>")
	filePath := dir + "/index.html"
	if err := os.WriteFile(filePath, dummyContent, 0644); err != nil {
		t.Fatalf("failed to write dummy swagger file: %v", err)
	}

	// Создаем обработчик FileServer, отдающий файлы из временной директории
	handler := http.StripPrefix("/swagger/", http.FileServer(http.Dir(dir)))

	// Формируем тестовый HTTP-запрос
	req := httptest.NewRequest("GET", "/swagger/index.html", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	// Проверяем статус код
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	// Считываем тело ответа
	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}

	// Выводим полученное тело для отладки
	t.Logf("Got body: %q", string(body))

	// Сравниваем полученный ответ с ожидаемым содержимым
	if string(body) != string(dummyContent) {
		t.Errorf("expected %q, got %q", dummyContent, body)
	}
}
