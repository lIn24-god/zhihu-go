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

// UpdateProfile 修改用户基础信息
func UpdateProfile(db *gorm.DB, userID uint, updates map[string]interface{}) error {
	return db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}
