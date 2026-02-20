package service

import (
	"zhihu-go/internal/model"
	"zhihu-go/internal/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"zhihu-go/internal/dto"

	"zhihu-go/internal/dao"
)

// CreateComment 创建评论
func CreateComment(db *gorm.DB, rdb *redis.Client, req *dto.CommentRequest, authorID uint) (*dto.CommentResponse, error) {
	comment := &model.Comment{
		PostID:   req.PostID,
		AuthorID: authorID,
		Content:  req.Content,
	}

	// 防刷机制
	allowed, err := utils.Allow(rdb, "评论", authorID, 5, 60)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrTooFrequent
	}

	err1 := dao.CreateComment(db, comment)

	response := &dto.CommentResponse{Content: req.Content}

	return response, err1
}
