package dao

import (
	"zhihu-go/internal/model"

	"gorm.io/gorm"
)

//创建用户关注关系

func FollowUser(db *gorm.DB, followerID, followeeID uint) error {
	follow := model.Follow{
		FolloweeID: followeeID,
		FollowerID: followerID,
	}
	return db.Create(&follow).Error
}

//获取用户的所有粉丝

func GetFollowers(db *gorm.DB, userID uint) ([]model.User, error) {
	var followers []model.User
	err := db.Joins("Join follows ON follows.follower_id = users.id").
		Where("follows.followee_id = ?", userID).
		Find(&followers).Error
	return followers, err
}

//获取用户的所有关注

func GetFollowees(db *gorm.DB, userID uint) ([]model.User, error) {
	var followees []model.User
	err := db.Joins("Join follows ON follows.followee_id = users.id").
		Where("follows.follower_id = ?", userID).
		Find(&followees).Error
	return followees, err
}
