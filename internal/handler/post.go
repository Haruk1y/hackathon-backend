// internal/handler/post.go

package handler

import (
	"net/http"
	"time"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Haruk1y/hackathon-backend/internal/model"
	"github.com/Haruk1y/hackathon-backend/internal/ai"
)

type CreatePostRequest struct {
	Content string `json:"content" binding:"required"`
}

// 投稿の作成
func CreatePost(c *gin.Context) {
    userID := c.GetString("uid")
    if userID == "" {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
        return
    }

    var req CreatePostRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    var user model.User
    if err := db.Where("firebase_uid = ?", userID).First(&user).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
        return
    }

    // コンテンツが50文字以上の場合のみ要約を生成
    var summary string
    if len(req.Content) > 50 {
        var err error
        summary, err = ai.SummarizeText(req.Content)
        if err != nil {
            // 要約の生成に失敗しても投稿自体は継続
            log.Printf("Failed to generate summary: %v", err)
        }
    }

    post := model.Post{
        ID:        uuid.New().String(),
        UserID:    user.ID,
        Content:   req.Content,
        Summary:   summary,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    if err := db.Create(&post).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
        return
    }

    if err := db.Preload("User").First(&post, "id = ?", post.ID).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch post"})
        return
    }

    c.JSON(http.StatusCreated, post)
}

// 投稿一覧の取得
func GetPosts(c *gin.Context) {
	var posts []model.Post

	// 最新の投稿から順に取得し、ユーザー情報も一緒に取得
	if err := db.Preload("User").Order("created_at desc").Limit(20).Find(&posts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch posts"})
		return
	}

	c.JSON(http.StatusOK, posts)
}

// 投稿の詳細取得
func GetPost(c *gin.Context) {
	id := c.Param("id")
	var post model.Post

	if err := db.Preload("User").First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	c.JSON(http.StatusOK, post)
}