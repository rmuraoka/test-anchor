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

type ProjectHandler struct {
	DB *gorm.DB
}

func NewProjectHandler(db *gorm.DB) *ProjectHandler {
	return &ProjectHandler{DB: db}
}

func (h *ProjectHandler) GetProjects(c *gin.Context) {
	var projects []model.Project
	result := h.DB.Find(&projects)
	if result.Error != nil {
		log.Fatal("Failed to retrieve records: ", result.Error)
	}

	// Prepare the response data
	projectResponses := []util.Project{}
	for _, project := range projects {
		projectResponses = append(projectResponses, util.Project{
			ID: project.ID, Title: project.Title, Code: project.Code, Description: project.Description,
		})
	}

	data := util.ProjectsResponseData{
		Projects: projectResponses,
	}

	c.JSON(http.StatusOK, data)
}

func (h *ProjectHandler) GetProject(c *gin.Context) {
	projectCode := c.Param("project_code")
	if projectCode == "" {
		handleError(c, http.StatusBadRequest, "Project Code is required", nil)
		return
	}

	var project model.Project
	if result := h.DB.Preload("Milestones").Preload("TestPlans.CreatedBy").Preload("TestPlans.UpdatedBy").Preload("TestPlans").Where("code = ?", projectCode).First(&project); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Project not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	milestones := make([]util.Milestone, len(project.Milestones))
	for i, m := range project.Milestones {
		milestones[i] = util.Milestone{
			ID:            m.ID,
			Title:         m.Title,
			Description:   m.Description,
			DueDate:       m.DueDate.Format("2006-01-02"),
			Status:        m.Status,
			TestCaseCount: 0, // このフィールド名は実際のmodel.Milestone型に応じて調整してください
		}
	}

	testRuns := make([]util.TestPlan, len(project.TestPlans))
	for i, tp := range project.TestPlans {
		var startedAtStr *string
		var completedAtStr *string
		if tp.StartedAt != nil && !tp.StartedAt.IsZero() {
			formattedStartedAt := tp.StartedAt.Format("2006-01-02 15:04")
			startedAtStr = &formattedStartedAt
		}
		if tp.CompletedAt != nil && !tp.CompletedAt.IsZero() {
			formattedCompletedAt := tp.CompletedAt.Format("2006-01-02 15:04")
			completedAtStr = &formattedCompletedAt
		}
		testRuns[i] = util.TestPlan{
			ID:          tp.ID,
			ProjectID:   tp.ProjectID,
			Title:       tp.Title,
			Status:      tp.Status,
			StartedAt:   startedAtStr,
			CompletedAt: completedAtStr,
			CreatedBy:   util.User{ID: tp.CreatedByID, Name: tp.CreatedBy.Name},
			UpdatedBy:   util.User{ID: tp.UpdatedByID, Name: tp.UpdatedBy.Name},
		}
	}

	response := util.ProjectResponseData{
		ID:          project.ID,
		Code:        project.Code,
		Title:       project.Title,
		Description: project.Description,
		Milestones:  milestones,
		TestPlans:   testRuns,
	}

	c.JSON(http.StatusOK, response)
}

func (h *ProjectHandler) PostProject(c *gin.Context) {
	var newProject model.Project
	if err := c.ShouldBindJSON(&newProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.DB.Create(&newProject); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, newProject)
}

func (h *ProjectHandler) PutProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updatedProject model.Project
	if err := c.ShouldBindJSON(&updatedProject); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.DB.Model(&model.Project{}).Where("id = ?", id).Updates(updatedProject); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedProject)
}

func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if result := h.DB.Delete(&model.Project{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
