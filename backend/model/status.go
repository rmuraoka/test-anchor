package model

import "gorm.io/gorm"

type Status struct {
	gorm.Model
	Name  string `json:"name"`  // 名前
	Color string `json:"color"` // 色
}
