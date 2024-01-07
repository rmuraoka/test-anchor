package model

import "gorm.io/gorm"

type Status struct {
	gorm.Model
	Name    string `json:"name"`
	Color   string `json:"color"`
	Default bool   `json:"default"`
}
