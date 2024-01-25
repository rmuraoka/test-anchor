package handler_test

import (
	"backend/handler"
	"backend/model"
	"backend/util"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupMockTestCaseHandler() (*handler.TestCaseHandler, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
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

	return handler.NewTestCaseHandler(gormDB), mock
}

func TestGetTestCases(t *testing.T) {
	h, mock := setupMockTestCaseHandler()
	gin.SetMode(gin.TestMode)

	// プロジェクトIDを設定
	projectID := 1
	projectCode := "PROJECT"

	// プロジェクトコードに基づいてプロジェクトを検索するクエリを期待
	rows := sqlmock.NewRows([]string{"id", "code"}).
		AddRow(projectID, projectCode)
	mock.ExpectQuery("^SELECT \\* FROM `projects`").
		WithArgs(projectCode).
		WillReturnRows(rows)

	// テストスイートのクエリを期待
	testSuiteRows := sqlmock.NewRows([]string{"id", "project_id", "name", "parent_id"}).
		AddRow(1, projectID, "TestSuites 1", nil)
	mock.ExpectQuery("^SELECT \\* FROM `test_suites`").
		WithArgs(projectID).
		WillReturnRows(testSuiteRows)
	// テストケースのクエリを期待
	testCaseRows := sqlmock.NewRows([]string{"id", "test_suite_id", "title", "content"}).
		AddRow(1, 1, "Test Case 1", "Test Content")
	mock.ExpectQuery("^SELECT \\* FROM `test_cases`").
		WithArgs(1).
		WillReturnRows(testCaseRows)
	mock.ExpectQuery("^SELECT \\* FROM `test_suites` WHERE `test_suites`.`parent_id` = ?").
		WithArgs(1).
		WillReturnRows(testSuiteRows)

	// Set up HTTP request
	r := gin.Default()
	r.GET("/protected/:project_code/cases", h.GetTestCases)
	req, _ := http.NewRequest("GET", "/protected/"+projectCode+"/cases", nil)
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check expectations
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスデータを検証
	var responseData util.TestCasesResponseData
	err = json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, uint(projectID), responseData.ProjectID)
	assert.Len(t, responseData.TestSuites, 1)
	assert.Equal(t, "TestSuites 1", responseData.TestSuites[0].Name)
}

func TestPostTestCase(t *testing.T) {
	h, mock := setupMockTestCaseHandler()
	gin.SetMode(gin.TestMode)

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `test_cases`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	newTestCase := model.TestCase{
		Title:   "New Test Case",
		Content: "New Content",
	}
	requestBody, _ := json.Marshal(newTestCase)
	r := gin.Default()
	r.POST("protected/cases", h.PostTestCase)
	req, _ := http.NewRequest("POST", "/protected/cases", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetTestCase(t *testing.T) {
	h, mock := setupMockTestCaseHandler()
	gin.SetMode(gin.TestMode)

	// Set up expected queries and responses
	id := 1
	rows := sqlmock.NewRows([]string{"id", "title", "content"}).
		AddRow(id, "Test Case 1", "Test Content")

	mock.ExpectQuery("^SELECT \\* FROM `test_cases`").
		WithArgs(id).
		WillReturnRows(rows)

	// Set up HTTP request
	r := gin.Default()
	r.GET("/protected/cases/:id", h.GetTestCase)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/protected/cases/%d", id), nil)
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check expectations
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
	var resp model.TestCase
	err = json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "Test Case 1", resp.Title)
}

// TestPutTestCase - テストケースを更新するテスト
func TestPutTestCase(t *testing.T) {
	h, mock := setupMockTestCaseHandler()
	gin.SetMode(gin.TestMode)

	id := 1

	// SELECT クエリを期待
	rows := sqlmock.NewRows([]string{"id", "title", "content"}).
		AddRow(id, "Original Title", "Original Content")
	mock.ExpectQuery("^SELECT \\* FROM `test_cases`").
		WithArgs(id).
		WillReturnRows(rows)
	// トランザクションの開始を期待
	mock.ExpectBegin()
	// UPDATE クエリを期待
	mock.ExpectExec("^UPDATE `test_cases`").
		WithArgs(sqlmock.AnyArg(), "Updated Title", "Updated Content", id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	// トランザクションのコミットを期待
	mock.ExpectCommit()

	mock.ExpectBegin()
	// UPDATE クエリを期待
	mock.ExpectExec("^UPDATE `test_cases`").
		WithArgs(sqlmock.AnyArg(), nil, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	// トランザクションのコミットを期待
	mock.ExpectCommit()

	updatedData := model.TestCase{
		Title:   "Updated Title",
		Content: "Updated Content",
	}
	requestBody, _ := json.Marshal(updatedData)

	r := gin.Default()
	r.PUT("/protected/cases/:id", h.PutTestCase)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/protected/cases/%d", id), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

// TestDeleteTestCase - テストケースを削除するテスト
func TestDeleteTestCase(t *testing.T) {
	h, mock := setupMockTestCaseHandler()
	gin.SetMode(gin.TestMode)

	id := 1

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `test_cases`").
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `test_run_cases`").
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	r := gin.Default()
	r.DELETE("/protected/cases/:id", h.DeleteTestCase)
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/protected/cases/%d", id), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPostTestSuite(t *testing.T) {
	h, mock := setupMockTestCaseHandler()
	gin.SetMode(gin.TestMode)

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `test_suites`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	newTestCase := model.TestSuite{
		Name:     "New Test Case",
		ParentID: nil,
	}
	requestBody, _ := json.Marshal(newTestCase)
	r := gin.Default()
	r.POST("/protected/suites", h.PostTestSuite)
	req, _ := http.NewRequest("POST", "/protected/suites", bytes.NewBuffer(requestBody))
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPutTestSuite(t *testing.T) {
	h, mock := setupMockTestCaseHandler()
	gin.SetMode(gin.TestMode)

	id := 1

	// SELECT クエリを期待
	rows := sqlmock.NewRows([]string{"id", "name", "parent_id"}).
		AddRow(id, "TestSuites Name", nil)
	mock.ExpectQuery("^SELECT \\* FROM `test_suites`").
		WithArgs(id).
		WillReturnRows(rows)
	// トランザクションの開始を期待
	mock.ExpectBegin()
	// UPDATE クエリを期待
	mock.ExpectExec("^UPDATE `test_suites`").
		WithArgs(sqlmock.AnyArg(), "Updated Title", id).
		WillReturnResult(sqlmock.NewResult(1, 1))

	// トランザクションのコミットを期待
	mock.ExpectCommit()

	updatedData := model.TestSuite{
		Name:     "Updated Title",
		ParentID: nil,
	}
	requestBody, _ := json.Marshal(updatedData)

	r := gin.Default()
	r.PUT("/protected/suites/:id", h.PutTestSuite)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/protected/suites/%d", id), bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTestSuites(t *testing.T) {
	h, mock := setupMockTestCaseHandler()
	gin.SetMode(gin.TestMode)

	id := 1

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `test_suites`").
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `test_cases`").
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `test_run_cases`").
		WithArgs(sqlmock.AnyArg(), id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	r := gin.Default()
	r.DELETE("/protected/suites/:id", h.DeleteTestSuite)
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/protected/suites/%d", id), nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPutTestCasesBulk(t *testing.T) {
	h, mock := setupMockTestCaseHandler()
	gin.SetMode(gin.TestMode)

	projectID := 1
	projectCode := "PROJECT"
	rows := sqlmock.NewRows([]string{"id", "code"}).
		AddRow(projectID, projectCode)
	mock.ExpectQuery("^SELECT \\* FROM `projects`").
		WithArgs(projectCode).
		WillReturnRows(rows)

	// モックデータの設定
	testCaseRequestRows := []handler.TestCasePutRequestRow{}
	testCaseRequestRows = append(testCaseRequestRows, handler.TestCasePutRequestRow{TestCaseID: uint(1), OrderIndex: 0})
	requestBody := handler.TestCasePutRequest{
		TestSuiteID:         uint(1),
		TestCaseRequestRows: testCaseRequestRows,
	}

	// TestRunCasesのINSERTクエリを模擬
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `test_cases`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), projectID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// HTTPリクエストの設定
	r := gin.Default()
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/protected/"+projectCode+"/cases/bulk", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// テスト実行
	r.POST("/protected/:project_code/cases/bulk", h.PutTestCaseBulk)
	r.ServeHTTP(w, req)

	// レスポンスと期待される結果を検証
	assert.Equal(t, http.StatusOK, w.Code)

	// モックの期待が満たされたか検証
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPutTestSuitesBulk(t *testing.T) {
	h, mock := setupMockTestCaseHandler()
	gin.SetMode(gin.TestMode)

	projectID := 1
	projectCode := "PROJECT"
	rows := sqlmock.NewRows([]string{"id", "code"}).
		AddRow(projectID, projectCode)
	mock.ExpectQuery("^SELECT \\* FROM `projects`").
		WithArgs(projectCode).
		WillReturnRows(rows)

	// モックデータの設定
	testCaseRequestRows := []handler.TestSuiteRequestRow{}
	testCaseRequestRows = append(testCaseRequestRows, handler.TestSuiteRequestRow{TestSuiteID: uint(1), OrderIndex: 0})
	requestBody := handler.TestSuiteRequest{
		ParentID:            nil,
		TestCaseRequestRows: testCaseRequestRows,
	}

	// TestRunCasesのINSERTクエリを模擬
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `test_suites`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), projectID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// HTTPリクエストの設定
	r := gin.Default()
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/protected/"+projectCode+"/suites/bulk", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// テスト実行
	r.POST("/protected/:project_code/suites/bulk", h.PutTestSuiteBulk)
	r.ServeHTTP(w, req)

	// レスポンスと期待される結果を検証
	assert.Equal(t, http.StatusOK, w.Code)

	// モックの期待が満たされたか検証
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPostTestCasesBulk(t *testing.T) {
	h, mock := setupMockTestCaseHandler()
	gin.SetMode(gin.TestMode)

	projectID := 1
	projectCode := "PROJECT"
	rows := sqlmock.NewRows([]string{"id", "code"}).
		AddRow(projectID, projectCode)
	mock.ExpectQuery("^SELECT \\* FROM `projects`").
		WithArgs(projectCode).
		WillReturnRows(rows)

	// モックデータの設定
	testCasePostRequests := []handler.TestCasePostRequest{}
	testCasePostRequests = append(testCasePostRequests, handler.TestCasePostRequest{TestSuiteName: "testSuite", Title: "title", Content: "content"})
	requestBody := handler.TestCasesPostRequest{
		TestCases: testCasePostRequests,
	}

	testSuiteRows := sqlmock.NewRows([]string{"id", "name"})
	mock.ExpectQuery("^SELECT \\* FROM `test_suites`").
		WithArgs(sqlmock.AnyArg(), projectID, projectID, sqlmock.AnyArg()).
		WillReturnRows(testSuiteRows)

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `test_suites`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `test_cases`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// HTTPリクエストの設定
	r := gin.Default()
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/protected/"+projectCode+"/cases/bulk", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// テスト実行
	r.POST("/protected/:project_code/cases/bulk", h.PostTestCasesBulk)
	r.ServeHTTP(w, req)

	// レスポンスと期待される結果を検証
	assert.Equal(t, http.StatusOK, w.Code)

	// モックの期待が満たされたか検証
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
