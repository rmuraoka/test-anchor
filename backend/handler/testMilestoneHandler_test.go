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

func setupMockMilestoneHandler() (*handler.MilestoneHandler, sqlmock.Sqlmock) {
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

	return handler.NewMilestoneHandler(gormDB), mock
}

func TestGetMilestone(t *testing.T) {
	h, mock := setupMockMilestoneHandler()
	gin.SetMode(gin.TestMode)

	milestoneID := 1
	rows := sqlmock.NewRows([]string{"id", "title", "description"}).
		AddRow(milestoneID, "Test Milestone 1", "Description")
	mock.ExpectQuery("^SELECT \\* FROM `milestones`").
		WithArgs(milestoneID).
		WillReturnRows(rows)

	req := httptest.NewRequest("GET", fmt.Sprintf("/protected/milestones/%d", milestoneID), nil)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/protected/milestones/:id", h.GetMilestone)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPostMilestone(t *testing.T) {
	h, mock := setupMockMilestoneHandler()
	gin.SetMode(gin.TestMode)

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `milestones`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "New Milestone", "New Description", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := strings.NewReader(`{"project_id": 1, "title": "New Milestone", "description": "New Description", "status": "Active", "due_date": "2024-01-01"}`)
	req := httptest.NewRequest("POST", "/protected/milestones", body)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/protected/milestones", h.PostMilestone)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPutMilestone(t *testing.T) {
	h, mock := setupMockMilestoneHandler()
	gin.SetMode(gin.TestMode)

	milestoneID := 1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `milestones`").
		WithArgs(sqlmock.AnyArg(), "Updated Milestone", "Updated Description", "Active", sqlmock.AnyArg(), milestoneID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	body := strings.NewReader(`{"title": "Updated Milestone", "description": "Updated Description", "status": "Active", "due_date": "2024-01-01"}`)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/protected/milestones/%d", milestoneID), body)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.PUT("/protected/milestones/:id", h.PutMilestone)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteMilestone(t *testing.T) {
	h, mock := setupMockMilestoneHandler()
	gin.SetMode(gin.TestMode)

	milestoneID := 1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `milestones`").
		WithArgs(sqlmock.AnyArg(), milestoneID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	req := httptest.NewRequest("DELETE", fmt.Sprintf("/protected/milestones/%d", milestoneID), nil)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.DELETE("/protected/milestones/:id", h.DeleteMilestone)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
