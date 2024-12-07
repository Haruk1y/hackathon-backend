// cmd/api/main.go

package main

import (
    "log"
    "os"
    
    "github.com/Haruk1y/hackathon-backend/internal/auth"
    "github.com/gin-gonic/gin"
    "github.com/joho/godotenv"
    "github.com/Haruk1y/hackathon-backend/internal/database"
    "github.com/Haruk1y/hackathon-backend/internal/handler"
    "github.com/Haruk1y/hackathon-backend/internal/middleware"
)

func main() {
    // カレントディレクトリを表示
    dir, _ := os.Getwd()
    log.Printf("Current working directory: %s", dir)

    // Load .env file
    if err := godotenv.Load(); err != nil {
        log.Printf("Error loading .env file: %v", err)
    } else {
        log.Println(".env file loaded successfully")
    }

    // Initialize database connection
    if err := database.InitDB(); err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // Initialize handler
    handler.InitHandler()  // 追加

    // Initialize Firebase
    if err := auth.InitFirebase(); err != nil {
        log.Fatalf("Failed to initialize Firebase: %v", err)
    }

    // Setup Gin
    r := gin.Default()

    // Add middleware
    r.Use(middleware.CORS())

    setupRoutes(r)

    // Start server
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    
    if err := r.Run(":" + port); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}

func setupRoutes(r *gin.Engine) {
    api := r.Group("/api")
    {
        // Auth routes
        auth := api.Group("/auth")
        {
            auth.POST("/signup", handler.Signup)
            auth.POST("/login", handler.Login)
        }

        // Protected routes
        protected := api.Group("")
        protected.Use(middleware.AuthMiddleware())
        {
            // Posts
            posts := protected.Group("/posts")
            {
                posts.GET("", handler.GetPosts)
                posts.POST("", handler.CreatePost)
                
                // Likes - これらを先に定義
                posts.POST("/:id/likes", handler.CreateLike)
                posts.DELETE("/:id/likes", handler.DeleteLike)
                posts.GET("/:id/likes", handler.GetPostLikes)
                posts.GET("/:id/like-status", handler.CheckLikeStatus)
                
                // 個別の投稿取得 - これを最後に定義
                posts.GET("/:id", handler.GetPost)
            }
        }
    }
}