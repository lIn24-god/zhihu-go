package model

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title    string `gorm:"size:100"`
	Content  string `gorm:"type:longtext"`
	AuthorID uint
	Status   string `gorm:"default:draft;index"`
}
