package main

import (
	"backend/handler"
	"backend/model"
	"backend/router"
	"backend/util"
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	dbUsername := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUsername, dbPassword, dbHost, dbName)
	var db *gorm.DB
	var err error
	// 最大試行回数
	maxAttempts := 10

	for attempts := 1; attempts <= maxAttempts; attempts++ {
		db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}

		fmt.Printf("データベース接続に失敗しました。再試行します... (%d/%d)\n", attempts, maxAttempts)
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Fatalf("データベース接続に失敗しました: %v", err)
	}
	db.AutoMigrate(
		&model.Status{},
		&model.Permission{},
		&model.Milestone{},
		&model.Role{},
		&model.RolePermission{},
		&model.User{},
		&model.Project{},
		&model.TestSuite{},
		&model.TestCase{},
		&model.TestPlan{},
		&model.TestRun{},
		&model.TestRunCase{},
		&model.Comment{})

	host := os.Getenv("MAIL_HOST")
	port, _ := strconv.Atoi(os.Getenv("MAIL_PORT"))
	username := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	fromEmail := os.Getenv("FROM_EMAIL")
	useTLS, _ := strconv.ParseBool(os.Getenv("USE_TLS"))
	emailSender := &util.SMTPSender{
		Host:      host,
		Port:      port,
		Username:  username,
		Password:  password,
		FromEmail: fromEmail,
		UseTLS:    useTLS,
	}

	// ハンドラーの初期化
	healthHandler := handler.NewHealthHandler()
	statusHandler := handler.NewStatusHandler(db)
	rolesHandler := handler.NewRoleHandler(db)
	authHandler := handler.NewAuthHandler(db)
	memberHandler := handler.NewMemberHandler(db, emailSender)
	testCaseHandler := handler.NewTestCaseHandler(db)
	testRunCaseHandler := handler.NewTestRunHandler(db)
	projectHandler := handler.NewProjectHandler(db)
	testPlanHandler := handler.NewTestPlanHandler(db)
	milestoneHandler := handler.NewMilestoneHandler(db)

	// ルータの初期化
	r := router.NewRouter(
		db,
		healthHandler,
		statusHandler,
		rolesHandler,
		authHandler,
		memberHandler,
		testCaseHandler,
		testRunCaseHandler,
		projectHandler,
		testPlanHandler,
		milestoneHandler,
	)

	createInitialData(db)
	createInitialUser(db, emailSender)

	// サーバを起動
	r.Run(":8000")
}

func createInitialUser(db *gorm.DB, sender util.EmailSender) {
	initialUserEmail := os.Getenv("INITIAL_USER_EMAIL")
	initialUserName := os.Getenv("INITIAL_USER_NAME")
	tempPassword := util.GenerateTempPassword(10)
	hashedPassword, err := util.HashPassword(tempPassword)
	if err != nil {
		log.Fatalf("パスワードのハッシュ化に失敗しました: %v", err)
	}
	var count int64
	db.Model(&model.User{}).Count(&count)
	if count == 0 {
		var role model.Role
		if err := db.Where("name = ?", "Administrator").First(&role).Error; err != nil {
			log.Fatalf("初期ユーザーの作成に失敗しました: %v", err)
		}

		user := model.User{
			Name:     initialUserName,
			Email:    initialUserEmail,
			Password: hashedPassword,
			Status:   "Active",
			Language: "en",
			RoleID:   role.ID,
		}
		if err := db.Create(&user).Error; err != nil {
			log.Fatalf("初期ユーザーの作成に失敗しました: %v", err)
		}

		subject := "Your Account"
		body := "Welcome " + user.Name + " Your Password is " + tempPassword
		if err := sender.SendMail([]string{user.Email}, subject, body); err != nil {
			log.Fatalf("初期ユーザーの招待に失敗しました: %v", err)
		}
	}

	db.Model(&model.User{}).Where("role_id IS NOT NULL").Count(&count)
	if count == 0 {
		var youngestUser model.User
		if err := db.Order("created_at ASC").First(&youngestUser).Error; err != nil {
			log.Fatalf("最初に登録したユーザーの取得に失敗しました: %v", err)
		}
		var role model.Role
		if err := db.Where("name = ?", "Administrator").First(&role).Error; err != nil {
			log.Fatalf("ロールの取得に失敗しました: %v", err)
		}
		if err := db.Model(&youngestUser).Update("role_id", role.ID).Error; err != nil {
			log.Fatalf("管理者権限の付与に失敗しました: %v", err)
		}
	}
}

func createInitialData(db *gorm.DB) {
	var statusCount int64
	db.Model(&model.Status{}).Count(&statusCount)
	if statusCount == 0 {
		var statuses []model.Status
		absPath, _ := filepath.Abs("config/initial_statuses.json")
		byteValue, err := os.ReadFile(absPath)
		if err != nil {
			log.Fatalf("Error reading statuses file: %v", err)
		}
		json.Unmarshal(byteValue, &statuses)

		for _, status := range statuses {
			db.Create(&status)
		}
	}

	var roleCount int64
	db.Model(&model.Role{}).Count(&roleCount)
	if roleCount == 0 {
		var roles []model.Role
		absPath, _ := filepath.Abs("config/initial_roles.json")
		byteValue, err := os.ReadFile(absPath)
		if err != nil {
			log.Fatalf("Error reading statuses file: %v", err)
		}
		json.Unmarshal(byteValue, &roles)
		for _, role := range roles {
			db.Create(&role)
		}
	}

	var permissionCount int64
	db.Model(&model.Permission{}).Count(&permissionCount)
	if permissionCount == 0 {
		var permissions []model.Permission
		absPath, _ := filepath.Abs("config/initial_permissions.json")
		byteValue, err := os.ReadFile(absPath)
		if err != nil {
			log.Fatalf("Error reading permissions file: %v", err)
		}
		json.Unmarshal(byteValue, &permissions)
		for _, permission := range permissions {
			db.Create(&permission)
		}
	}

	var rolePermissionCount int64
	db.Model(&model.RolePermission{}).Count(&rolePermissionCount)
	if rolePermissionCount == 0 {
		var rolePermissions []model.RolePermission
		absPath, _ := filepath.Abs("config/initial_role_permissions.json")
		byteValue, err := os.ReadFile(absPath)
		if err != nil {
			log.Fatalf("Error reading rolePermissions file: %v", err)
		}
		json.Unmarshal(byteValue, &rolePermissions)
		for _, rolePermission := range rolePermissions {
			db.Create(&rolePermission)
		}
	}
}
