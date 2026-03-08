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
