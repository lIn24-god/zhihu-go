package service

import (
	"context"
	"zhihu-go/internal/model"
	"zhihu-go/pkg/ratelimit"

	"zhihu-go/internal/dto"

	"github.com/redis/go-redis/v9"

	"zhihu-go/internal/dao"
)

// LikeService 定义点赞相关的数据访问接口
type LikeService interface {
	CreateLike(ctx context.Context, req dto.LikeRequest, userID uint) error
}

// 结构体定义
type likeService struct {
	likeDAO dao.LikeDAO
	rdb     *redis.Client
}

// NewLikeService 构造函数
func NewLikeService(likeDAO dao.LikeDAO, rdb *redis.Client) LikeService {
	return &likeService{
		likeDAO: likeDAO,
		rdb:     rdb,
	}
}

// CreateLike 点赞文章
func (s *likeService) CreateLike(ctx context.Context, req dto.LikeRequest, userID uint) error {
	like := model.Like{
		PostID: req.PostID,
		UserID: userID,
	}

	// 防刷机制
	allowed, err := ratelimit.Allow(s.rdb, "点赞", userID, 5, 60)
	if err != nil {
		return err
	}
	if !allowed {
		return ErrTooFrequent
	}

	return s.likeDAO.CreateLike(ctx, &like)
}
