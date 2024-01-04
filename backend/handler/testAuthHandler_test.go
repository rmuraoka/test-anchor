package handler_test

import (
	"backend/handler"
	"backend/util"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupMockAuthHandler() (*handler.AuthHandler, sqlmock.Sqlmock) {
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

	return handler.NewAuthHandler(gormDB), mock
}

func TestPostLogin(t *testing.T) {
	h, mock := setupMockAuthHandler()
	gin.SetMode(gin.TestMode)

	// モックユーザーデータの設定
	email := "test@example.com"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	userRows := sqlmock.NewRows([]string{"id", "email", "password", "name", "language"}).
		AddRow(1, email, hashedPassword, "Test User", "en")
	mock.ExpectQuery("^SELECT \\* FROM `users`").
		WithArgs(email).
		WillReturnRows(userRows)

	// HTTPリクエストの設定
	loginInput := handler.LoginInput{Email: email, Password: "password"}
	body, _ := json.Marshal(loginInput)
	r := gin.Default()
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// テスト実行
	r.POST("/login", h.PostLogin)
	r.ServeHTTP(w, req)

	// レスポンスの検証
	assert.Equal(t, http.StatusOK, w.Code)
	var responseData util.Auth
	err := json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, uint(1), responseData.User.ID)
}
