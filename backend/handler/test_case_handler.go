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

type TestCaseHandler struct {
	DB *gorm.DB
}

func NewTestCaseHandler(db *gorm.DB) *TestCaseHandler {
	return &TestCaseHandler{DB: db}
}

func handleError(c *gin.Context, status int, msg string, err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
	c.JSON(status, gin.H{"error": msg})
}

func createTestCaseResponse(testCase model.TestCase) util.TestCase {
	var milestone *util.TestCaseMilestone
	if testCase.MilestoneID != nil {
		milestone = &util.TestCaseMilestone{ID: *testCase.MilestoneID, Title: testCase.Milestone.Title}
	}

	return util.TestCase{
		ID:        testCase.ID,
		Title:     testCase.Title,
		Content:   testCase.Content,
		Milestone: milestone,
		CreatedBy: util.User{ID: testCase.CreatedByID, Name: testCase.CreatedBy.Name},
		UpdatedBy: util.User{ID: testCase.UpdatedByID, Name: testCase.UpdatedBy.Name},
	}
}

func (h *TestCaseHandler) GetTestCases(c *gin.Context) {
	projectCode := c.Param("project_code")
	if projectCode == "" {
		handleError(c, http.StatusBadRequest, "Project ID is required", nil)
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

	testSuites := []model.TestSuite{}
	if result := h.DB.Where("project_id = ?", project.ID).
		Order("test_suites.order_index ASC, test_suites.id ASC"). // TestSuite の並び順をここで指定
		Preload("TestCases", func(db *gorm.DB) *gorm.DB {
			return db.Order("test_cases.order_index ASC, test_cases.id ASC")
		}).
		Preload("TestCases.Milestone").
		Preload("TestCases.CreatedBy").
		Preload("TestCases.UpdatedBy").
		Preload("TestSuites.TestCases", func(db *gorm.DB) *gorm.DB {
			return db.Order("test_cases.order_index ASC, test_cases.id ASC")
		}).
		Preload("TestSuites.TestCases.Milestone").
		Find(&testSuites); result.Error != nil {
		handleError(c, http.StatusInternalServerError, "Failed to retrieve records", result.Error)
		return
	}

	testSuiteMap := make(map[uint][]model.TestSuite)
	for _, testSuite := range testSuites {
		var parentID uint
		if testSuite.ParentID != nil {
			parentID = *testSuite.ParentID
		}
		testSuiteMap[parentID] = append(testSuiteMap[parentID], testSuite)
	}

	topLevelTestSuites := testSuiteMap[0]
	jsonTestSuites := convertToJSONTestSuites(topLevelTestSuites, testSuiteMap)
	jsonOnlyTestSuites := convertToJSONOnlyTestSuites(topLevelTestSuites, testSuiteMap)

	responseData := util.TestCasesResponseData{
		ProjectID:      project.ID,
		TestSuites:     jsonTestSuites,
		OnlyTestSuites: jsonOnlyTestSuites,
	}

	c.JSON(http.StatusOK, responseData)
}

func convertToJSONTestSuites(testSuites []model.TestSuite, testSuiteMap map[uint][]model.TestSuite) []util.JSONTestSuite {
	jsonTestSuites := []util.JSONTestSuite{}
	for _, testSuite := range testSuites {
		jsonTestCases := []util.TestCase{}
		for _, testCase := range testSuite.TestCases {
			jsonTestCases = append(jsonTestCases, createTestCaseResponse(testCase))
		}

		subTestSuites := convertToJSONTestSuites(testSuiteMap[testSuite.ID], testSuiteMap)
		jsonTestSuites = append(jsonTestSuites, util.JSONTestSuite{
			ID:         testSuite.ID,
			Name:       testSuite.Name,
			TestSuites: subTestSuites,
			TestCases:  jsonTestCases,
		})
	}
	return jsonTestSuites
}

func convertToJSONOnlyTestSuites(testSuites []model.TestSuite, testSuiteMap map[uint][]model.TestSuite) []util.JSONOnlyTestSuite {
	jsonTestSuites := []util.JSONOnlyTestSuite{}

	for _, testSuite := range testSuites {
		subTestSuites := []util.JSONOnlyTestSuite{}
		if _, ok := testSuiteMap[testSuite.ID]; ok {
			subTestSuites = convertToJSONOnlyTestSuites(testSuiteMap[testSuite.ID], testSuiteMap)
		}
		jsonTestSuites = append(jsonTestSuites, util.JSONOnlyTestSuite{
			ID:       testSuite.ID,
			Title:    testSuite.Name,
			Children: subTestSuites,
		})
	}

	return jsonTestSuites
}

func (h *TestCaseHandler) GetTestCase(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	var testCase model.TestCase
	if result := h.DB.Preload("CreatedBy").Preload("UpdatedBy").Preload("Milestone").First(&testCase, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			handleError(c, http.StatusNotFound, "Test case not found", result.Error)
		} else {
			handleError(c, http.StatusInternalServerError, "Database error", result.Error)
		}
		return
	}

	c.JSON(http.StatusOK, createTestCaseResponse(testCase))
}

func (h *TestCaseHandler) PostTestCase(c *gin.Context) {
	var newTestCase model.TestCase
	if err := c.ShouldBindJSON(&newTestCase); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	if result := h.DB.Create(&newTestCase); result.Error != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create test case", result.Error)
		return
	}

	c.JSON(http.StatusCreated, createTestCaseResponse(newTestCase))
}

func (h *TestCaseHandler) PutTestCase(c *gin.Context) {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		// エラー処理
	}
	id := uint(idInt)

	var updatedTestCase model.TestCase
	if err := c.ShouldBindJSON(&updatedTestCase); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	var existingTestCase model.TestCase
	if result := h.DB.First(&existingTestCase, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			handleError(c, http.StatusNotFound, "Test case not found", result.Error)
		} else {
			handleError(c, http.StatusInternalServerError, "Database error", result.Error)
		}
		return
	}

	if result := h.DB.Model(&existingTestCase).Updates(updatedTestCase); result.Error != nil {
		handleError(c, http.StatusInternalServerError, "Failed to update test case", result.Error)
		return
	}

	if updatedTestCase.MilestoneID == nil {
		if result := h.DB.Model(&existingTestCase).Select("MilestoneID").Updates(updatedTestCase); result.Error != nil {
			handleError(c, http.StatusInternalServerError, "Failed to update test case", result.Error)
			return
		}
	}

	c.JSON(http.StatusOK, createTestCaseResponse(existingTestCase))
}

func (h *TestCaseHandler) DeleteTestCase(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	if result := h.DB.Delete(&model.TestCase{}, id); result.Error != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete test case", result.Error)
		return
	}

	if result := h.DB.Where("test_case_id = ?", id).Delete(&model.TestRunCase{}); result.Error != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete test run case", result.Error)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *TestCaseHandler) PostTestSuite(c *gin.Context) {
	var newTestSuite model.TestSuite
	if err := c.ShouldBindJSON(&newTestSuite); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	if result := h.DB.Create(&newTestSuite); result.Error != nil {
		handleError(c, http.StatusInternalServerError, "Failed to create test case", result.Error)
		return
	}

	c.JSON(http.StatusCreated, newTestSuite)
}

func (h *TestCaseHandler) PutTestSuite(c *gin.Context) {
	idStr := c.Param("id")
	idInt, err := strconv.Atoi(idStr)
	if err != nil {
		// エラー処理
	}
	id := uint(idInt)

	var updatedTestSuite model.TestSuite
	if err := c.ShouldBindJSON(&updatedTestSuite); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	var existingTestSuite model.TestSuite
	if result := h.DB.First(&existingTestSuite, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			handleError(c, http.StatusNotFound, "TestSuite not found", result.Error)
		} else {
			handleError(c, http.StatusInternalServerError, "Database error", result.Error)
		}
		return
	}

	if result := h.DB.Model(&existingTestSuite).Updates(updatedTestSuite); result.Error != nil {
		handleError(c, http.StatusInternalServerError, "Failed to update test case", result.Error)
		return
	}

	c.JSON(http.StatusOK, existingTestSuite)
}

func (h *TestCaseHandler) DeleteTestSuite(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	if result := h.DB.Delete(&model.TestSuite{}, id); result.Error != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete test suite", result.Error)
		return
	}

	if result := h.DB.Where("test_suite_id = ?", id).Delete(&model.TestCase{}); result.Error != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete test case", result.Error)
		return
	}

	if result := h.DB.Where("test_case_id IN (SELECT id FROM test_cases WHERE test_suite_id = ?)", id).Delete(&model.TestRunCase{}); result.Error != nil {
		handleError(c, http.StatusInternalServerError, "Failed to delete test run case", result.Error)
		return
	}

	c.Status(http.StatusNoContent)
}

type TestCaseRequestRow struct {
	TestCaseID uint `json:"test_case_id"`
	OrderIndex int  `json:"index"`
}

type TestCaseRequest struct {
	TestSuiteID         uint                 `json:"test_suite_id"`
	TestCaseRequestRows []TestCaseRequestRow `json:"test_cases"`
}

func (h *TestCaseHandler) PutTestCaseBulk(c *gin.Context) {
	var req TestCaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, row := range req.TestCaseRequestRows {
		update := map[string]interface{}{
			"test_suite_id": req.TestSuiteID,
			"order_index":   row.OrderIndex,
		}
		if err := h.DB.Model(&model.TestCase{}).
			Where("id = ?", row.TestCaseID).
			Updates(update).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

type TestSuiteRequestRow struct {
	TestSuiteID uint `json:"test_suite_id"`
	OrderIndex  int  `json:"index"`
}

type TestSuiteRequest struct {
	ParentID            *uint                 `json:"parent_id"`
	TestCaseRequestRows []TestSuiteRequestRow `json:"test_suites"`
}

func (h *TestCaseHandler) PutTestSuiteBulk(c *gin.Context) {
	var req TestSuiteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	for _, row := range req.TestCaseRequestRows {
		update := map[string]interface{}{
			"parent_id":   req.ParentID,
			"order_index": row.OrderIndex,
		}
		if err := h.DB.Model(&model.TestSuite{}).
			Where("id = ?", row.TestSuiteID).
			Updates(update).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Update failed"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
