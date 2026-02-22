package service

import (
	"errors"
	"zhihu-go/internal/dto"

	"gorm.io/gorm"

	"zhihu-go/internal/dao"
)

// FollowService 定义关注数据访问接口
type FollowService interface {
	FollowUser(followeeID, followerID uint) error
	UnfollowUser(followeeID, followerID uint) error
	GetFollowers(userID uint) ([]dto.FollowUserInfo, error)
	GetFollowees(userID uint) ([]dto.FollowUserInfo, error)
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
func (s *followService) FollowUser(followeeID, followerID uint) error {
	//不能关注自己
	if followerID == followeeID {
		return ErrCannotFollowSelf
	}

	//不能关注不存在的人
	_, err := s.userDAO.GetUserByID(followeeID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err //其他数据库错误
	}

	//不能关注已关注的人
	exists, err := s.followDAO.CheckFollowExists(followeeID, followerID)
	if exists {
		return ErrAlreadyFollowed
	}

	return s.followDAO.FollowUser(followerID, followeeID)
}

// UnfollowUser 取关用户
func (s *followService) UnfollowUser(followeeID, followerID uint) error {
	return s.followDAO.UnfollowUser(followerID, followeeID)
}

// GetFollowers 获取用户的所有粉丝
func (s *followService) GetFollowers(userID uint) ([]dto.FollowUserInfo, error) {
	followers, err := s.followDAO.GetFollowers(userID)
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
func (s *followService) GetFollowees(userID uint) ([]dto.FollowUserInfo, error) {
	followees, err := s.followDAO.GetFollowees(userID)
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
