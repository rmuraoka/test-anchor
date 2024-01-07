package handler

import (
	"backend/model"
	"backend/util"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
)

type StatusHandler struct {
	DB *gorm.DB
}

func NewStatusHandler(db *gorm.DB) *StatusHandler {
	return &StatusHandler{DB: db}
}

func (h *StatusHandler) GetStatues(c *gin.Context) {
	var statuses []model.Status
	result := h.DB.Find(&statuses)
	if result.Error != nil {
		log.Fatal("Failed to retrieve records: ", result.Error)
	}

	statusResponses := []util.Status{}
	var defaultId uint
	for _, status := range statuses {
		if status.Default {
			defaultId = status.ID
		}
		statusResponses = append(statusResponses, util.Status{
			ID:    status.ID,
			Name:  status.Name,
			Color: status.Color,
		})
	}

	data := util.StatusesResponseData{
		DefaultID: defaultId,
		Statuses:  statusResponses,
	}

	c.JSON(http.StatusOK, data)
}
