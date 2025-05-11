package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"cz.Finance/backend/configs"
	"cz.Finance/backend/utils"
)

// AuthMiddleware проверяет JWT токен в запросе
func AuthMiddleware(jwtConfig configs.JWTConfig) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", "Отсутствует токен авторизации")
				return
			}

			// Проверяем формат Bearer Token
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Некорректный формат токена", "Ожидается формат: Bearer <token>")
				return
			}

			// Проверяем валидность JWT токена
			token := parts[1]
			fmt.Printf("Получен токен: %s\n", token)
			tokenParts := strings.Split(token, ".")
			if len(tokenParts) != 3 {
				utils.RespondWithError(w, http.StatusUnauthorized, "Недействительный токен", "Некорректный формат JWT")
				return
			}

			claims, err := utils.ValidateJWT(token, jwtConfig.Secret)
			if err != nil {
				fmt.Printf("Ошибка валидации токена: %v\n", err)
				utils.RespondWithError(w, http.StatusUnauthorized, "Недействительный токен", err.Error())
				return
			}

			fmt.Printf("Токен валиден. UserID: %d, Email: %s\n", claims.UserID, claims.Email)

			// Добавляем информацию о пользователе в контекст запроса
			ctx := context.WithValue(r.Context(), utils.UserIDKey, claims.UserID)
			if claims.Email != "" {
				ctx = context.WithValue(ctx, utils.EmailKey, claims.Email)
			}

			// Передаем запрос дальше с обновленным контекстом
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RoleMiddleware проверяет роль пользователя
func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := utils.GetUserIDFromContext(r)
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Требуется авторизация", err.Error())
			return
		}

		next.ServeHTTP(w, r)
	})
}
