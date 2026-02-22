package dao

import (
	"zhihu-go/internal/model"

	"gorm.io/gorm"
)

// FollowDAO 定义关注数据访问接口
type FollowDAO interface {
	FollowUser(followerID, followeeID uint) error
	UnfollowUser(followerID, followeeID uint) error
	GetFollowers(userID uint) ([]model.User, error)
	GetFollowees(userID uint) ([]model.User, error)
	CheckFollowExists(followeeID, followerID uint) (bool, error)
}

// 结构体定义
type followDAO struct {
	db *gorm.DB
}

// NewFollowDAO 构造函数
func NewFollowDAO(db *gorm.DB) FollowDAO { return &followDAO{db: db} }

// FollowUser 创建用户关注关系
func (u *followDAO) FollowUser(followerID, followeeID uint) error {
	follow := model.Follow{
		FolloweeID: followeeID,
		FollowerID: followerID,
	}
	return u.db.Create(&follow).Error
}

// UnfollowUser 取消关注用户
func (u *followDAO) UnfollowUser(followerID, followeeID uint) error {
	return u.db.Where("follower_id = ? AND followee_id = ?", followerID, followeeID).
		Delete(&model.Follow{}).Error
}

// GetFollowers 获取用户的所有粉丝
func (u *followDAO) GetFollowers(userID uint) ([]model.User, error) {
	var followers []model.User
	err := u.db.Joins("Join follows ON follows.follower_id = users.id").
		Where("follows.followee_id = ?", userID).
		Find(&followers).Error
	return followers, err
}

// GetFollowees 获取用户的所有关注
func (u *followDAO) GetFollowees(userID uint) ([]model.User, error) {
	var followees []model.User
	err := u.db.Joins("Join follows ON follows.followee_id = users.id").
		Where("follows.follower_id = ?", userID).
		Find(&followees).Error
	return followees, err
}

func (u *followDAO) CheckFollowExists(followeeID, followerID uint) (bool, error) {
	var count int64
	err := u.db.Model(&model.Follow{}).
		Where("followee_id = ? AND follower_id = ?", followeeID, followerID).
		Count(&count).Error

	return count > 0, err
}
