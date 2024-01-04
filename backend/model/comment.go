package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	TestRunCaseID uint        `json:"test_run_case_id"`
	StatusID      *uint       `json:"status_id"`
	Content       string      `json:"content"`
	CreatedByID   uint        `json:"created_by_id"`
	UpdatedByID   uint        `json:"updated_by_id"`
	TestRunCase   TestRunCase `gorm:"foreignKey:TestRunCaseID"`
	Status        *Status     `json:"status" gorm:"foreignKey:StatusID"`
	CreatedBy     User        `json:"created_by" gorm:"foreignKey:CreatedByID"`
	UpdatedBy     User        `json:"updated_by" gorm:"foreignKey:UpdatedByID"`
}
