package handler

import (
	"backend/middleware"
	"backend/model"
	"backend/util"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"net/http"
)

type AuthHandler struct {
	DB *gorm.DB
}

func NewAuthHandler(db *gorm.DB) *AuthHandler {
	return &AuthHandler{DB: db}
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) PostLogin(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user model.User
	if result := h.DB.Preload("Role").Preload("Role.RolePermissions").Preload("Role.RolePermissions.Permission").Where("email = ?", input.Email).First(&user); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			handleError(c, http.StatusUnauthorized, "User not found", result.Error)
			return
		}
		handleError(c, http.StatusInternalServerError, "Failed to retrieve project", result.Error)
		return
	}
	userHashedPassword := user.Password

	// パスワードのチェック
	if !checkPasswordHash(input.Password, userHashedPassword) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "認証に失敗しました"})
		return
	}

	token, err := middleware.GenerateJWT(input.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "トークンの生成に失敗しました"})
		return
	}
	var permissions []string
	for _, permission := range user.Role.RolePermissions {
		permissions = append(permissions, permission.Permission.Name)
	}
	userJson := util.Auth{
		Token: token,
		User: util.LoginUser{
			ID:          user.ID,
			Name:        user.Name,
			Language:    user.Language,
			Permissions: permissions,
		},
	}

	c.JSON(http.StatusOK, userJson)
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
