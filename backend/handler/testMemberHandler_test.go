package handler_test

import (
	"backend/handler"
	"backend/model"
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
	"strings"
	"testing"
)

type MockEmailSender struct {
	SentEmails []string
}

func (m *MockEmailSender) SendMail(to []string, subject, body string) error {
	m.SentEmails = append(m.SentEmails, subject)
	return nil
}

func setupMockMemberHandler() (*handler.MemberHandler, sqlmock.Sqlmock, *MockEmailSender) {
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

	mockEmailSender := &MockEmailSender{}
	memberHandler := &handler.MemberHandler{
		DB:          gormDB,
		EmailSender: mockEmailSender, // モックされたメール送信機能を注入
	}

	return memberHandler, mock, mockEmailSender
}

func TestGetMembers(t *testing.T) {
	h, mock, _ := setupMockMemberHandler()
	gin.SetMode(gin.TestMode)

	rows := sqlmock.NewRows([]string{"id", "email", "name", "status"})
	rows.AddRow(1, "john.doe@example.com", "John Doe", "active")
	mock.ExpectQuery("^SELECT \\* FROM `users`").WillReturnRows(rows)

	req := httptest.NewRequest("GET", "/protected/members", nil)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.GET("/protected/members", h.GetMembers)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var response util.MembersResponseData
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Len(t, response.Members, 1)
	assert.Equal(t, "John Doe", response.Members[0].Name)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPostMember(t *testing.T) {
	h, mock, mockEmailSender := setupMockMemberHandler()
	gin.SetMode(gin.TestMode)

	// データベースの操作をモック
	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `users`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "Jane Doe", "jane.doe@example.com", sqlmock.AnyArg(), "active", sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// HTTPリクエストとレコーダーの設定
	body := strings.NewReader(`{"email": "jane.doe@example.com", "name": "Jane Doe", "status": "active"}`)
	req := httptest.NewRequest("POST", "/protected/members", body)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.POST("/protected/members", h.PostMember)
	router.ServeHTTP(w, req)

	// 検証
	assert.Equal(t, http.StatusCreated, w.Code)
	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe", response.Name)
	assert.Len(t, mockEmailSender.SentEmails, 1)

	// 期待されるデータベースの操作が行われたか確認
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPutMember(t *testing.T) {
	h, mock, _ := setupMockMemberHandler()
	gin.SetMode(gin.TestMode)

	// データベースの操作をモック
	userID := 1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `users`").
		WithArgs(sqlmock.AnyArg(), "Jane Doe Updated", "active", userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// HTTPリクエストとレコーダーの設定
	body := strings.NewReader(`{"name": "Jane Doe Updated", "status": "active"}`)
	req := httptest.NewRequest("PUT", fmt.Sprintf("/protected/members/%d", userID), body)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.PUT("/protected/members/:id", h.PutMember)
	router.ServeHTTP(w, req)

	// 検証
	assert.Equal(t, http.StatusOK, w.Code)
	var response model.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, "Jane Doe Updated", response.Name)

	// 期待されるデータベースの操作が行われたか確認
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

// TestDeleteMember - メンバーを削除するエンドポイントのテスト
func TestDeleteMember(t *testing.T) {
	h, mock, _ := setupMockMemberHandler()
	gin.SetMode(gin.TestMode)

	// データベースの操作をモック
	userID := 1
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `users`").
		WithArgs(sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// HTTPリクエストとレコーダーの設定
	req := httptest.NewRequest("DELETE", fmt.Sprintf("/protected/members/%d", userID), nil)
	w := httptest.NewRecorder()
	router := gin.Default()
	router.DELETE("/protected/members/:id", h.DeleteMember)
	router.ServeHTTP(w, req)

	// 検証
	assert.Equal(t, http.StatusNoContent, w.Code)

	// 期待されるデータベースの操作が行われたか確認
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
