package model

import "gorm.io/gorm"

type Permission struct {
	gorm.Model
	Name        string `json:"name"`        // 権限名
	Description string `json:"description"` // 権限の説明
}
