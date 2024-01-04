package model

import (
	"gorm.io/gorm"
	"time"
)

type TestRun struct {
	gorm.Model
	ProjectID          uint          `json:"project_id"`
	TestPlanID         uint          `json:"test_plan_id"`
	TestRunCases       []TestRunCase `json:"test_run_cases" gorm:"foreignKey:TestRunID"`
	Title              string        `json:"title"`
	StartedAt          *time.Time    `json:"started_at"`
	CompletedAt        *time.Time    `json:"completed_at"`
	CreatedByID        uint          `json:"created_by_id"`
	UpdatedByID        uint          `json:"updated_by_id"`
	Project            Project       `gorm:"foreignKey:ProjectID"`
	CreatedBy          User          `json:"created_by" gorm:"foreignKey:CreatedByID"`
	UpdatedBy          User          `json:"updated_by" gorm:"foreignKey:UpdatedByID"`
	Status             string        `json:"status"`
	FinalizedTestCases string        `json:"finalized_test_cases"`
}
