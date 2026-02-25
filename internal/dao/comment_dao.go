package dao

import (
	"context"
	"zhihu-go/internal/model"

	"gorm.io/gorm"
)

// CommentDAO 定义评论数据访问接口
type CommentDAO interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
}

// 结构体定义
type commentDAO struct {
	db *gorm.DB
}

// NewCommentDAO 构造函数
func NewCommentDAO(db *gorm.DB) CommentDAO { return &commentDAO{db: db} }

// CreateComment 创建评论
func (u *commentDAO) CreateComment(ctx context.Context, comment *model.Comment) error {
	return u.db.WithContext(ctx).Create(comment).Error
}
