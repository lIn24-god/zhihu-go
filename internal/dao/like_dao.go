package dao

import (
	"context"
	"zhihu-go/internal/model"

	"gorm.io/gorm"
)

// LikeDAO 定义点赞数据访问接口
type LikeDAO interface {
	CreateLike(ctx context.Context, like *model.Like) error
}

// 结构体定义
type likeDAO struct {
	db *gorm.DB
}

// NewLikeDAO 构造函数
func NewLikeDAO(db *gorm.DB) LikeDAO { return &likeDAO{db: db} }

// CreateLike 对文章进行点赞
func (u *likeDAO) CreateLike(ctx context.Context, like *model.Like) error {
	return u.db.WithContext(ctx).Create(like).Error
}
