package dao

import (
	"zhihu-go/internal/model"

	"gorm.io/gorm"
)

//创建新用户

func CreateUser(db *gorm.DB, user *model.User) error {
	return db.Create(user).Error
}

//查找用户名是否存在

func GetUserByUsername(db *gorm.DB, username string) (*model.User, error) {
	var user model.User
	err := db.Where("username = ?", username).First(&user).Error
	return &user, err
}

//创建用户关注关系

func FollowUser(db *gorm.DB, followerID, followeeID uint) error {
	follow := model.Follow{
		FolloweeID: followeeID,
		FollowerID: followerID,
	}
	return db.Create(&follow).Error
}

//获取用户的所有关注者

func GetFollowers(db *gorm.DB, userID uint) ([]model.User, error) {
	var followers []model.User
	err := db.Joins("Join follows ON follows.follower_id = user.id").Where("follows.follower_id = ?", userID).Find(&followers).Error
	return followers, err
}
