// internal/middleware/cors.go

package middleware

import (
    "os"
    "strings"
    "github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
    // 環境変数から許可するオリジンを取得
    allowedOrigins := strings.Split(os.Getenv("ALLOWED_ORIGINS"), ",")
    if len(allowedOrigins) == 0 {
        // デフォルトでVercelのドメインを許可
        allowedOrigins = []string{"https://hackathon-frontend-taupe.vercel.app"}
    }

    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")
        
        // オリジンの検証
        allowOrigin := ""
        for _, allowed := range allowedOrigins {
            if origin == allowed {
                allowOrigin = origin
                break
            }
        }

        if allowOrigin != "" {
            c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
            c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
            c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
        }

        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}