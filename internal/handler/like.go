// internal/handler/like.go

package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Haruk1y/hackathon-backend/internal/model"
	"gorm.io/gorm"
)

// いいねの追加
func CreateLike(c *gin.Context) {
	userID := c.GetString("uid")
	id := c.Param("id")  // postID から id に変更

	// ユーザー情報の取得
	var user model.User
	if err := db.Where("firebase_uid = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 投稿の存在確認
	var post model.Post
	if err := db.First(&post, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}

	// いいねの重複チェック
	var existingLike model.Like
	err := db.Where("user_id = ? AND post_id = ?", user.ID, id).First(&existingLike).Error
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Like already exists"})
		return
	}

	// トランザクション開始
	err = db.Transaction(func(tx *gorm.DB) error {
		// いいねの作成
		like := model.Like{
			ID:        uuid.New().String(),
			UserID:    user.ID,
			PostID:    id,  // postID から id に変更
			CreatedAt: time.Now(),
		}
		if err := tx.Create(&like).Error; err != nil {
			return err
		}

		// いいね数の更新
		if err := tx.Model(&post).Update("like_count", post.LikeCount+1).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create like"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Like created successfully"})
}

// いいねの削除
func DeleteLike(c *gin.Context) {
	userID := c.GetString("uid")
	id := c.Param("id")  // postID から id に変更

	// ユーザー情報の取得
	var user model.User
	if err := db.Where("firebase_uid = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// トランザクション開始
	err := db.Transaction(func(tx *gorm.DB) error {
		// いいねの存在確認と削除
		var like model.Like
		if err := tx.Where("user_id = ? AND post_id = ?", user.ID, id).First(&like).Error; err != nil {
			return err
		}

		if err := tx.Delete(&like).Error; err != nil {
			return err
		}

		// いいね数の更新
		var post model.Post
		if err := tx.First(&post, "id = ?", id).Error; err != nil {
			return err
		}

		if err := tx.Model(&post).Update("like_count", post.LikeCount-1).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Like deleted successfully"})
}

// 投稿のいいね一覧取得
func GetPostLikes(c *gin.Context) {
    id := c.Param("id")

    var likes []model.Like
    if err := db.Preload("User").Preload("Post").Preload("Post.User").Where("post_id = ?", id).Find(&likes).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch likes"})
        return
    }

    c.JSON(http.StatusOK, likes)
}

// ユーザーがいいねしているかどうかの確認
func CheckLikeStatus(c *gin.Context) {
	userID := c.GetString("uid")
	id := c.Param("id")  // postID から id に変更

	var user model.User
	if err := db.Where("firebase_uid = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var like model.Like
	isLiked := db.Where("user_id = ? AND post_id = ?", user.ID, id).First(&like).Error == nil

	c.JSON(http.StatusOK, gin.H{"is_liked": isLiked})
}