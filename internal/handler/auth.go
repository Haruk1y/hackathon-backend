package handler

import (
	"net/http"
	"strings"

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
	idToken := strings.Replace(authHeader, "Bearer ", "", 1)

	// Firebaseトークンの検証
	token, err := auth.VerifyIDToken(idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// ユーザーの作成
	user := model.User{
		ID:          uuid.New().String(),
		FirebaseUID: token.UID,
		Username:    req.Username,
		DisplayName: req.DisplayName,
	}

	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func Login(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
		return
	}

	idToken := strings.Replace(authHeader, "Bearer ", "", 1)

	// トークンの検証
	token, err := auth.VerifyIDToken(idToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// ユーザー情報の取得
	var user model.User
	if err := db.Where("firebase_uid = ?", token.UID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}