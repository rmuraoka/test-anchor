package handler_test

import (
	"backend/handler"
	"backend/model"
	"backend/util"
	"bytes"
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

func setupMockTestRunHandler() (*handler.TestRunHandler, sqlmock.Sqlmock) {
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

	return handler.NewTestRunHandler(gormDB), mock
}

func TestGetTestRuns(t *testing.T) {
	h, mock := setupMockTestRunHandler()
	gin.SetMode(gin.TestMode)

	projectCode := "PROJECT"
	projectID := 1

	// プロジェクトコードに基づいてプロジェクトを検索するクエリを期待
	rows := sqlmock.NewRows([]string{"id", "code"}).
		AddRow(projectID, projectCode)
	mock.ExpectQuery("^SELECT \\* FROM `projects`").
		WithArgs(projectCode).
		WillReturnRows(rows)

	// TestRunテーブルからのクエリを期待
	testRunRows := sqlmock.NewRows([]string{"id", "project_id", "status"}).
		AddRow(1, projectID, "In Progress")
	mock.ExpectQuery("^SELECT \\* FROM `test_runs`").
		WithArgs(projectID).
		WillReturnRows(testRunRows)

	// Set up HTTP request
	r := gin.Default()
	r.GET("/protected/:project_code/runs", h.GetTestRuns)
	req, _ := http.NewRequest("GET", "/protected/"+projectCode+"/runs", nil)
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check expectations
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetTestRunCase(t *testing.T) {
	h, mock := setupMockTestRunHandler()
	gin.SetMode(gin.TestMode)

	testRunCaseID := 1
	statusID := 2
	assignedToID := 3
	createdByUserID := 4
	updatedByUserID := 5
	testCaseID := 6
	commentID := 7

	// TestRunCaseテーブルからのクエリを期待
	testRunCaseRows := sqlmock.NewRows([]string{"id", "status_id", "assigned_to_id", "test_case_id"}).
		AddRow(testRunCaseID, statusID, assignedToID, testCaseID)
	mock.ExpectQuery("^SELECT \\* FROM `test_run_cases`").
		WithArgs(testRunCaseID).
		WillReturnRows(testRunCaseRows)

	// AssignedToテーブルからのクエリを期待
	assignedToRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(assignedToID, "Assigned User")
	mock.ExpectQuery("^SELECT \\* FROM `users`").
		WithArgs(assignedToID).
		WillReturnRows(assignedToRows)

	commentRows := sqlmock.NewRows([]string{"id", "content", "test_run_case_id"}).
		AddRow(commentID, "Content", testRunCaseID)
	mock.ExpectQuery("^SELECT \\* FROM `comments`").
		WithArgs(testRunCaseID).
		WillReturnRows(commentRows)

	// Statusテーブルからのクエリを期待
	statusRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(statusID, "Untested")
	mock.ExpectQuery("^SELECT \\* FROM `statuses`").
		WithArgs(statusID).
		WillReturnRows(statusRows)

	// TestCaseテーブルからのクエリを期待
	testCaseRows := sqlmock.NewRows([]string{"id", "title", "content", "created_by_id", "updated_by_id"}).
		AddRow(testCaseID, "Test Case Title", "Test Content", createdByUserID, updatedByUserID)
	mock.ExpectQuery("^SELECT \\* FROM `test_cases`").
		WithArgs(testCaseID).
		WillReturnRows(testCaseRows)

	// CreatedByとUpdatedByユーザーのクエリも同様に期待
	createdByUserRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(createdByUserID, "Creator User")
	updatedByUserRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(updatedByUserID, "Updater User")
	mock.ExpectQuery("^SELECT \\* FROM `users`").
		WithArgs(createdByUserID).
		WillReturnRows(createdByUserRows)
	mock.ExpectQuery("^SELECT \\* FROM `users`").
		WithArgs(updatedByUserID).
		WillReturnRows(updatedByUserRows)

	// Set up HTTP request
	r := gin.Default()
	r.GET("/protected/runs/cases/:id", h.GetTestRunCase)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/protected/runs/cases/%d", testRunCaseID), nil)
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check expectations
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスデータを検証
	var responseData util.TestRunCase
	err = json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	// 追加: レスポンスデータの検証
	assert.Equal(t, uint(testRunCaseID), responseData.ID)
	assert.Equal(t, "Untested", responseData.Status.Name)
	assert.Equal(t, "Assigned User", responseData.AssignedTo.Name)
	assert.Equal(t, "Test Case Title", responseData.Title)
	assert.Equal(t, "Test Content", responseData.Content)
}

func TestGetTestRunCases(t *testing.T) {
	h, mock := setupMockTestRunHandler()
	gin.SetMode(gin.TestMode)

	testRunID := 1
	testCaseID := 2
	testSuiteID := 3
	statusID := 4
	assignedToID := 5
	createdByUserID := 6
	updatedByUserID := 7
	testRunCaseID := 8
	commentID := 9
	projectID := 10

	testRunRow := sqlmock.NewRows([]string{"id", "title", "project_id"}).
		AddRow(testRunID, "Test Run Title", projectID)
	mock.ExpectQuery("^SELECT \\* FROM `test_runs`").
		WithArgs(testRunID).
		WillReturnRows(testRunRow)

	testRunCaseStatusRows := sqlmock.NewRows([]string{"id", "test_run_id", "test_case_id", "status_id", "assigned_to_id"}).
		AddRow(testRunCaseID, testRunID, testCaseID, statusID, assignedToID)
	mock.ExpectQuery("^SELECT \\* FROM `test_run_cases`").
		WithArgs(testRunID).
		WillReturnRows(testRunCaseStatusRows)

	testRunStatusRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(statusID, "Status Name")
	mock.ExpectQuery("^SELECT \\* FROM `statuses`").
		WithArgs(statusID).
		WillReturnRows(testRunStatusRows)

	// TestRunCaseテーブルからのクエリを期待
	testRunCaseRows := sqlmock.NewRows([]string{"id", "test_run_id", "test_case_id", "status_id", "assigned_to_id"}).
		AddRow(testRunCaseID, testRunID, testCaseID, statusID, assignedToID)
	mock.ExpectQuery("^SELECT \\* FROM `test_run_cases`").
		WithArgs(testRunID).
		WillReturnRows(testRunCaseRows)

	// AssignedToユーザーのクエリを期待
	assignedToRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(assignedToID, "Assigned User Name")
	mock.ExpectQuery("^SELECT \\* FROM `users`").
		WithArgs(assignedToID).
		WillReturnRows(assignedToRows)

	commentRows := sqlmock.NewRows([]string{"id", "content", "test_run_case_id"}).
		AddRow(commentID, "Content", testRunCaseID)
	mock.ExpectQuery("^SELECT \\* FROM `comments`").
		WithArgs(testRunCaseID).
		WillReturnRows(commentRows)

	// Statusテーブルからのクエリを期待
	statusRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(statusID, "Status Name")
	mock.ExpectQuery("^SELECT \\* FROM `statuses`").
		WithArgs(statusID).
		WillReturnRows(statusRows)

	// TestCaseテーブルからのクエリを期待
	testCaseRows := sqlmock.NewRows([]string{"id", "test_suite_id", "created_by_id", "updated_by_id", "title", "content"}).
		AddRow(testCaseID, testSuiteID, createdByUserID, updatedByUserID, "Test Case Title", "Test Case Content")
	mock.ExpectQuery("^SELECT \\* FROM `test_cases`").
		WithArgs(testCaseID).
		WillReturnRows(testCaseRows)

	// TestSuiteテーブルからのクエリを期待
	testSuiteRows := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(testSuiteID, "TestSuites Name")
	mock.ExpectQuery("^SELECT \\* FROM `test_suites`").
		WithArgs(testSuiteID).
		WillReturnRows(testSuiteRows)

	// TestSuiteテーブルからのクエリを期待
	testSuiteNewRows := sqlmock.NewRows([]string{"id", "name", "project_id", "parent_id"}).
		AddRow(3, "TestSuites Name", 1, nil)
	mock.ExpectQuery("^SELECT \\* FROM `test_suites`").
		WithArgs(testSuiteID).
		WillReturnRows(testSuiteNewRows)

	mock.ExpectQuery("^SELECT \\* FROM `test_suites`").
		WithArgs(projectID).
		WillReturnRows(testSuiteNewRows)

	// Set up HTTP request
	r := gin.Default()
	r.GET("/protected/runs/:id/cases", h.GetTestRunCases)
	req, _ := http.NewRequest("GET", fmt.Sprintf("/protected/runs/%d/cases", testRunID), nil)
	w := httptest.NewRecorder()

	// Perform the request
	r.ServeHTTP(w, req)

	// Check expectations
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, w.Code)

	// レスポンスデータを検証
	var responseData util.TestRunCasesResponseData
	err = json.Unmarshal(w.Body.Bytes(), &responseData)
	assert.NoError(t, err)
	assert.Equal(t, uint(testRunID), responseData.TestRunID)
	assert.Len(t, responseData.TestSuites, 1)
	assert.Equal(t, "TestSuites Name", responseData.TestSuites[0].Name)
	assert.Len(t, responseData.TestSuites[0].TestCases, 1)
	assert.Equal(t, "Test Case Title", responseData.TestSuites[0].TestCases[0].Title)
	assert.Equal(t, "Test Case Content", responseData.TestSuites[0].TestCases[0].Content)
	assert.Equal(t, "Status Name", responseData.TestSuites[0].TestCases[0].Status.Name)
	assert.Equal(t, "Assigned User Name", responseData.TestSuites[0].TestCases[0].AssignedTo.Name)
}

func TestPostTestRun(t *testing.T) {
	h, mock := setupMockTestRunHandler()
	gin.SetMode(gin.TestMode)

	// 新しいTestRunのモックデータ
	newTestRun := model.TestRun{
		ProjectID: 1,
		Title:     "New Test Run",
	}
	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `test_runs`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// テスト用のHTTPリクエストとレスポンスライターの作成
	r := gin.Default()
	body, _ := json.Marshal(newTestRun)
	req, _ := http.NewRequest("POST", "/protected/runs", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// テストを実行
	r.POST("/protected/runs", h.PostTestRun)
	r.ServeHTTP(w, req)

	// レスポンスと期待される結果を検証
	assert.Equal(t, http.StatusCreated, w.Code)
	var returnedTestRun model.TestRun
	err := json.Unmarshal(w.Body.Bytes(), &returnedTestRun)
	assert.NoError(t, err)
	assert.Equal(t, newTestRun.Title, returnedTestRun.Title)

	// モックの期待が満たされたか検証
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPostTestRunCase(t *testing.T) {
	h, mock := setupMockTestRunHandler()
	gin.SetMode(gin.TestMode)

	// 新しいTestRunのモックデータ
	newTestRun := model.TestRunCase{
		TestCaseID:   1,
		StatusID:     1,
		Comments:     nil,
		AssignedToID: nil,
		TestRunID:    1,
	}
	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `test_run_cases`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// テスト用のHTTPリクエストとレスポンスライターの作成
	r := gin.Default()
	body, _ := json.Marshal(newTestRun)
	req, _ := http.NewRequest("POST", "/protected/runs/cases", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// テストを実行
	r.POST("/protected/runs/cases", h.PostTestRunCase)
	r.ServeHTTP(w, req)

	// レスポンスと期待される結果を検証
	assert.Equal(t, http.StatusCreated, w.Code)
	var returnedTestRun model.TestRun
	err := json.Unmarshal(w.Body.Bytes(), &returnedTestRun)
	assert.NoError(t, err)

	// モックの期待が満たされたか検証
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPutTestRun(t *testing.T) {
	h, mock := setupMockTestRunHandler()
	gin.SetMode(gin.TestMode)

	// テスト対象のTestRunのID
	testRunID := 1

	// 更新データ
	updatedData := model.TestRun{
		ProjectID: 1,
		Title:     "Updated Test Run",
	}

	// モックの期待を設定
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `test_runs`").
		WithArgs(sqlmock.AnyArg(), updatedData.ProjectID, updatedData.Title, testRunID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// HTTPリクエストの設定
	r := gin.Default()
	body, _ := json.Marshal(updatedData)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/protected/runs/%d", testRunID), bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// テスト実行
	r.PUT("/protected/runs/:id", h.PutTestRun)
	r.ServeHTTP(w, req)

	// レスポンスの検証
	assert.Equal(t, http.StatusOK, w.Code)

	// モックの期待が満たされたか検証
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPutTestRunCase(t *testing.T) {
	h, mock := setupMockTestRunHandler()
	gin.SetMode(gin.TestMode)

	testRunCaseID := 1
	updatedTestRunCase := model.TestRunCase{
		TestCaseID:   2,
		TestRunID:    3,
		Comments:     nil,
		AssignedToID: nil,
		StatusID:     5,
	}

	// モックの期待を設定
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `test_run_cases`").
		WithArgs(sqlmock.AnyArg(), updatedTestRunCase.TestCaseID, updatedTestRunCase.TestRunID, updatedTestRunCase.StatusID, testRunCaseID).
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// HTTPリクエストの設定
	r := gin.Default()
	body, _ := json.Marshal(updatedTestRunCase)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/protected/runs/cases/%d", testRunCaseID), bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// テスト実行
	r.PUT("/protected/runs/cases/:id", h.PutTestRunCase)
	r.ServeHTTP(w, req)

	// レスポンスの検証
	assert.Equal(t, http.StatusOK, w.Code)
	var returnedTestRunCase model.TestRunCase
	err := json.Unmarshal(w.Body.Bytes(), &returnedTestRunCase)
	assert.NoError(t, err)
	assert.Equal(t, updatedTestRunCase.Comments, returnedTestRunCase.Comments)

	// モックの期待が満たされたか検証
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteTestRunCase(t *testing.T) {
	h, mock := setupMockTestRunHandler()
	gin.SetMode(gin.TestMode)

	testRunCaseID := 1

	// モックの期待を設定
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `test_run_cases`").
		WithArgs(sqlmock.AnyArg(), testRunCaseID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// HTTPリクエストの設定
	r := gin.Default()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/protected/runs/cases/%d", testRunCaseID), nil)
	w := httptest.NewRecorder()

	// テスト実行
	r.DELETE("/protected/runs/cases/:id", h.DeleteTestRunCase)
	r.ServeHTTP(w, req)

	// レスポンスの検証
	assert.Equal(t, http.StatusNoContent, w.Code)

	// モックの期待が満たされたか検証
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteTestRun(t *testing.T) {
	h, mock := setupMockTestRunHandler()
	gin.SetMode(gin.TestMode)

	testRunID := 1

	// モックの期待を設定
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `test_runs`").
		WithArgs(sqlmock.AnyArg(), testRunID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	// HTTPリクエストの設定
	r := gin.Default()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/protected/runs/%d", testRunID), nil)
	w := httptest.NewRecorder()

	// テスト実行
	r.DELETE("/protected/runs/:id", h.DeleteTestRun)
	r.ServeHTTP(w, req)

	// レスポンスの検証
	assert.Equal(t, http.StatusNoContent, w.Code)

	// モックの期待が満たされたか検証
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPostTestRunCaseBulk(t *testing.T) {
	h, mock := setupMockTestRunHandler()
	gin.SetMode(gin.TestMode)

	// モックデータの設定
	requestBody := handler.TestRunCasesRequest{
		TestCaseIDs: []uint{1, 2, 3},
		TestRunID:   1,
	}

	// 既存のTestRunCasesのSELECTクエリを模擬
	existingTestRunCasesRows := sqlmock.NewRows([]string{"id", "test_run_id", "test_case_id", "status_id"})
	// 既存のデータがある場合、ここで行を追加してください。例えば:
	// .AddRow(1, 1, 2, 1)  // 既存のTestRunCaseの例
	mock.ExpectQuery("^SELECT \\* FROM `test_run_cases`").
		WithArgs(requestBody.TestRunID).
		WillReturnRows(existingTestRunCasesRows)

	// TestRunCasesのINSERTクエリを模擬
	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `test_run_cases`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 3)) // 3件の挿入を模擬
	mock.ExpectCommit()

	// HTTPリクエストの設定
	r := gin.Default()
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/protected/runs/cases/bulk", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// テスト実行
	r.POST("/protected/runs/cases/bulk", h.PostTestRunCaseBulk)
	r.ServeHTTP(w, req)

	// レスポンスと期待される結果を検証
	assert.Equal(t, http.StatusCreated, w.Code)

	// モックの期待が満たされたか検証
	err := mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestPostTestRunCaseComment(t *testing.T) {
	h, mock := setupMockTestRunHandler()
	gin.SetMode(gin.TestMode)

	// テストケースコメントのリクエストボディを定義
	requestBody := model.Comment{
		// 必要に応じてリクエストボディを設定
		TestRunCaseID: 1,
		Content:       "Test Comment",
	}

	// データベースへのコメント挿入を模擬するモック期待を設定
	mock.ExpectBegin()
	mock.ExpectExec("^INSERT INTO `comments`").
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), requestBody.TestRunCaseID, sqlmock.AnyArg(), requestBody.Content, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1行が挿入されたと仮定
	mock.ExpectCommit()

	commentRows := sqlmock.NewRows([]string{"id", "test_run_case_id", "content"}).
		AddRow(1, requestBody.TestRunCaseID, requestBody.Content)
	mock.ExpectQuery("^SELECT \\* FROM `comments`").
		WithArgs(1).
		WillReturnRows(commentRows)

	// HTTPリクエストとレコーダの設定
	r := gin.Default()
	body, _ := json.Marshal(requestBody)
	req, _ := http.NewRequest("POST", "/protected/runs/cases/comments", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// リクエストの実行
	r.POST("/protected/runs/cases/comments", h.PostTestRunCaseComment)
	r.ServeHTTP(w, req)

	// レスポンスと期待値の検証
	assert.Equal(t, http.StatusCreated, w.Code)
	var returnedComment model.Comment
	err := json.Unmarshal(w.Body.Bytes(), &returnedComment)
	assert.NoError(t, err)
	assert.Equal(t, requestBody.Content, returnedComment.Content)

	// モックの期待が満たされたか確認
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestPutTestRunCaseComment(t *testing.T) {
	h, mock := setupMockTestRunHandler()
	gin.SetMode(gin.TestMode)

	// テスト対象のコメントID
	commentID := 1

	// 更新するコメントデータ
	updatedComment := model.Comment{
		Content:       "Updated Comment",
		TestRunCaseID: 2,
	}

	// データベースへのコメント更新を模擬するモック期待を設定
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `comments`").
		WithArgs(sqlmock.AnyArg(), updatedComment.TestRunCaseID, updatedComment.Content, commentID).
		WillReturnResult(sqlmock.NewResult(1, 1)) // 1行が更新されたと仮定
	mock.ExpectCommit()

	// HTTPリクエストとレコーダの設定
	r := gin.Default()
	body, _ := json.Marshal(updatedComment)
	req, _ := http.NewRequest("PUT", fmt.Sprintf("/protected/runs/cases/comments/%d", commentID), bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	// リクエストの実行
	r.PUT("/protected/runs/cases/comments/:id", h.PutTestRunCaseComment)
	r.ServeHTTP(w, req)

	// レスポンスと期待値の検証
	assert.Equal(t, http.StatusOK, w.Code)
	var returnedComment model.Comment
	err := json.Unmarshal(w.Body.Bytes(), &returnedComment)
	assert.NoError(t, err)
	assert.Equal(t, updatedComment.Content, returnedComment.Content)
	assert.Equal(t, updatedComment.TestRunCaseID, returnedComment.TestRunCaseID)

	// モックの期待が満たされたか確認
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteTestRunCaseComment(t *testing.T) {
	h, mock := setupMockTestRunHandler()
	gin.SetMode(gin.TestMode)

	// テスト対象のコメントID
	commentID := 1

	// データベースからのコメント削除を模擬するモック期待を設定
	mock.ExpectBegin()
	mock.ExpectExec("^UPDATE `comments`").
		WithArgs(sqlmock.AnyArg(), commentID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1行が削除されたと仮定
	mock.ExpectCommit()

	// HTTPリクエストとレコーダの設定
	r := gin.Default()
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/protected/runs/cases/comments/%d", commentID), nil)
	w := httptest.NewRecorder()

	// リクエストの実行
	r.DELETE("/protected/runs/cases/comments/:id", h.DeleteTestRunCaseComment)
	r.ServeHTTP(w, req)

	// レスポンスと期待値の検証
	assert.Equal(t, http.StatusNoContent, w.Code)

	// モックの期待が満たされたか確認
	assert.NoError(t, mock.ExpectationsWereMet())
}
