package middleware

import (
	"Test_App/pkg/jwt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. достаём заголовок Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "отсутствует токен авторизации"})
			c.Abort()
			return
		}

		// 2. заголовок должен быть вида "Bearer <token>"
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "неверный формат токена"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		// 3. парсим и проверяем токен
		claims, err := jwt.Parse(tokenStr, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "невалидный или истёкший токен"})
			c.Abort()
			return
		}

		// 4. кладём userID и role в контекст Gin
		// теперь любой хэндлер достанет их через c.GetInt64("userID")
		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)

		// 5. передаём управление следующему хэндлеру
		c.Next()
	}
}
