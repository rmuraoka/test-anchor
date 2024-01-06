package router

import (
	"backend/handler"
	"backend/middleware"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
)

func NewRouter(
	db *gorm.DB,
	authHandler *handler.AuthHandler,
	memberHandler *handler.MemberHandler,
	testCaseHandler *handler.TestCaseHandler,
	testRunHandler *handler.TestRunHandler,
	projectHandler *handler.ProjectHandler,
	testPlanHandler *handler.TestPlanHandler,
	milestoneHandler *handler.MilestoneHandler,
) *gin.Engine {
	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_ORIGIN"))
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	r.POST("/login", authHandler.PostLogin)

	protected := r.Group("/protected")
	protected.Use(middleware.AuthenticateJWT(), middleware.CheckUserStatus(db))
	{
		protected.GET("/members", memberHandler.GetMembers)
		protected.GET("/members/:id", memberHandler.GetMember)
		protected.POST("/members", memberHandler.PostMember)
		protected.PUT("/members/:id", memberHandler.PutMember)
		protected.DELETE("/members/:id", memberHandler.DeleteMember)

		protected.GET("/projects", projectHandler.GetProjects)
		protected.GET("/projects/:project_code", projectHandler.GetProject)
		protected.POST("/projects", projectHandler.PostProject)
		protected.PUT("/projects/:id", projectHandler.PutProject)
		protected.DELETE("/projects/:id", projectHandler.DeleteProject)

		protected.GET("/:project_code/cases", testCaseHandler.GetTestCases)
		protected.GET("/cases/:id", testCaseHandler.GetTestCase)
		protected.POST("/cases", testCaseHandler.PostTestCase)
		protected.PUT("/cases/:id", testCaseHandler.PutTestCase)
		protected.DELETE("/cases/:id", testCaseHandler.DeleteTestCase)
		protected.POST("/suites", testCaseHandler.PostTestSuite)
		protected.PUT("/suites/:id", testCaseHandler.PutTestSuite)
		protected.DELETE("/suites/:id", testCaseHandler.DeleteTestSuite)
		protected.PUT("/:project_code/cases/bulk", testCaseHandler.PutTestCaseBulk)

		protected.GET("/:project_code/:test_plan_id/runs", testRunHandler.GetTestRuns)
		protected.GET("/runs/:id", testRunHandler.GetTestRunCases)
		protected.GET("/runs/cases/:id", testRunHandler.GetTestRunCase)
		protected.POST("/runs", testRunHandler.PostTestRun)
		protected.POST("/runs/cases", testRunHandler.PostTestRunCase)
		protected.POST("/runs/cases/bulk", testRunHandler.PostTestRunCaseBulk)
		protected.POST("/runs/cases/comments", testRunHandler.PostTestRunCaseComment)
		protected.PUT("/runs/:id", testRunHandler.PutTestRun)
		protected.PUT("/runs/cases/:id", testRunHandler.PutTestRunCase)
		protected.PUT("/runs/cases/comments/:id", testRunHandler.PutTestRunCaseComment)
		protected.DELETE("/runs/cases/:id", testRunHandler.DeleteTestRunCase)
		protected.DELETE("/runs/:id", testRunHandler.DeleteTestRun)
		protected.DELETE("/runs/cases/comments/:id", testRunHandler.DeleteTestRunCaseComment)

		protected.GET("/:project_code/plans", testPlanHandler.GetTestPlans)
		protected.GET("/plans/:id", testPlanHandler.GetTestPlan)
		protected.POST("/plans", testPlanHandler.PostTestPlan)
		protected.PUT("/plans/:id", testPlanHandler.PutTestPlan)
		protected.DELETE("/plans/:id", testPlanHandler.DeletePlan)

		protected.GET("/:project_code/milestones", milestoneHandler.GetMilestones)
		protected.GET("/milestones/:id", milestoneHandler.GetMilestone)
		protected.POST("/milestones", milestoneHandler.PostMilestone)
		protected.PUT("/milestones/:id", milestoneHandler.PutMilestone)
		protected.DELETE("/milestones/:id", milestoneHandler.DeleteMilestone)
	}

	return r
}
