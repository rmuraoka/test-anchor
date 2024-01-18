package router

import (
	"backend/handler"
	"backend/middleware"
	"backend/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"os"
)

func NewRouter(
	db *gorm.DB,
	healthHandler *handler.HealthHandler,
	statusHandler *handler.StatusHandler,
	roleHandler *handler.RoleHandler,
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

	api := r.Group("/api")
	api.GET("/health", healthHandler.GetHealth)
	api.POST("/login", authHandler.PostLogin)

	protected := api.Group("/protected")
	protected.Use(middleware.AuthenticateJWT(), middleware.CheckUserStatus(db))
	{
		protected.GET("/statuses", statusHandler.GetStatues)
		protected.GET("/roles", checkPermission("admin", db), roleHandler.GetRoles)

		protected.GET("/members", memberHandler.GetMembers)
		protected.GET("/admin/members", checkPermission("admin", db), memberHandler.GetAdminMembers)
		protected.GET("/members/:id", checkPermission("admin", db), memberHandler.GetMember)
		protected.POST("/members", checkPermission("admin", db), memberHandler.PostMember)
		protected.PUT("/members/:id", memberHandler.PutMember)
		protected.DELETE("/members/:id", checkPermission("admin", db), memberHandler.DeleteMember)

		protected.GET("/projects", projectHandler.GetProjects)
		protected.GET("/projects/:project_code", projectHandler.GetProject)
		protected.POST("/projects", checkPermission("edit", db), projectHandler.PostProject)
		protected.PUT("/projects/:id", checkPermission("edit", db), projectHandler.PutProject)
		protected.DELETE("/projects/:id", checkPermission("edit", db), projectHandler.DeleteProject)

		protected.GET("/:project_code/cases", testCaseHandler.GetTestCases)
		protected.GET("/cases/:id", testCaseHandler.GetTestCase)
		protected.POST("/cases", checkPermission("edit", db), testCaseHandler.PostTestCase)
		protected.PUT("/cases/:id", checkPermission("edit", db), testCaseHandler.PutTestCase)
		protected.DELETE("/cases/:id", checkPermission("edit", db), testCaseHandler.DeleteTestCase)
		protected.POST("/suites", checkPermission("edit", db), testCaseHandler.PostTestSuite)
		protected.PUT("/suites/:id", checkPermission("edit", db), testCaseHandler.PutTestSuite)
		protected.DELETE("/suites/:id", checkPermission("edit", db), testCaseHandler.DeleteTestSuite)
		protected.PUT("/:project_code/cases/bulk", checkPermission("edit", db), testCaseHandler.PutTestCaseBulk)
		protected.PUT("/:project_code/suites/bulk", checkPermission("edit", db), testCaseHandler.PutTestSuiteBulk)

		protected.GET("/:project_code/:test_plan_id/runs", testRunHandler.GetTestRuns)
		protected.GET("/runs/:id", testRunHandler.GetTestRunCases)
		protected.GET("/runs/cases/:id", testRunHandler.GetTestRunCase)
		protected.POST("/runs", checkPermission("edit", db), testRunHandler.PostTestRun)
		protected.POST("/runs/cases", checkPermission("edit", db), testRunHandler.PostTestRunCase)
		protected.POST("/runs/cases/bulk", checkPermission("edit", db), testRunHandler.PostTestRunCaseBulk)
		protected.POST("/runs/cases/comments", testRunHandler.PostTestRunCaseComment)
		protected.PUT("/runs/:id", testRunHandler.PutTestRun)
		protected.PUT("/runs/cases/:id", testRunHandler.PutTestRunCase)
		protected.PUT("/runs/cases/comments/:id", testRunHandler.PutTestRunCaseComment)
		protected.DELETE("/runs/cases/:id", testRunHandler.DeleteTestRunCase)
		protected.DELETE("/runs/:id", checkPermission("edit", db), testRunHandler.DeleteTestRun)
		protected.DELETE("/runs/cases/comments/:id", checkPermission("edit", db), testRunHandler.DeleteTestRunCaseComment)

		protected.GET("/:project_code/plans", testPlanHandler.GetTestPlans)
		protected.GET("/plans/:id", testPlanHandler.GetTestPlan)
		protected.POST("/plans", checkPermission("edit", db), testPlanHandler.PostTestPlan)
		protected.PUT("/plans/:id", testPlanHandler.PutTestPlan)
		protected.DELETE("/plans/:id", checkPermission("edit", db), testPlanHandler.DeletePlan)

		protected.GET("/:project_code/milestones", milestoneHandler.GetMilestones)
		protected.GET("/milestones/:id", milestoneHandler.GetMilestone)
		protected.POST("/milestones", checkPermission("edit", db), milestoneHandler.PostMilestone)
		protected.PUT("/milestones/:id", checkPermission("edit", db), milestoneHandler.PutMilestone)
		protected.DELETE("/milestones/:id", checkPermission("edit", db), milestoneHandler.DeleteMilestone)
	}

	return r
}

func checkPermission(permission string, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.GetString("email")
		var user model.User
		if err := db.Preload("Role").Preload("Role.RolePermissions").Preload("Role.RolePermissions.Permission").Where("email = ?", email).First(&user).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "内部サーバーエラー"})
			return
		}

		hasPermission := false
		for _, rolePermission := range user.Role.RolePermissions {
			if rolePermission.Permission.Name == permission {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "Forbidden - Insufficient permissions"})
			return
		}

		c.Next()
	}
}
