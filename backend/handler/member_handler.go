package handler

import (
	"backend/model"
	"backend/util"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"strconv"
)

type MemberHandler struct {
	DB          *gorm.DB
	EmailSender util.EmailSender
}

func NewMemberHandler(db *gorm.DB, emailSender util.EmailSender) *MemberHandler {
	return &MemberHandler{DB: db, EmailSender: emailSender}
}

func (h *MemberHandler) GetMembers(c *gin.Context) {
	var users []model.User
	result := h.DB.Find(&users)
	if result.Error != nil {
		log.Fatal("Failed to retrieve records: ", result.Error)
	}

	// Prepare the response data
	memberResponses := []util.Member{}
	for _, user := range users {
		memberResponses = append(memberResponses, util.Member{
			ID:     user.ID,
			Email:  user.Email,
			Name:   user.Name,
			Status: user.Status,
		})
	}

	data := util.MembersResponseData{Members: memberResponses}

	c.JSON(http.StatusOK, data)
}

func (h *MemberHandler) GetMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var user model.User
	result := h.DB.First(&user, id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Member is not found"})
		return
	}

	memberResponses := util.Member{
		ID:     user.ID,
		Email:  user.Email,
		Name:   user.Name,
		Status: user.Status,
	}

	c.JSON(http.StatusOK, memberResponses)
}

func (h *MemberHandler) PostMember(c *gin.Context) {
	var newUser model.User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tempPassword := util.GenerateTempPassword(10)
	hashedPassword, err := util.HashPassword(tempPassword)
	if err != nil {
		log.Fatalf("パスワードのハッシュ化に失敗しました: %v", err)
	}
	newUser.Password = hashedPassword

	if result := h.DB.Create(&newUser); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	memberResponses := util.Member{
		ID:     newUser.ID,
		Email:  newUser.Email,
		Name:   newUser.Name,
		Status: newUser.Status,
	}

	subject := "Your Account"
	body := "Welcome " + newUser.Name + " Your Password is " + tempPassword
	to := []string{newUser.Email}

	// メール送信の実行
	err = h.EmailSender.SendMail(to, subject, body)
	if err != nil {
		log.Printf("ユーザーの招待に失敗しました: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email"})
		return
	}

	c.JSON(http.StatusCreated, memberResponses)
}

func (h *MemberHandler) PutMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	currentPassword := c.PostForm("current_password")
	if currentPassword != "" {
		var user model.User
		if result := h.DB.First(&user, id); result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				handleError(c, http.StatusUnauthorized, "User not found", result.Error)
				return
			}
			handleError(c, http.StatusInternalServerError, "Failed to retrieve project", result.Error)
			return
		}
		userHashedPassword := user.Password

		// パスワードのチェック
		if !checkPasswordHash(currentPassword, userHashedPassword) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "認証に失敗しました"})
			return
		}
	}

	var updatedMember model.User
	if err := c.ShouldBindJSON(&updatedMember); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if updatedMember.Password != "" {
		tempPassword := updatedMember.Password
		hashedPassword, err := util.HashPassword(tempPassword)
		if err != nil {
			log.Fatalf("パスワードのハッシュ化に失敗しました: %v", err)
		}
		updatedMember.Password = hashedPassword
	}

	if result := h.DB.Model(&model.User{}).Where("id = ?", id).Updates(updatedMember); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	memberResponses := util.Member{
		ID:     updatedMember.ID,
		Email:  updatedMember.Email,
		Name:   updatedMember.Name,
		Status: updatedMember.Status,
	}

	c.JSON(http.StatusOK, memberResponses)
}

func (h *MemberHandler) DeleteMember(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if result := h.DB.Delete(&model.User{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
