package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string `gorm:"unique; size:50"`
	Email      string `gorm:"unique; default:null; size:100"`
	Bio        string `gorm:"size:200; default:null"`
	Password   string
	IsAdmin    bool       `gorm:"default:false"`
	MutedUntil *time.Time `gorm:"default:null"`
}

type Post struct {
	gorm.Model
	Title    string `gorm:"size:100"`
	Content  string `gorm:"type:longtext"`
	AuthorID uint
	Status   string `gorm:"default:draft;index"`
}

type Follow struct {
	ID         uint `gorm:"primaryKey"`
	FolloweeID uint
	FollowerID uint
	CreatedAt  time.Time
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
