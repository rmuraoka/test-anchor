package util

type Auth struct {
	Token string    `json:"token"`
	User  LoginUser `json:"user"`
}

type Project struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Code        string `json:"code"`
	Description string `json:"description"`
}

type TestCaseMilestone struct {
	ID    uint   `json:"id"`
	Title string `json:"title"`
}

type Milestone struct {
	ID            uint   `json:"id"`
	Title         string `json:"title"`
	Description   string `json:"description"`
	DueDate       string `json:"due_date"`
	Status        string `json:"status"`
	TestCaseCount int    `json:"test_case_count"`
}

type Chart struct {
	Name  string `json:"name"`
	Color string `json:"color"`
	Count int    `json:"count"`
}

type ProjectsResponseData struct {
	Projects []Project `json:"entities"`
}

type TestCase struct {
	ID        uint               `json:"id"`
	Title     string             `json:"title"`
	Content   string             `json:"content"`
	Milestone *TestCaseMilestone `json:"milestone"`
	CreatedBy User               `json:"created_by"`
	UpdatedBy User               `json:"updated_by"`
}

type JSONTestSuite struct {
	ID         uint            `json:"id"`
	Name       string          `json:"name"`
	TestSuites []JSONTestSuite `json:"test_suites"`
	TestCases  []TestCase      `json:"test_cases"`
}

type JSONOnlyTestSuite struct {
	ID       uint                `json:"key"`
	Title    string              `json:"title"`
	Children []JSONOnlyTestSuite `json:"children"`
}

type TestRunCasesTestSuite struct {
	Name       string                  `json:"name"`
	TestSuites []TestRunCasesTestSuite `json:"test_suites"`
	TestCases  []TestRunCase           `json:"test_cases"`
}

type TestCasesResponseData struct {
	ProjectID      uint                `json:"project_id"`
	TestSuites     []JSONTestSuite     `json:"entities"`
	OnlyTestSuites []JSONOnlyTestSuite `json:"folders"`
}

type TestRunCasesResponseData struct {
	ProjectID      uint                    `json:"project_id"`
	TestRunID      uint                    `json:"test_run_id"`
	TestPlanId     uint                    `json:"test_plan_id"`
	Status         string                  `json:"status"`
	TestSuites     []TestRunCasesTestSuite `json:"entities"`
	OnlyTestSuites []JSONOnlyTestSuite     `json:"folders"`
}

type TestRunCase struct {
	ID         uint      `json:"id"`
	TestCaseId uint      `json:"test_case_id"`
	Title      string    `json:"title"`
	Content    string    `json:"content"`
	Status     Status    `json:"status"`
	AssignedTo *User     `json:"assigned_to"`
	Comments   []Comment `json:"comments"`
}

type Status struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Comment struct {
	ID        uint    `json:"id"`
	Status    *Status `json:"status"`
	Content   string  `json:"content"`
	CreatedBy User    `json:"created_by"`
	UpdatedBy User    `json:"updated_by"`
}

type User struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type LoginUser struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	Language string `json:"language"`
}

type TestRun struct {
	ID          uint    `json:"id"`
	ProjectID   uint    `json:"project_id"`
	Title       string  `json:"title"`
	Count       int     `json:"count"`
	Status      string  `json:"status"`
	StartedAt   *string `json:"started_at"`
	CompletedAt *string `json:"completed_at"`
	TestCaseIDs []uint  `json:"test_case_ids"`
	CreatedBy   User    `json:"created_by"`
	UpdatedBy   User    `json:"updated_by"`
}

type TestPlan struct {
	ID          uint    `json:"id"`
	ProjectID   uint    `json:"project_id"`
	Title       string  `json:"title"`
	Status      string  `json:"status"`
	StartedAt   *string `json:"started_at"`
	CompletedAt *string `json:"completed_at"`
	CreatedBy   User    `json:"created_by"`
	UpdatedBy   User    `json:"updated_by"`
}

type TestPlanDetail struct {
	ID        uint      `json:"id"`
	ProjectID uint      `json:"project_id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	Charts    []Chart   `json:"charts"`
	TestRuns  []TestRun `json:"test_runs"`
	CreatedBy User      `json:"created_by"`
	UpdatedBy User      `json:"updated_by"`
}

type TestRunsResponseData struct {
	TestRuns []TestRun `json:"entities"`
}

type TestRunCaseResponseData struct {
	TestRunCase TestRunCase `json:"entities"`
}

type TestPlansResponseData struct {
	ProjectId uint       `json:"project_id"`
	TestPlans []TestPlan `json:"entities"`
}

type MilestonesResponseData struct {
	ProjectId  uint        `json:"project_id"`
	Milestones []Milestone `json:"entities"`
}

type ProjectResponseData struct {
	ID          uint        `json:"id"`
	Code        string      `json:"code"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
	Milestones  []Milestone `json:"milestones"`
	TestPlans   []TestPlan  `json:"test_plans"`
}

type Member struct {
	ID     uint   `json:"id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type MembersResponseData struct {
	Members []Member `json:"entities"`
}

type StatusesResponseData struct {
	DefaultID uint     `json:"default_id"`
	Statuses  []Status `json:"entities"`
}
