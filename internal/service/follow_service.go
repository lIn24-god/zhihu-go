package service

import (
	"context"
	"errors"
	"zhihu-go/internal/dto"

	"gorm.io/gorm"

	"zhihu-go/internal/dao"
)

// FollowService 定义关注数据访问接口
type FollowService interface {
	FollowUser(ctx context.Context, followeeID, followerID uint) error
	UnfollowUser(ctx context.Context, followeeID, followerID uint) error
	GetFollowers(ctx context.Context, userID uint) ([]dto.FollowUserInfo, error)
	GetFollowees(ctx context.Context, userID uint) ([]dto.FollowUserInfo, error)
}

// 定义结构体
type followService struct {
	followDAO dao.FollowDAO
	userDAO   dao.UserDAO
}

// NewFollowService 构造函数
func NewFollowService(followDAO dao.FollowDAO, userDAO dao.UserDAO) FollowService {
	return &followService{
		followDAO: followDAO,
		userDAO:   userDAO,
	}
}

// FollowUser 关注用户
func (s *followService) FollowUser(ctx context.Context, followeeID, followerID uint) error {
	//不能关注自己
	if followerID == followeeID {
		return ErrCannotFollowSelf
	}

	//不能关注不存在的人
	_, err := s.userDAO.GetUserByID(ctx, followeeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err //其他数据库错误
	}

	//不能关注已关注的人
	exists, err := s.followDAO.CheckFollowExists(ctx, followeeID, followerID)
	if exists {
		return ErrAlreadyFollowed
	}

	return s.followDAO.FollowUser(ctx, followerID, followeeID)
}

// UnfollowUser 取关用户
func (s *followService) UnfollowUser(ctx context.Context, followeeID, followerID uint) error {
	return s.followDAO.UnfollowUser(ctx, followerID, followeeID)
}

// GetFollowers 获取用户的所有粉丝
func (s *followService) GetFollowers(ctx context.Context, userID uint) ([]dto.FollowUserInfo, error) {
	followers, err := s.followDAO.GetFollowers(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []dto.FollowUserInfo
	for _, f := range followers {
		result = append(result, dto.FollowUserInfo{
			Username: f.Username,
			ID:       f.ID,
		})
	}

	return result, err
}

// GetFollowees 获取用户的所有关注
func (s *followService) GetFollowees(ctx context.Context, userID uint) ([]dto.FollowUserInfo, error) {
	followees, err := s.followDAO.GetFollowees(ctx, userID)
	if err != nil {
		return nil, err
	}

	var result []dto.FollowUserInfo
	for _, f := range followees {
		result = append(result, dto.FollowUserInfo{
			Username: f.Username,
			ID:       f.ID,
		})
	}

	return result, err
}
