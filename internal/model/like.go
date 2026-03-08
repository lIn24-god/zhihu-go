package model

import "gorm.io/gorm"

type Like struct {
	gorm.Model
	PostID uint
	UserID uint
}
