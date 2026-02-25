package service

import (
	"context"
	"errors"
	"fmt"
	"zhihu-go/internal/model"
	"zhihu-go/pkg/ratelimit"

	"zhihu-go/internal/dto"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"zhihu-go/internal/dao"
)

// CommentService 定义评论数据访问接口
type CommentService interface {
	CreateComment(ctx context.Context, req *dto.CommentRequest, authorID uint) (*dto.CommentResponse, error)
}

// 结构体定义
type commentService struct {
	commentDAO  dao.CommentDAO
	userService UserService
	rdb         *redis.Client
}

// NewCommentService 构造函数
func NewCommentService(commentDAO dao.CommentDAO, userService UserService, rdb *redis.Client) CommentService {
	return &commentService{
		commentDAO:  commentDAO,
		userService: userService,
		rdb:         rdb,
	}
}

// CreateComment 创建评论
func (s *commentService) CreateComment(ctx context.Context, req *dto.CommentRequest, authorID uint) (*dto.CommentResponse, error) {
	_, err := s.userService.GetUserProfile(ctx, authorID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("find user failed: %w", err)
	}
	if s.userService.CheckMuted(ctx, authorID) != nil { // 假设 User 模型有 IsMuted 方法或字段
		return nil, ErrUserMuted
	}

	comment := &model.Comment{
		PostID:   req.PostID,
		AuthorID: authorID,
		Content:  req.Content,
	}

	// 防刷机制
	allowed, err := ratelimit.Allow(s.rdb, "评论", authorID, 5, 60)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrTooFrequent
	}

	err1 := s.commentDAO.CreateComment(ctx, comment)

	response := &dto.CommentResponse{Content: req.Content}

	return response, err1
}
