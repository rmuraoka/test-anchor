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
	"time"
)

func setupMockProjectHandler() (*handler.ProjectHandler, sqlmock.Sqlmock) {
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

	return handler.NewProjectHandler(gormDB), mock
}

func TestGetProject(t *testing.T) {
	h, mock := setupMockProjectHandler()
	gin.SetMode(gin.TestMode)

	ProjectID := 1
	ProjectCode := "CODE"
	createdByUserID := 1
	updatedByUserID := 2
	rows := sqlmock.NewRows([]string{"id", "code", "title", "description"}).
		AddRow(ProjectID, ProjectCode, "Test Project 1", "Description")
	mock.ExpectQuery("^SELECT \\* FROM `projects`").
		WithArgs(ProjectCode).
		WillReturnRows(rows)

	milestoneRow := sqlmock.NewRows([]string{"id", "project_id", "title", "description", "status", "due_date"}).
		AddRow(1, ProjectID, "Milestone Title", "Milestone Description", "Open", time.Now())
	mock.ExpectQuery("^SELECT \\* FROM `milestones`").
		WithArgs(ProjectID).
		WillReturnRows(milestoneRow)

	testPlanRow := sqlmock.NewRows([]string{"id", "project_id", "test_plan_id", "title", "status", "started_at", "completed_at", "created_by_id", "updated_by_id"}).
		AddRow(1, ProjectID, nil, "Test Plan Title", "Not Started", time.Now(), time.Now(), 1, 2)
	mock.ExpectQuery("^SELECT \\* FROM `test_plans`").
		WithArgs(ProjectID).
		WillReturnRows(testPlanRow)

	createdByUserRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(createdByUserID, "Created User Name")
	mock.ExpectQuery("^SELECT \\* FROM `users`").
		WithArgs(createdByUserID).
		WillReturnRows(createdByUserRows)

	updatedByUserRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(updatedByUserID, "Updated User Name")
	mock.ExpectQuery("^SELECT \\* FROM `users`").
		WithArgs(updatedByUserID).
		WillReturnRows(updatedByUserRows)

	req := httptest.NewRequest("GET", fmt.Sprintf("/protected/projects/%s", ProjectCode), nil)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/protected/projects/:project_code", h.GetProject)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPostProject(t *testing.T) {
	h, mock := setupMockProjectHandler()
	gin.SetMode(gin.TestMode)

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `projects`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "CODE", "New Project", "New Description").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := strings.NewReader(`{"code": "CODE", "title": "New Project", "description": "New Description"}`)
	req := httptest.NewRequest("POST", "/protected/projects", body)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/protected/projects", h.PostProject)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPutProject(t *testing.T) {
	h, mock := setupMockProjectHandler()
	gin.SetMode(gin.TestMode)

	ProjectID := 1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `projects`").
		WithArgs(sqlmock.AnyArg(), "Updated Project", "Updated Description", ProjectID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := strings.NewReader(`{"title": "Updated Project", "description": "Updated Description"}`)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/protected/projects/%d", ProjectID), body)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.PUT("/protected/projects/:id", h.PutProject)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteProject(t *testing.T) {
	h, mock := setupMockProjectHandler()
	gin.SetMode(gin.TestMode)

	ProjectID := 1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `projects`").
		WithArgs(sqlmock.AnyArg(), ProjectID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/protected/projects/%d", ProjectID), nil)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.DELETE("/protected/projects/:id", h.DeleteProject)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
