// internal/handler/auth.go

package handler

import (
    "log"
    "net/http"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/google/uuid"
    "github.com/Haruk1y/hackathon-backend/internal/auth"
    "github.com/Haruk1y/hackathon-backend/internal/model"
)

type SignupRequest struct {
    Username    string `json:"username" binding:"required"`
    DisplayName string `json:"displayName" binding:"required"`
}

func Signup(c *gin.Context) {
    var req SignupRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // Authorizationヘッダーからトークンを取得
    authHeader := c.GetHeader("Authorization")
    if authHeader == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
        return
    }

    // "Bearer "を除去してトークンを取得
    idToken := strings.TrimPrefix(authHeader, "Bearer ")
    if idToken == authHeader {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
        return
    }

    // Firebaseトークンの検証
    token, err := auth.VerifyIDToken(idToken)
    if err != nil {
        log.Printf("Token verification failed: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    // ユーザーの作成
    user := model.User{
        ID:          uuid.New().String(),
        FirebaseUID: token.UID,
        Username:    req.Username,
        DisplayName: req.DisplayName,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
        IsActive:    true,
    }

    if err := db.Create(&user).Error; err != nil {
        log.Printf("Failed to create user: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
        return
    }

    c.JSON(http.StatusCreated, user)
}

func Login(c *gin.Context) {
    authHeader := c.GetHeader("Authorization")
    log.Printf("Auth header received: %s", authHeader) // デバッグ用

    if authHeader == "" || authHeader == "*" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Valid Authorization header is required"})
        return
    }

    idToken := strings.TrimPrefix(authHeader, "Bearer ")
    if idToken == authHeader {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
        return
    }

    // Firebaseトークンの検証
    token, err := auth.VerifyIDToken(idToken)
    if err != nil {
        log.Printf("Token verification failed: %v", err)
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
        return
    }

    // ユーザー情報の取得
    var user model.User
    if err := db.Where("firebase_uid = ?", token.UID).First(&user).Error; err != nil {
        log.Printf("User not found: %v", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    c.JSON(http.StatusOK, user)
}