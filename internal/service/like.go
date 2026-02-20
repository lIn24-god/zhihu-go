package service

import (
	"zhihu-go/internal/model"

	"gorm.io/gorm"

	"zhihu-go/internal/dto"

	"zhihu-go/internal/dao"
)

// CreateLike 点赞文章
func CreateLike(db *gorm.DB, req dto.LikeRequest, userID uint) error {
	like := model.Like{
		PostID: req.PostID,
		UserID: userID,
	}

	return dao.CreateLike(db, &like)
}
