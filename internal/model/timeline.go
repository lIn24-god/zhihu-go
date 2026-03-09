package model

import "gorm.io/gorm"

type Timeline struct {
	gorm.Model
	UserID   uint `gorm:"index:idx_user_created,priority:1;not null;comment:收件人用户ID"`
	PostID   uint `gorm:"not null;comment:动态ID（文章ID）"`
	AuthorID uint `gorm:"not null;comment:动态作者ID"`
	IsOwn    bool `gorm:"default:false;comment:是否自己的动态"`
}
