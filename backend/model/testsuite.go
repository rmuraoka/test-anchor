package model

import "gorm.io/gorm"

type TestSuite struct {
	gorm.Model
	ProjectID  uint        `json:"project_id"` // プロジェクトID
	TestCases  []TestCase  `gorm:"foreignKey:TestSuiteID"`
	TestSuites []TestSuite `gorm:"foreignKey:ParentID;references:ID"`
	Name       string      `json:"name"`      // フォルダーの名前
	ParentID   *uint       `json:"parent_id"` // 親フォルダーのID（ルートフォルダーの場合はnil）
}
