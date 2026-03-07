package service

import (
	"context"
	"errors"
	"fmt"
	"time"
	"zhihu-go/internal/dao"
	"zhihu-go/internal/dto"
	"zhihu-go/internal/model"
	"zhihu-go/pkg/encrypt"
	"zhihu-go/pkg/jwtutil"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 定义用户相关的业务接口
type UserService interface {
	InitAdmin(ctx context.Context, adminUsername, adminPassword string) error
	RegisterUser(ctx context.Context, username, password string) (*model.User, error)
	LoginUser(ctx context.Context, username, password string) (string, *model.User, error)
	GetUserProfile(ctx context.Context, userID uint) (*model.User, error)
	UpdateProfile(ctx context.Context, userID uint, req *dto.UpdateProfileRequest) (*model.User, error)
	MuteUser(ctx context.Context, targetUserID uint, hours int) error
	CheckMuted(ctx context.Context, userID uint) error
	IsAdmin(ctx context.Context, userID uint) (bool, error)
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
func (s *userService) InitAdmin(ctx context.Context, adminUsername, adminPassword string) error {
	//检查是否有管理员
	count, err := s.userDAO.CountAdmin(ctx)
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
	return s.userDAO.CreateUser(ctx, admin)
}

// RegisterUser 用户注册
func (s *userService) RegisterUser(ctx context.Context, username, password string) (*model.User, error) {
	existing, err := s.userDAO.GetUserByUsername(ctx, username)
	if err == nil && existing != nil {
		return nil, ErrUserAlreadyExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err // 数据库错误（系统错误）
	}

	hashedPassword, err := encrypt.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: hashedPassword,
	}

	if err := s.userDAO.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, err
}

// LoginUser  用户登录
func (s *userService) LoginUser(ctx context.Context, username, password string) (string, *model.User, error) {
	// 1. 查询用户
	user, err := s.userDAO.GetUserByUsername(ctx, username)
	if err != nil {
		// 区分“用户不存在”和其他数据库错误
		// 注意：GORM 返回 gorm.ErrRecordNotFound 表示记录不存在
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", nil, ErrUserNotFound
		}
		return "", nil, err // 其他数据库错误（如连接失败），属于系统错误，将走 default 分支
	}

	// 2. 验证密码
	if !encrypt.CheckPasswordHash(password, user.Password) {
		return "", nil, ErrInvalidPassword
	}

	// 3. 生成 JWT token
	token, err := jwtutil.GenerateToken(user.ID)
	if err != nil {
		return "", nil, err // token 生成失败，系统错误
	}

	return token, user, nil
}

// GetUserProfile 获取用户最新信息
func (s *userService) GetUserProfile(ctx context.Context, userID uint) (*model.User, error) {
	user, err := s.userDAO.GetUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// UpdateProfile 用户信息修改
func (s *userService) UpdateProfile(ctx context.Context, userID uint, req *dto.UpdateProfileRequest) (*model.User, error) {
	updates := make(map[string]interface{})

	//用户名修改
	if req.Username != "" {

		exists, err := s.userDAO.CheckUsernameExists(ctx, req.Username, userID)
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
		exists, err := s.userDAO.CheckEmailExists(ctx, req.Email, userID)
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
		return s.GetUserProfile(ctx, userID)
	}

	//更新信息
	if err := s.userDAO.UpdateProfile(ctx, userID, updates); err != nil {
		return nil, err
	}

	//返回更新后的用户信息
	return s.GetUserProfile(ctx, userID)
}

// MuteUser 禁言或解禁用户
// hours 大于零表示禁言hours小时， 否则则为解除禁言
func (s *userService) MuteUser(ctx context.Context, targetUserID uint, hours int) error {
	// 检查目标用户是否存在
	target, err := s.userDAO.GetUserByID(ctx, targetUserID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("find target user failed: %w", err)
	}

	// 检查目标是否为管理员
	if target.IsAdmin {
		return ErrCannotMuteAdmin
	}

	var mutedUntil *time.Time
	if hours > 0 {
		h := time.Now().Add(time.Duration(hours) * time.Hour)
		mutedUntil = &h
	} else {
		mutedUntil = nil
	}
	if err := s.userDAO.UpdateUserMutedUntil(ctx, targetUserID, mutedUntil); err != nil {
		return ErrMuteFailed
	}

	return nil
}

// CheckMuted 检查用户是否被禁言
func (s *userService) CheckMuted(ctx context.Context, userID uint) error {
	user, err := s.userDAO.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user.MutedUntil != nil && user.MutedUntil.After(time.Now()) {
		return fmt.Errorf("用户已被禁言至 %s", user.MutedUntil.Format("2006-01-02 15:04:05"))
	}

	return nil
}

func (s *userService) IsAdmin(ctx context.Context, userID uint) (bool, error) {
	user, err := s.userDAO.GetUserByID(ctx, userID)
	if err != nil {
		return false, err
	}
	return user.IsAdmin, nil
}
