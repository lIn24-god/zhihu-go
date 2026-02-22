package dao

import (
	"time"
	"zhihu-go/internal/model"

	"gorm.io/gorm"
)

// UserDAO 定义用户数据访问接口
type UserDAO interface {
	CreateUser(user *model.User) error
	GetUserByUsername(username string) (*model.User, error)
	GetUserByID(userID uint) (*model.User, error)
	GetUsersByIDs(userIDs []uint) ([]model.User, error)
	UpdateProfile(userID uint, updates map[string]interface{}) error
	UpdateUserMutedUntil(userID uint, mutedUntil *time.Time) error
	CountAdmin() (int64, error)
	CheckUsernameExists(username string, excludeUserID uint) (bool, error)
	CheckEmailExists(email string, excludeUserID uint) (bool, error)
}

// 结构体定义
type userDAO struct {
	db *gorm.DB
}

// NewUserDAO 构造函数
func NewUserDAO(db *gorm.DB) UserDAO {
	return &userDAO{db: db}
}

// CreateUser 创建新用户
func (u *userDAO) CreateUser(user *model.User) error {
	return u.db.Create(user).Error
}

// GetUserByUsername 查找用户名是否存在
func (u *userDAO) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	err := u.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

// GetUserByID 靠id获取用户
func (u *userDAO) GetUserByID(userID uint) (*model.User, error) {
	var user model.User
	err := u.db.Where("id = ?", userID).First(&user).Error
	return &user, err
}

// GetUsersByIDs 靠一堆id获取一堆用户
func (u *userDAO) GetUsersByIDs(userIDs []uint) ([]model.User, error) {
	var result []model.User
	err := u.db.Where("id IN ?", userIDs).Find(&result).Error
	return result, err
}

// UpdateProfile 修改用户基础信息
func (u *userDAO) UpdateProfile(userID uint, updates map[string]interface{}) error {
	return u.db.Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
}

// UpdateUserMutedUntil 更新用户禁言时间
func (u *userDAO) UpdateUserMutedUntil(userID uint, mutedUntil *time.Time) error {
	return u.db.Model(&model.User{}).Where("id = ?", userID).
		Update("muted_until", mutedUntil).Error
}

// CountAdmin 统计管理员数量
func (u *userDAO) CountAdmin() (int64, error) {
	var count int64
	err := u.db.Model(&model.User{}).Where("is_admin = ?", true).Count(&count).Error
	return count, err
}

// CheckUsernameExists 检查用户名是否已存在
func (u *userDAO) CheckUsernameExists(username string, excludeUserID uint) (bool, error) {
	var count int64
	err := u.db.Model(&model.User{}).
		Where("username = ? AND id != ?", username, excludeUserID).
		Count(&count).Error
	return count > 0, err
}

// CheckEmailExists 检查邮箱是否已存在
func (u *userDAO) CheckEmailExists(email string, excludeUserID uint) (bool, error) {
	var count int64
	err := u.db.Model(&model.User{}).
		Where("email = ? AND id != ?", email, excludeUserID).
		Count(&count).Error
	return count > 0, err
}
