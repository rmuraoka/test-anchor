package handler

import (
	"backend/model"
	"backend/util"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type TestPlanHandler struct {
	DB *gorm.DB
}

func NewTestPlanHandler(db *gorm.DB) *TestPlanHandler {
	return &TestPlanHandler{DB: db}
}

func (h *TestPlanHandler) GetTestPlans(c *gin.Context) {
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

	var testPlans []model.TestPlan
	result := h.DB.Preload("CreatedBy").Preload("UpdatedBy").Where("project_id = ?", project.ID).Find(&testPlans)
	if result.Error != nil {
		log.Fatal("Failed to retrieve records: ", result.Error)
	}

	// Prepare the response data
	testPlanResponses := []util.TestPlan{}
	for _, testPlan := range testPlans {
		var startedAtStr *string
		var completedAtStr *string
		if testPlan.StartedAt != nil && !testPlan.StartedAt.IsZero() {
			formattedStartedAt := testPlan.StartedAt.Format("2006-01-02 15:04")
			startedAtStr = &formattedStartedAt
		}
		if testPlan.CompletedAt != nil && !testPlan.CompletedAt.IsZero() {
			formattedCompletedAt := testPlan.CompletedAt.Format("2006-01-02 15:04")
			completedAtStr = &formattedCompletedAt
		}

		testPlanResponses = append(testPlanResponses, util.TestPlan{
			ID:          testPlan.ID,
			ProjectID:   testPlan.ProjectID,
			Title:       testPlan.Title,
			Status:      testPlan.Status,
			StartedAt:   startedAtStr,
			CompletedAt: completedAtStr,
			CreatedBy: util.User{
				ID:   testPlan.CreatedByID,
				Name: testPlan.CreatedBy.Name,
			},
			UpdatedBy: util.User{
				ID:   testPlan.UpdatedByID,
				Name: testPlan.UpdatedBy.Name,
			},
		})
	}

	data := util.TestPlansResponseData{
		ProjectId: project.ID,
		TestPlans: testPlanResponses,
	}

	c.JSON(http.StatusOK, data)
}

func (h *TestPlanHandler) GetTestPlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var testPlan model.TestPlan
	if result := h.DB.Preload("CreatedBy").Preload("UpdatedBy").Preload("TestRuns").Preload("TestRuns.TestRunCases.Status").First(&testPlan, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Test Plan not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}

	testRuns := []util.TestRun{}
	for _, run := range testPlan.TestRuns {
		// 特定のTestRunに紐づくTestCaseのIDの一覧を取得
		var testCaseIDs []uint
		h.DB.Model(&model.TestRunCase{}).Where("test_run_id = ?", run.ID).Pluck("test_case_id", &testCaseIDs)
		var startedAtStr *string
		var completedAtStr *string

		if run.StartedAt != nil && !run.StartedAt.IsZero() {
			formattedStartedAt := run.StartedAt.Format("2006-01-02 15:04")
			startedAtStr = &formattedStartedAt
		}

		if run.CompletedAt != nil && !run.CompletedAt.IsZero() {
			formattedCompletedAt := run.CompletedAt.Format("2006-01-02 15:04")
			completedAtStr = &formattedCompletedAt
		}
		testRuns = append(testRuns, util.TestRun{
			ID:          run.ID,
			ProjectID:   run.ProjectID,
			Title:       run.Title,
			Count:       len(testCaseIDs),
			TestCaseIDs: testCaseIDs,
			Status:      run.Status,
			StartedAt:   startedAtStr,
			CompletedAt: completedAtStr,
			CreatedBy: util.User{
				ID:   run.CreatedByID,
				Name: run.CreatedBy.Name,
			},
			UpdatedBy: util.User{
				ID:   run.UpdatedByID,
				Name: run.UpdatedBy.Name,
			},
		})
	}

	charts := aggregateStatusCounts(testPlan.TestRuns)
	testPlanResponses := util.TestPlanDetail{
		ID:        testPlan.ID,
		ProjectID: testPlan.ProjectID,
		Title:     testPlan.Title,
		Status:    testPlan.Status,
		TestRuns:  testRuns,
		Charts:    charts,
		CreatedBy: util.User{
			ID:   testPlan.CreatedByID,
			Name: testPlan.CreatedBy.Name,
		},
		UpdatedBy: util.User{
			ID:   testPlan.UpdatedByID,
			Name: testPlan.UpdatedBy.Name,
		},
		Percentage: calcPercentage(testPlan.TestRuns),
	}

	c.JSON(http.StatusOK, testPlanResponses)
}

type statusInfo struct {
	ID    uint
	Name  string
	Color string
	Count int
}

func aggregateStatusCounts(testRuns []model.TestRun) []util.Chart {
	statusCounts := make(map[string]int)
	statusColors := make(map[string]string)
	statusIDs := make(map[string]uint)

	for _, run := range testRuns {
		for _, testRunCase := range run.TestRunCases {
			status := testRunCase.Status.Name
			color := testRunCase.Status.Color
			id := testRunCase.Status.ID // ステータスからIDを取得
			statusCounts[status]++
			statusColors[status] = color
			statusIDs[status] = id
		}
	}

	var infos []statusInfo
	for status, count := range statusCounts {
		info := statusInfo{
			ID:    statusIDs[status],
			Name:  status,
			Color: statusColors[status],
			Count: count,
		}
		infos = append(infos, info)
	}

	sort.Slice(infos, func(i, j int) bool {
		return infos[i].ID > infos[j].ID
	})

	charts := make([]util.Chart, len(infos))
	for i, info := range infos {
		charts[i] = util.Chart{
			Name:  info.Name,
			Color: info.Color,
			Count: info.Count,
		}
	}

	return charts
}

func (h *TestPlanHandler) PostTestPlan(c *gin.Context) {
	var newTestPlan model.TestPlan
	if err := c.ShouldBindJSON(&newTestPlan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.DB.Create(&newTestPlan); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTestPlan)
}

func (h *TestPlanHandler) PutTestPlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updatedTestPlan model.TestPlan
	if err := c.ShouldBindJSON(&updatedTestPlan); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.DB.Model(&model.TestPlan{}).Where("id = ?", id).Updates(updatedTestPlan); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTestPlan)
}

func (h *TestPlanHandler) DeletePlan(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	if result := h.DB.Delete(&model.TestPlan{}, id); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
