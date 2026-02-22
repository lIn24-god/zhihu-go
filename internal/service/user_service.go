package service

import (
	"errors"
	"fmt"
	"time"
	"zhihu-go/internal/dao"
	"zhihu-go/internal/dto"
	"zhihu-go/internal/model"

	"golang.org/x/crypto/bcrypt"
)

// UserService 定义用户相关的业务接口
type UserService interface {
	InitAdmin(adminUsername, adminPassword string) error
	RegisterUser(username, password string) (*model.User, error)
	LoginUser(username, password string) (*model.User, error)
	GetUserProfile(userID uint) (*model.User, error)
	UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*model.User, error)
	MuteUser(targetUserID uint, hours int) error
	CheckMuted(userID uint) error
	IsAdmin(userID uint) (bool, error)
}

// userService 结构体定义
type userService struct {
	userDAO dao.UserDAO
}

// NewUserService 构造函数
func NewUserService(userDAO dao.UserDAO) UserService {
	return &userService{userDAO: userDAO}
}

// InitAdmin 初始化管理员账号
func (s *userService) InitAdmin(adminUsername, adminPassword string) error {
	//检查是否有管理员
	count, err := s.userDAO.CountAdmin()
	if err != nil {
		return err
	}

	if count > 0 {
		//已有管理员，不用创建
		return nil
	}
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(adminPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &model.User{
		Username: adminUsername,
		Password: string(hashedPassword),
		IsAdmin:  true,
	}
	return s.userDAO.CreateUser(admin)
}

// RegisterUser 用户注册
func (s *userService) RegisterUser(username, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
	}

	if err := s.userDAO.CreateUser(user); err != nil {
		return nil, err
	}

	return user, err
}

// LoginUser 用户登录
func (s *userService) LoginUser(username, password string) (*model.User, error) {
	user, err := s.userDAO.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, err
}

// GetUserProfile 获取用户最新信息
func (s *userService) GetUserProfile(userID uint) (*model.User, error) {
	user, err := s.userDAO.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateProfile 用户信息修改
func (s *userService) UpdateProfile(userID uint, req *dto.UpdateProfileRequest) (*model.User, error) {
	updates := make(map[string]interface{})

	//用户名修改
	if req.Username != "" {

		exists, err := s.userDAO.CheckUsernameExists(req.Username, userID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("username already taken")
		}
		updates["username"] = req.Username
	}

	//用户邮箱修改
	if req.Email != "" {
		exists, err := s.userDAO.CheckEmailExists(req.Email, userID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already taken")
		}
		updates["email"] = req.Email
	}

	//用户简介修改
	if req.Bio != "" {
		updates["bio"] = req.Bio
	}

	//如果没有要修改的信息，则直接返回
	if len(updates) == 0 {
		return s.GetUserProfile(userID)
	}

	//更新信息
	if err := s.userDAO.UpdateProfile(userID, updates); err != nil {
		return nil, err
	}

	//返回更新后的用户信息
	return s.GetUserProfile(userID)
}

// MuteUser 禁言或解禁用户
// hours 大于零表示禁言hours小时， 否则则为解除禁言
func (s *userService) MuteUser(targetUserID uint, hours int) error {
	var mutedUntil *time.Time
	if hours > 0 {
		h := time.Now().Add(time.Duration(hours) * time.Hour)
		mutedUntil = &h
	} else {
		mutedUntil = nil
	}

	return s.userDAO.UpdateUserMutedUntil(targetUserID, mutedUntil)
}

// CheckMuted 检查用户是否被禁言
func (s *userService) CheckMuted(userID uint) error {
	user, err := s.userDAO.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user.MutedUntil != nil && user.MutedUntil.After(time.Now()) {
		return fmt.Errorf("用户已被禁言至 %s", user.MutedUntil.Format("2006-01-02 15:04:05"))
	}

	return nil
}

func (s *userService) IsAdmin(userID uint) (bool, error) {
	user, err := s.userDAO.GetUserByID(userID)
	if err != nil {
		return false, err
	}
	return user.IsAdmin, nil
}
