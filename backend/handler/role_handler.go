package handler

import (
	"backend/model"
	"backend/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type RoleHandler struct {
	DB *gorm.DB
}

func NewRoleHandler(db *gorm.DB) *RoleHandler {
	return &RoleHandler{DB: db}
}

func (h *RoleHandler) GetRoles(c *gin.Context) {
	var statuses []model.Role
	result := h.DB.Find(&statuses)
	if result.Error != nil {
		log.Fatal("Failed to retrieve records: ", result.Error)
	}

	roleResponses := []util.Role{}
	for _, role := range statuses {
		roleResponses = append(roleResponses, util.Role{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
		})
	}

	data := util.RolesResponseData{
		Roles: roleResponses,
	}

	c.JSON(http.StatusOK, data)
}
