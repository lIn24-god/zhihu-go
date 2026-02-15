package dao

import (
	"zhihu-go/internal/model"

	"gorm.io/gorm"
)

// CreateComment 创建评论
func CreateComment(db *gorm.DB, comment *model.Comment) error { return db.Create(comment).Error }
