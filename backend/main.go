package main

import (
	"backend/handler"
	"backend/model"
	"backend/router"
	"backend/util"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	// DSN (Data Source Name) を定義
	dsn := "user:password@tcp(db:3306)/testanchorDB?charset=utf8mb4&parseTime=True&loc=Local"
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
	emailSender := &util.SMTPSender{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		UseTLS:   false,
	}

	// ハンドラーの初期化
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
		authHandler,
		memberHandler,
		testCaseHandler,
		testRunCaseHandler,
		projectHandler,
		testPlanHandler,
		milestoneHandler,
	)
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
		user := model.User{
			Name:     initialUserName,
			Email:    initialUserEmail,
			Password: hashedPassword,
			Status:   "Active",
			Language: "en",
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
}
