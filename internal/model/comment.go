package model

import "gorm.io/gorm"

type Comment struct {
	gorm.Model
	PostID   uint
	AuthorID uint
	Content  string `gorm:"type:text"`
}
