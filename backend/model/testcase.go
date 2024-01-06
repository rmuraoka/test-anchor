package model

import "gorm.io/gorm"

type TestCase struct {
	gorm.Model
	TestSuiteID *uint     `json:"test_suite_id"`
	ProjectID   uint      `json:"project_id"`
	MilestoneID *uint     `json:"milestone_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	OrderIndex  int       `json:"order_index"`
	TestSuite   TestSuite `gorm:"foreignKey:TestSuiteID;pointer"`
	Project     Project   `gorm:"foreignKey:ProjectID"`
	Milestone   Milestone `gorm:"foreignKey:MilestoneID"`

	CreatedByID uint `json:"created_by_id"`
	UpdatedByID uint `json:"updated_by_id"`

	CreatedBy User `gorm:"foreignKey:CreatedByID"`
	UpdatedBy User `gorm:"foreignKey:UpdatedByID"`
}
