package dao

import (
	"zhihu-go/internal/model"

	"gorm.io/gorm"
)

// CreateLike 对文章进行点赞
func CreateLike(db *gorm.DB, like *model.Like) error { return db.Create(like).Error }
