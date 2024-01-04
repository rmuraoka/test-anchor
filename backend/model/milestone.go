package model

import (
	"gorm.io/gorm"
	"time"
)

type Milestone struct {
	gorm.Model
	ProjectID   uint      `json:"project_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	DueDate     time.Time `json:"due_date" sql:"not null;type:date"`
	Project     Project   `gorm:"foreignKey:ProjectID"`
}
