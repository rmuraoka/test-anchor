package model

import "gorm.io/gorm"

type Project struct {
	gorm.Model
	TestSuites  []TestSuite `gorm:"foreignKey:ProjectID"`
	TestCases   []TestCase  `gorm:"foreignKey:ProjectID"`
	Milestones  []Milestone `gorm:"foreignKey:ProjectID"`
	TestPlans   []TestPlan  `gorm:"foreignKey:ProjectID"`
	Code        string      `json:"code" gorm:"unique;not null"`
	Title       string      `json:"title"`
	Description string      `json:"description"`
}
