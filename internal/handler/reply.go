// internal/handler/reply.go

package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Haruk1y/hackathon-backend/internal/model"
	"gorm.io/gorm"
)

type CreateReplyRequest struct {
	Content string `json:"content" binding:"required"`
}

// リプライの作成
func CreateReply(c *gin.Context) {
	userID := c.GetString("uid")
	parentID := c.Param("id") // 親投稿のID

	// リクエストのバインド
	var req CreateReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザー情報の取得
	var user model.User
	if err := db.Where("firebase_uid = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 親投稿の存在確認
	var parentPost model.Post
	if err := db.First(&parentPost, "id = ?", parentID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Parent post not found"})
		return
	}

	// トランザクション開始
	err := db.Transaction(func(tx *gorm.DB) error {
		// リプライの作成
		reply := model.Post{
			ID:           uuid.New().String(),
			UserID:       user.ID,
			Content:      req.Content,
			IsReply:      true,
			ParentPostID: &parentID,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := tx.Create(&reply).Error; err != nil {
			return err
		}

		// 親投稿のリプライ数を更新
		if err := tx.Model(&parentPost).Update("reply_count", parentPost.ReplyCount+1).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reply"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Reply created successfully"})
}

// 投稿に対するリプライ一覧の取得
func GetReplies(c *gin.Context) {
	postID := c.Param("id")

	var replies []model.Post
	if err := db.Preload("User").
		Where("is_reply = ? AND parent_post_id = ?", true, postID).
		Order("created_at desc").
		Find(&replies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch replies"})
		return
	}

	c.JSON(http.StatusOK, replies)
}

// 投稿の詳細取得（リプライ情報を含む）
func GetPostWithReplies(c *gin.Context) {
	id := c.Param("id")

	var post model.Post
	if err := db.Preload("User").First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	var replies []model.Post
	if err := db.Preload("User").
		Where("is_reply = ? AND parent_post_id = ?", true, id).
		Order("created_at desc").
		Find(&replies).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch replies"})
		return
	}

	response := gin.H{
		"post":    post,
		"replies": replies,
	}

	c.JSON(http.StatusOK, response)
}