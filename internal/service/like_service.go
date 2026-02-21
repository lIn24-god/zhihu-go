package service

import (
	"zhihu-go/internal/model"
	"zhihu-go/internal/utils"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"zhihu-go/internal/dto"

	"zhihu-go/internal/dao"
)

// CreateLike 点赞文章
func CreateLike(db *gorm.DB, rdb *redis.Client, req dto.LikeRequest, userID uint) error {
	like := model.Like{
		PostID: req.PostID,
		UserID: userID,
	}

	// 防刷机制
	allowed, err := utils.Allow(rdb, "点赞", userID, 5, 60)
	if err != nil {
		return err
	}
	if !allowed {
		return ErrTooFrequent
	}

	return dao.CreateLike(db, &like)
}
