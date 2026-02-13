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

//修改用户基础信息
