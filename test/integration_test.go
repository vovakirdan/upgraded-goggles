package tests

import (
	"io/ioutil"
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

	dummyContent := []byte("<html><body>Swagger UI</body></html>")
	err = os.WriteFile(dir+"/index.html", dummyContent, 0644)
	if err != nil {
		t.Fatalf("failed to write dummy swagger file: %v", err)
	}

	// Создаем обработчик для статических файлов
	handler := http.StripPrefix("/swagger/", http.FileServer(http.Dir(dir)))
	req := httptest.NewRequest("GET", "/swagger/index.html", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatalf("failed to read response: %v", err)
	}
	if string(body) != string(dummyContent) {
		t.Errorf("expected %s, got %s", string(dummyContent), string(body))
	}
}
