package handler_test

import (
	"backend/handler"
	"backend/util"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupMockStatusHandler() (*handler.StatusHandler, sqlmock.Sqlmock) {
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

	memberHandler := &handler.StatusHandler{
		DB: gormDB,
	}

	return memberHandler, mock
}

func TestGetStatuses(t *testing.T) {
	h, mock := setupMockStatusHandler()
	gin.SetMode(gin.TestMode)

	rows := sqlmock.NewRows([]string{"id", "name", "color"})
	rows.AddRow(1, "Passed", "green")
	mock.ExpectQuery("^SELECT \\* FROM `statuses`").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/protected/statuses", nil)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/protected/statuses", h.GetStatues)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response util.StatusesResponseData
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Statuses, 1)
	assert.Equal(t, "Passed", response.Statuses[0].Name)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
