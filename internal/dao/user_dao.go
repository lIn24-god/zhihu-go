package dao

import (
	"time"
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

//查找用户名是否存在(靠id)

func GetUserByID(db *gorm.DB, userID uint) (*model.User, error) {
	var user model.User
	err := db.Where("id = ?", userID).First(&user).Error
	return &user, err
}

// UpdateProfile 修改用户基础信息
func UpdateProfile(db *gorm.DB, userID uint, updates map[string]interface{}) error {
	return db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

// UpdateUserMutedUntil 更新用户禁言时间
func UpdateUserMutedUntil(db *gorm.DB, userID uint, mutedUntil *time.Time) error {
	return db.Model(&model.User{}).Where("id = ?", userID).
		Update("muted_until", mutedUntil).Error
}

// CountAdmin 统计管理员数量
func CountAdmin(db *gorm.DB) (int64, error) {
	var count int64
	err := db.Model(&model.User{}).Where("is_admin = ?", true).Count(&count).Error
	return count, err
}
