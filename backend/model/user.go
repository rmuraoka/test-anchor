package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name"`     // ユーザー名
	Email    string `json:"email"`    // メールアドレス
	Password string `json:"password"` // ハッシュ化されたパスワード
	Status   string `json:"status"`
	Roles    []Role `json:"roles" gorm:"many2many:user_roles;"`
	Language string `json:"language"`
}
