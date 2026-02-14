package model

import (
	"time"
)

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Username  string `gorm:"unique; size:50"`
	Email     string `gorm:"unique; default:null; size:100"`
	Bio       string `gorm:"size:200"`
	Password  string
	CreatedAt time.Time
}

type Follow struct {
	ID         uint `gorm:"primaryKey"`
	FolloweeID uint
	FollowerID uint
	CreatedAt  time.Time
}

type Post struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Content   string `gorm:"type:text"`
	AuthorID  uint
	CreatedAt time.Time
}

type Comment struct {
	ID        uint `gorm:"primaryKey"`
	PostID    uint
	AuthorID  uint
	Content   string `gorm:"type:text"`
	CreatedAt time.Time
}

type Like struct {
	ID        uint `gorm:"primaryKey"`
	PostID    uint
	UserID    uint
	CreatedAt time.Time
}
