package handler

import (
	"backend/model"
	"backend/util"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
	"net/http"
	"sort"
	"strconv"
)

type TestRunHandler struct {
	DB *gorm.DB
}

func NewTestRunHandler(db *gorm.DB) *TestRunHandler {
	return &TestRunHandler{DB: db}
}

func (h *TestRunHandler) GetTestRuns(c *gin.Context) {
	projectCode := c.Param("project_code")
	if projectCode == "" {
		handleError(c, http.StatusBadRequest, "Project Code is required", nil)
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

	testRuns := []model.TestRun{}
	if result := h.DB.Preload("CreatedBy").Preload("UpdatedBy").Where("project_id = ?", project.ID).Find(&testRuns); result.Error != nil {
		handleError(c, http.StatusInternalServerError, "Failed to retrieve records", result.Error)
		return
	}

	var responseTestRuns []util.TestRun
	for _, testRun := range testRuns {
		var startedAtStr *string
		var completedAtStr *string

		if testRun.StartedAt != nil && !testRun.StartedAt.IsZero() {
			formattedStartedAt := testRun.StartedAt.Format("2006-01-02 15:04")
			startedAtStr = &formattedStartedAt
		}

		if testRun.CompletedAt != nil && !testRun.CompletedAt.IsZero() {
			formattedCompletedAt := testRun.CompletedAt.Format("2006-01-02 15:04")
			completedAtStr = &formattedCompletedAt
		}
		responseTestRun := util.TestRun{
			ID:          testRun.ID,
			ProjectID:   testRun.ProjectID,
			Count:       1,
			Status:      testRun.Status,
			StartedAt:   startedAtStr,
			CompletedAt: completedAtStr,
			Title:       testRun.Title,
			CreatedBy:   convertUserToJSON(testRun.CreatedBy),
			UpdatedBy:   convertUserToJSON(testRun.UpdatedBy),
		}
		responseTestRuns = append(responseTestRuns, responseTestRun)
	}

	responseData := util.TestRunsResponseData{
		TestRuns: responseTestRuns,
	}

	c.JSON(http.StatusOK, responseData)
}

func (h *TestRunHandler) GetTestRunCase(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Invalid ID format", err)
		return
	}

	var testRunCase model.TestRunCase
	if result := h.DB.Preload("Status").
		Preload("AssignedTo").
		Preload("Comments.CreatedBy").
		Preload("Comments.UpdatedBy").
		Preload("Comments").
		Preload("TestCase.CreatedBy").
		Preload("TestCase.UpdatedBy").
		First(&testRunCase, id); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			handleError(c, http.StatusNotFound, "Test run not found", result.Error)
		} else {
			handleError(c, http.StatusInternalServerError, "Database error", result.Error)
		}
		return
	}

	commentsJSON := []util.Comment{}
	for _, comment := range testRunCase.Comments {
		commentsJSON = append(commentsJSON, convertCommentToJSON(comment))
	}

	var assignedToJSON *util.User
	if testRunCase.AssignedTo != nil {
		userJSON := convertUserToJSON(*testRunCase.AssignedTo)
		assignedToJSON = &userJSON // userJSON のアドレスを assignedToJSON に代入
	}
	response := util.TestRunCase{
		ID:         testRunCase.ID,
		TestCaseId: testRunCase.TestCaseID,
		Title:      testRunCase.TestCase.Title,
		Content:    testRunCase.TestCase.Content,
		Status:     util.Status{ID: testRunCase.StatusID, Name: testRunCase.Status.Name},
		AssignedTo: assignedToJSON,
		Comments:   commentsJSON,
	}

	c.JSON(http.StatusOK, response)
}

func convertCommentToJSON(comment model.Comment) util.Comment {
	var status *util.Status
	if comment.Status != nil {
		status = &util.Status{ID: comment.Status.ID, Name: comment.Status.Name, Color: comment.Status.Color}
	} else {
		status = nil
	}
	return util.Comment{
		ID:        comment.ID,
		Status:    status,
		Content:   comment.Content,
		CreatedBy: convertUserToJSON(comment.CreatedBy),
		UpdatedBy: convertUserToJSON(comment.UpdatedBy),
	}
}

func (h *TestRunHandler) GetTestRunCases(c *gin.Context) {
	testRunIDStr := c.Param("id")
	if testRunIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	// 文字列をintに変換
	testRunIDInt, err := strconv.Atoi(testRunIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id format"})
		return
	}

	// intをuintにキャスト
	testRunID := uint(testRunIDInt)

	var testRun = model.TestRun{}
	preResult := h.DB.First(&testRun, testRunID)
	if preResult.Error != nil {
		if errors.Is(preResult.Error, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Test Run not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		}
		return
	}
	if testRun.Status == "Completed" {
		// 完了している場合は内部のデータを参照する
		var data util.TestRunCasesResponseData
		err := json.Unmarshal([]byte(testRun.FinalizedTestCases), &data)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		// 念の為置き換える
		data.TestPlanId = testRun.ProjectID
		data.TestRunID = testRunID
		data.TestPlanId = testRun.TestPlanID
		data.Status = testRun.Status

		c.JSON(http.StatusOK, data)
		return
	}

	var testRunCases []model.TestRunCase
	result := h.DB.
		Preload("TestCase.TestSuite", func(db *gorm.DB) *gorm.DB {
			return db.Order("test_suites.order_index ASC, test_suites.id ASC")
		}).
		Preload("Status").
		Preload("AssignedTo").
		Preload("Comments.CreatedBy").
		Preload("Comments.UpdatedBy").
		Preload("Comments").
		Where("test_run_id = ?", testRunID).
		Find(&testRunCases)
	if result.Error != nil {
		log.Fatal("Failed to retrieve records: ", result.Error)
	}

	testSuiteMap := make(map[uint][]model.TestRunCase)
	testSuiteHierarchyMap := make(map[uint][]model.TestSuite)
	addedTestSuites := make(map[uint]bool)
	for _, trc := range testRunCases {
		var testSuiteID uint
		if trc.TestCase.TestSuiteID != nil {
			testSuiteID = *trc.TestCase.TestSuiteID
		}

		testSuite := trc.TestCase.TestSuite
		parentID := uint(0)
		if testSuite.ParentID != nil {
			parentID = *testSuite.ParentID
		}

		testSuiteMap[testSuiteID] = append(testSuiteMap[testSuiteID], trc)

		if _, exists := addedTestSuites[testSuite.ID]; !exists {
			testSuiteHierarchyMap[parentID] = append(testSuiteHierarchyMap[parentID], testSuite)
			addedTestSuites[testSuite.ID] = true
		}
	}

	topLevelTestSuites := []model.TestSuite{}
	for testSuiteID := range testSuiteMap {
		if testSuiteID != 0 {
			var testSuite model.TestSuite
			h.DB.First(&testSuite, testSuiteID)
			if testSuite.ParentID == nil {
				topLevelTestSuites = append(topLevelTestSuites, testSuite)
			}
		}
	}

	sort.Slice(topLevelTestSuites, func(i, j int) bool {
		if topLevelTestSuites[i].OrderIndex == topLevelTestSuites[j].OrderIndex {
			return topLevelTestSuites[i].ID < topLevelTestSuites[j].ID
		}
		return topLevelTestSuites[i].OrderIndex < topLevelTestSuites[j].OrderIndex
	})

	var allTestSuites []model.TestSuite
	resultAllSuites := h.DB.Where("project_id = ?", testRun.ProjectID).Find(&allTestSuites)
	if resultAllSuites.Error != nil {
		log.Fatal("Failed to retrieve TestSuites: ", result.Error)
	}
	childTestSuitesMap := make(map[uint][]model.TestSuite)
	for _, testSuite := range allTestSuites {
		if testSuite.ParentID != nil {
			parentID := *testSuite.ParentID
			childTestSuitesMap[parentID] = append(childTestSuitesMap[parentID], testSuite)
		}
	}

	jsonTestSuites := convertToTestRunCaseTestSuites(topLevelTestSuites, testSuiteMap, childTestSuitesMap)
	jsonOnlyTestSuites := convertToJSONOnlyTestSuites(topLevelTestSuites, testSuiteHierarchyMap)

	responseData := util.TestRunCasesResponseData{
		ProjectID:      testRun.ProjectID,
		TestRunID:      testRunID,
		TestPlanId:     testRun.TestPlanID,
		Status:         testRun.Status,
		TestSuites:     jsonTestSuites,
		OnlyTestSuites: jsonOnlyTestSuites,
	}

	c.JSON(http.StatusOK, responseData)
}

func convertToTestRunCaseTestSuites(testSuites []model.TestSuite, testRunCaseMap map[uint][]model.TestRunCase, childTestSuitesMap map[uint][]model.TestSuite) []util.TestRunCasesTestSuite {
	sort.Slice(testSuites, func(i, j int) bool {
		return testSuites[i].OrderIndex < testSuites[j].OrderIndex
	})

	jsonTestSuites := []util.TestRunCasesTestSuite{}
	for _, testSuite := range testSuites {
		jsonTestRunCases := []util.TestRunCase{}
		for _, trc := range testRunCaseMap[testSuite.ID] {
			jsonTestRunCases = append(jsonTestRunCases, convertToJSONTestRunCase(trc))
		}

		childTestSuites := childTestSuitesMap[testSuite.ID]

		jsonTestSuites = append(jsonTestSuites, util.TestRunCasesTestSuite{
			Name:       testSuite.Name,
			TestSuites: convertToTestRunCaseTestSuites(childTestSuites, testRunCaseMap, childTestSuitesMap),
			TestCases:  jsonTestRunCases,
		})
	}
	return jsonTestSuites
}

func (h *TestRunHandler) PostTestRun(c *gin.Context) {
	var newTestRun model.TestRun
	if err := c.ShouldBindJSON(&newTestRun); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.DB.Create(&newTestRun); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTestRun)
}

func (h *TestRunHandler) PostTestRunCase(c *gin.Context) {
	var newTestRunCase model.TestRunCase
	if err := c.ShouldBindJSON(&newTestRunCase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.DB.Create(&newTestRunCase); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, newTestRunCase)
}

type TestRunCasesRequest struct {
	TestCaseIDs []uint `json:"test_case_ids"`
	TestRunID   uint   `json:"test_run_id"`
}

func (h *TestRunHandler) PostTestRunCaseBulk(c *gin.Context) {
	var req TestRunCasesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 既存のTestRunCaseを取得
	var existingTestRunCases []model.TestRunCase
	h.DB.Where("test_run_id = ?", req.TestRunID).Find(&existingTestRunCases)

	// 既存のTestCaseIDのマップを作成
	existingTestCaseIDs := make(map[uint]bool)
	for _, trc := range existingTestRunCases {
		existingTestCaseIDs[trc.TestCaseID] = true
	}

	// 新しいTestRunCaseを作成
	status := model.Status{}
	h.DB.Where("statuses.default = ?", 1).First(&status)

	newTestRunCases := []model.TestRunCase{}
	for _, testCaseID := range req.TestCaseIDs {
		if _, exists := existingTestCaseIDs[testCaseID]; !exists {
			newTestRunCase := model.TestRunCase{
				TestRunID:  req.TestRunID,
				TestCaseID: testCaseID,
				StatusID:   status.ID,
			}
			newTestRunCases = append(newTestRunCases, newTestRunCase)
		}
	}

	// 新しいTestRunCaseをバルクインサート
	if len(newTestRunCases) > 0 {
		if result := h.DB.Create(&newTestRunCases); result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}
	}

	// 不要になったTestRunCaseを削除
	for _, existingCase := range existingTestRunCases {
		if _, existsInRequest := findInSlice(req.TestCaseIDs, existingCase.TestCaseID); !existsInRequest {
			h.DB.Delete(&existingCase)
		}
	}

	c.JSON(http.StatusCreated, newTestRunCases)
}

func findInSlice(slice []uint, value uint) (int, bool) {
	for i, item := range slice {
		if item == value {
			return i, true
		}
	}
	return -1, false
}

func (h *TestRunHandler) PostTestRunCaseComment(c *gin.Context) {
	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.DB.Create(&comment); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	var newComment model.Comment
	if result := h.DB.Preload("CreatedBy").Preload("UpdatedBy").First(&newComment, comment.ID); result.Error != nil {
		handleError(c, http.StatusInternalServerError, "Failed to retrieve records", result.Error)
		return
	}

	var status *util.Status
	if comment.Status != nil {
		status = &util.Status{ID: comment.Status.ID, Name: comment.Status.Name, Color: comment.Status.Color}
	} else {
		status = nil
	}
	commentJson := util.Comment{
		ID:        newComment.ID,
		Status:    status,
		Content:   newComment.Content,
		CreatedBy: convertUserToJSON(newComment.CreatedBy),
		UpdatedBy: convertUserToJSON(newComment.UpdatedBy),
	}

	c.JSON(http.StatusCreated, commentJson)
}

func (h *TestRunHandler) PutTestRun(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updatedTestRun model.TestRun
	if err := c.ShouldBindJSON(&updatedTestRun); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.DB.Model(&model.TestRun{}).Where("id = ?", id).Updates(updatedTestRun); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTestRun)
}

func (h *TestRunHandler) PutTestRunCase(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var updatedTestRunCase model.TestRunCase
	if err := c.ShouldBindJSON(&updatedTestRunCase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.DB.Model(&model.TestRunCase{}).Where("id = ?", id).Updates(updatedTestRunCase); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedTestRunCase)
}

func (h *TestRunHandler) PutTestRunCaseComment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	var comment model.Comment
	if err := c.ShouldBindJSON(&comment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := h.DB.Model(&model.Comment{}).Where("id = ?", id).Updates(comment); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(http.StatusOK, comment)
}

func (h *TestRunHandler) DeleteTestRunCase(c *gin.Context) {
	// URLからidを取得
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 指定されたIDを持つTestRunCaseを削除
	if result := h.DB.Delete(&model.TestRunCase{}, id); result.Error != nil {
		// データベースエラーの場合
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 削除が成功した場合、204 No Contentステータスを返す
	c.Status(http.StatusNoContent)
}

func (h *TestRunHandler) DeleteTestRun(c *gin.Context) {
	// URLからidを取得
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 指定されたIDを持つTestRunCaseを削除
	if result := h.DB.Delete(&model.TestRun{}, id); result.Error != nil {
		// データベースエラーの場合
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 削除が成功した場合、204 No Contentステータスを返す
	c.Status(http.StatusNoContent)
}

func (h *TestRunHandler) DeleteTestRunCaseComment(c *gin.Context) {
	// URLからidを取得
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID format"})
		return
	}

	// 指定されたIDを持つTestRunCaseを削除
	if result := h.DB.Delete(&model.Comment{}, id); result.Error != nil {
		// データベースエラーの場合
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	// 削除が成功した場合、204 No Contentステータスを返す
	c.Status(http.StatusNoContent)
}

func convertToJSONTestRunCase(trc model.TestRunCase) util.TestRunCase {
	commentsJSON := []util.Comment{}
	for _, comment := range trc.Comments {
		commentsJSON = append(commentsJSON, convertCommentToJSON(comment))
	}

	var assignedToJSON *util.User
	if trc.AssignedTo != nil {
		userJSON := convertUserToJSON(*trc.AssignedTo)
		assignedToJSON = &userJSON // userJSON のアドレスを assignedToJSON に代入
	}
	return util.TestRunCase{
		ID:         trc.ID,
		TestCaseId: trc.TestCaseID,
		Title:      trc.TestCase.Title,
		Content:    trc.TestCase.Content,
		Comments:   commentsJSON,
		Status: util.Status{
			ID:    trc.Status.ID,
			Name:  trc.Status.Name,
			Color: trc.Status.Color,
		},
		AssignedTo: assignedToJSON,
	}
}

func convertUserToJSON(user model.User) util.User {
	return util.User{
		ID:   user.ID,
		Name: user.Name,
	}
}
