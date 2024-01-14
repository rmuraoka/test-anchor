package model

import "gorm.io/gorm"

type RolePermission struct {
	gorm.Model
	RoleID       uint       `json:"role_id"`
	PermissionID uint       `json:"permission_id"`
	Role         Role       `gorm:"foreignKey:RoleID"`
	Permission   Permission `gorm:"foreignKey:PermissionID"`
}
