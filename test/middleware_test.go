package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"upgraded-goggles/api/gateway"
)

func TestAuthMiddleware_MissingHeader(t *testing.T) {
	handler := gateway.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestAuthMiddleware_InvalidFormat(t *testing.T) {
	handler := gateway.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

// TestAuthMiddleware_ValidToken проверяет, что middleware пускает с валидным JWT
func TestAuthMiddleware_ValidToken(t *testing.T) {
    // 1. Генерируем реальный JWT-токен
    secretKey := []byte("your-secret-key") // ! Должен совпадать с gateway.secretKey
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "exp": time.Now().Add(time.Hour).Unix(),
        "iat": time.Now().Unix(),
    })
    tokenString, _ := token.SignedString(secretKey)

    // 2. Создаем middleware с тестовым обработчиком
    handler := gateway.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    }))

    req := httptest.NewRequest("GET", "/test", nil)
    // 3. Передаем реальный токен в заголовке
    req.Header.Set("Authorization", "Bearer "+tokenString)
    rec := httptest.NewRecorder()

    // 4. Выполняем запрос
    handler.ServeHTTP(rec, req)

    if rec.Code != http.StatusOK {
        t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
    }
}
