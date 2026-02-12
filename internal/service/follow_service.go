package service

import (
	"zhihu-go/internal/dto"

	"gorm.io/gorm"

	"zhihu-go/internal/dao"
)

//关注用户

func FollowUser(db *gorm.DB, followeeID, followerID uint) error {
	return dao.FollowUser(db, followerID, followeeID)
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
