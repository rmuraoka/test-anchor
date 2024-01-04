package model

import (
	"gorm.io/gorm"
	"time"
)

type TestPlan struct {
	gorm.Model
	ProjectID   uint       `json:"project_id"`   // プロジェクトID
	Status      string     `json:"status"`       // テストの実行状況（例：'未実行', '成功', '失敗'）
	Title       string     `json:"title"`        // テストスイートのタイトル
	StartedAt   *time.Time `json:"started_at"`   // 開始日時
	CompletedAt *time.Time `json:"completed_at"` // 終了日時
	TestRuns    []TestRun  `json:"test_runs" gorm:"foreignKey:TestPlanID"`
	CreatedByID uint       `json:"created_by_id"` // 作成者のユーザーID
	UpdatedByID uint       `json:"updated_by_id"` // 更新者のユーザーID
	CreatedBy   User       `json:"created_by" gorm:"foreignKey:CreatedByID"`
	UpdatedBy   User       `json:"updated_by" gorm:"foreignKey:UpdatedByID"`
}
