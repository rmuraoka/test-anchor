package handler_test

import (
	"backend/handler"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupMockTestPlanHandler() (*handler.TestPlanHandler, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherRegexp))
	if err != nil {
		panic(fmt.Sprintf("An error '%s' was not expected when opening a stub database connection", err))
	}

	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn:                      db,
		SkipInitializeWithVersion: true,
	}), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("An error '%s' was not expected when opening a gorm database connection", err))
	}

	return handler.NewTestPlanHandler(gormDB), mock
}

func TestGetTestPlan(t *testing.T) {
	h, mock := setupMockTestPlanHandler()
	gin.SetMode(gin.TestMode)

	testPlanID := 1
	projectID := 2
	CreatedUserId := 3
	UpdatedUserId := 3
	testRunId := 4
	statusId := 5
	rows := sqlmock.NewRows([]string{"id", "project_id", "title", "description"}).
		AddRow(testPlanID, projectID, "Test Plan 1", "Description")
	mock.ExpectQuery("^SELECT \\* FROM `test_plans`").
		WithArgs(testPlanID).
		WillReturnRows(rows)

	// TestRunに関連するデータをモック
	testRunRows := sqlmock.NewRows([]string{"id", "project_id", "test_plan_id", "title", "status", "created_by_id", "updated_by_id"}).
		AddRow(testRunId, projectID, testPlanID, "Test Run 1", "In Progress", CreatedUserId, UpdatedUserId)
	mock.ExpectQuery("^SELECT \\* FROM `test_runs`").
		WithArgs(testPlanID).
		WillReturnRows(testRunRows)

	// TestRunCaseテーブルからのクエリを期待
	testRunCaseRows := sqlmock.NewRows([]string{"id", "test_run_id", "test_case_id", "assigned_to_id", "status_id"}).
		AddRow(1, testRunId, 1, 3, statusId)
	mock.ExpectQuery("^SELECT \\* FROM `test_run_cases`").
		WithArgs(testRunId).
		WillReturnRows(testRunCaseRows)

	statusRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(statusId, "Status Name")
	mock.ExpectQuery("^SELECT \\* FROM `statuses`").
		WithArgs(statusId).
		WillReturnRows(statusRows)

	req := httptest.NewRequest("GET", fmt.Sprintf("/protected/plans/%d", testPlanID), nil)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/protected/plans/:id", h.GetTestPlan)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPostTestPlan(t *testing.T) {
	h, mock := setupMockTestPlanHandler()
	gin.SetMode(gin.TestMode)

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `test_plans`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1, "Active", "New Plan", sqlmock.AnyArg(), sqlmock.AnyArg(), 1, 1).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := strings.NewReader(`{"project_id": 1, "title": "New Plan", "status": "Active", "started_at": null, "completed_at": null, "created_by_id": 1, "updated_by_id": 1}`)
	req := httptest.NewRequest("POST", "/plans", body)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/plans", h.PostTestPlan)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPutTestPlan(t *testing.T) {
	h, mock := setupMockTestPlanHandler()
	gin.SetMode(gin.TestMode)

	testPlanID := 1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `test_plans`").
		WithArgs(sqlmock.AnyArg(), "Active", "Updated Plan", 1, testPlanID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := strings.NewReader(`{"title": "Updated Plan", "status": "Active", "updated_by_id": 1}`)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/protected/plans/%d", testPlanID), body)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.PUT("/protected/plans/:id", h.PutTestPlan)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeletePlan(t *testing.T) {
	h, mock := setupMockTestPlanHandler()
	gin.SetMode(gin.TestMode)

	testPlanID := 1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `test_plans`").
		WithArgs(sqlmock.AnyArg(), testPlanID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/protected/plans/%d", testPlanID), nil)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.DELETE("/protected/plans/:id", h.DeletePlan)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
