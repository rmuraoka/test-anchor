package model

import "gorm.io/gorm"

type TestRunCase struct {
	gorm.Model
	TestCaseID   uint      `json:"test_case_id"`
	TestRunID    uint      `json:"test_run_id"`
	AssignedToID *uint     `json:"assigned_to_id"`
	StatusID     uint      `json:"status_id"`
	Comments     []Comment `json:"comments" gorm:"foreignKey:TestRunCaseID;hasMany"`
	TestRun      TestRun   `json:"test_run" gorm:"foreignKey:TestRunID"`
	TestCase     TestCase  `json:"test_case" gorm:"foreignKey:TestCaseID"`
	AssignedTo   *User     `json:"assigned_to" gorm:"foreignKey:AssignedToID"`
	Status       *Status   `json:"status" gorm:"foreignKey:StatusID"`
}
