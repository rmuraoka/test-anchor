package model

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	Name        string `json:"name"`        // ロール名
	Description string `json:"description"` // ロールの説明
}
