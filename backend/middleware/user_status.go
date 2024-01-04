// middleware/user_status.go

package middleware

import (
	"backend/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

func CheckUserStatus(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// JWTからユーザーIDを取得
		email := c.GetString("email")

		var user model.User
		if err := db.Where("email = ?", email).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "内部サーバーエラー"})
			return
		}

		// ユーザーの状態を確認
		if user.Status == "InActive" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "アカウントが停止されています"})
			return
		}

		c.Next()
	}
}
