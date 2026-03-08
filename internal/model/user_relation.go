package model

import (
	"gorm.io/gorm"
)

type UserRelation struct {
	gorm.Model
	UserID       uint   `gorm:"index;not null;comment:主动方用户ID"`
	TargetID     uint   `gorm:"index;not null;comment:被动方用户ID"`
	RelationType string `gorm:"type:varchar(20);default:'follow';comment:关系类型"` //类型有关注，拉黑等
}
