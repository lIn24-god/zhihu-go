package service

import (
	"zhihu-go/internal/model"
	"zhihu-go/internal/utils"

	"zhihu-go/internal/dto"

	"github.com/redis/go-redis/v9"

	"zhihu-go/internal/dao"
)

// CommentService 定义评论数据访问接口
type CommentService interface {
	CreateComment(req *dto.CommentRequest, authorID uint) (*dto.CommentResponse, error)
}

// 结构体定义
type commentService struct {
	commentDAO dao.CommentDAO
	rdb        *redis.Client
}

// NewCommentService 构造函数
func NewCommentService(commentDAO dao.CommentDAO, rdb *redis.Client) CommentService {
	return &commentService{
		commentDAO: commentDAO,
		rdb:        rdb,
	}
}

// CreateComment 创建评论
func (s *commentService) CreateComment(req *dto.CommentRequest, authorID uint) (*dto.CommentResponse, error) {
	comment := &model.Comment{
		PostID:   req.PostID,
		AuthorID: authorID,
		Content:  req.Content,
	}

	// 防刷机制
	allowed, err := utils.Allow(s.rdb, "评论", authorID, 5, 60)
	if err != nil {
		return nil, err
	}
	if !allowed {
		return nil, ErrTooFrequent
	}

	err1 := s.commentDAO.CreateComment(comment)

	response := &dto.CommentResponse{Content: req.Content}

	return response, err1
}
