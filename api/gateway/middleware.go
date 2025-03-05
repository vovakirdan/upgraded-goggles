package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// secretKey используется для подписи и проверки JWT токенов
var secretKey = []byte("your-secret-key") // TODO: change to env variable

// AuthMiddleware проверяет наличие и корректность JWT-токена в заголовке Authorization
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем заголовок Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		// Ожидается формат "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]

		// Парсим и валидируем JWT токен
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Проверяем, что используется правильный алгоритм подписи
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return secretKey, nil
		})

		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		if !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Опционально: можно получить claims из токена
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			// Добавляем claims в контекст запроса для использования в следующих обработчиках
			r = r.WithContext(AddClaimsToContext(r.Context(), claims))
		}

		// Передаем управление следующему обработчику
		next.ServeHTTP(w, r)
	})
}

// AddClaimsToContext добавляет JWT claims в контекст запроса
func AddClaimsToContext(ctx context.Context, claims jwt.MapClaims) context.Context {
	return context.WithValue(ctx, "claims", claims)
}
