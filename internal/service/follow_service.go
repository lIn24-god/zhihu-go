package service

import (
	"errors"
	"zhihu-go/internal/dto"
	"zhihu-go/internal/model"

	"gorm.io/gorm"

	"zhihu-go/internal/dao"
)

//关注用户

func FollowUser(db *gorm.DB, followeeID, followerID uint) error {
	//不能关注自己
	if followerID == followeeID {
		return ErrCannotFollowSelf
	}

	//不能关注不存在的人
	var user model.User
	if err := db.First(&user, followeeID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err //其他数据库错误
	}

	//不能关注已关注的人
	var count int64
	if err := db.Model(&model.Follow{}).
		Where("followee_id = ? AND follower_id = ?", followeeID, followerID).
		Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return ErrAlreadyFollowed
	}

	return dao.FollowUser(db, followerID, followeeID)
}

//取关用户

func UnfollowUser(db *gorm.DB, followeeID, followerID uint) error {
	return dao.UnfollowUser(db, followerID, followeeID)
}

//获取用户的所有粉丝

func GetFollowers(db *gorm.DB, userID uint) ([]dto.FollowUserInfo, error) {
	followers, err := dao.GetFollowers(db, userID)
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

//获取用户的所有关注

func GetFollowees(db *gorm.DB, userID uint) ([]dto.FollowUserInfo, error) {
	followees, err := dao.GetFollowees(db, userID)
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
