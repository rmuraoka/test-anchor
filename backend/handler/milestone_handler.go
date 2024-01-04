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
	"time"
)

type MilestoneHandler struct {
	DB *gorm.DB
}

func NewMilestoneHandler(db *gorm.DB) *MilestoneHandler {
	return &MilestoneHandler{DB: db}
}

func (h *MilestoneHandler) GetMilestones(c *gin.Context) {
	projectCode := c.Param("project_code")
	if projectCode == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Project Code is required"})
		return
	}

	var project model.Project
	if result := h.DB.Where("code = ?", projectCode).First(&project); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			handleError(c, http.StatusNotFound, "Project not found", result.Error)
			return
		}
		handleError(c, http.StatusInternalServerError, "Failed to retrieve project", result.Error)
		return
	}

	var Milestones []model.Milestone
	result := h.DB.Where("project_id = ?", project.ID).Find(&Milestones)
	if result.Error != nil {
		log.Fatal("Failed to retrieve records: ", result.Error)
	}

	// Prepare the response data
	milestoneResponses := []util.Milestone{}
	for _, Milestone := range Milestones {
		milestoneResponses = append(milestoneResponses, util.Milestone{
			ID:            Milestone.ID,
			Title:         Milestone.Title,
			Description:   Milestone.Description,
			DueDate:       Milestone.DueDate.Format("2006-01-02"),
			Status:        Milestone.Status,
			TestCaseCount: 0,
		})
	}

	data := util.MilestonesResponseData{
		ProjectId:  project.ID,
		Milestones: milestoneResponses,
	}

	c.JSON(http.StatusOK, data)
}

func (h *MilestoneHandler) GetMilestone(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var Milestone model.Milestone
	if result := h.DB.First(&Milestone, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Milestone not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	c.JSON(http.StatusOK, Milestone)
}

type jsonMilestone struct {
	ProjectID   uint   `json:"project_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date"`
}

func (h *MilestoneHandler) PostMilestone(c *gin.Context) {
	var jsonMilestone jsonMilestone
	if err := c.ShouldBindJSON(&jsonMilestone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", jsonMilestone.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	newMilestone := model.Milestone{
		ProjectID:   jsonMilestone.ProjectID,
		Title:       jsonMilestone.Title,
		Description: jsonMilestone.Description,
		Status:      jsonMilestone.Status,
		DueDate:     parsedDate,
	}

	if result := h.DB.Create(&newMilestone); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, newMilestone)
}

func (h *MilestoneHandler) PutMilestone(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var jsonMilestone jsonMilestone
	if err := c.ShouldBindJSON(&jsonMilestone); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	parsedDate, err := time.Parse("2006-01-02", jsonMilestone.DueDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format"})
		return
	}

	updatedMilestone := model.Milestone{
		ProjectID:   jsonMilestone.ProjectID,
		Title:       jsonMilestone.Title,
		Description: jsonMilestone.Description,
		Status:      jsonMilestone.Status,
		DueDate:     parsedDate,
	}

	if result := h.DB.Model(&model.Milestone{}).Where("id = ?", id).Updates(updatedMilestone); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedMilestone)
}

func (h *MilestoneHandler) DeleteMilestone(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if result := h.DB.Delete(&model.Milestone{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
