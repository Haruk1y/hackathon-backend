// internal/middleware/auth.go

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/Haruk1y/hackathon-backend/internal/auth"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		idToken := strings.Replace(authHeader, "Bearer ", "", 1)
		token, err := auth.VerifyIDToken(idToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// トークンが有効な場合、ユーザーIDをコンテキストに設定
		c.Set("uid", token.UID)
		c.Next()
	}
}