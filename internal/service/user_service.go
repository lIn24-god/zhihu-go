package service

import (
	"errors"
	"zhihu-go/internal/dao"
	"zhihu-go/internal/dto"
	"zhihu-go/internal/model"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

//用户注册

func RegisterUser(db *gorm.DB, username, password string) (*model.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
	}

	if err := dao.CreateUser(db, user); err != nil {
		return nil, err
	}

	return user, err
}

//用户登录

func LoginUser(db *gorm.DB, username, password string) (*model.User, error) {
	user, err := dao.GetUserByUsername(db, username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, err
	}

	return user, err
}

//获取用户最新信息

func GetUserProfile(db *gorm.DB, userID uint) (*dto.UpdateProfileResponse, error) {
	var user model.User
	if err := db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &dto.UpdateProfileResponse{
		ID:       userID,
		Username: user.Username,
		Email:    user.Email,
		Bio:      user.Bio,
	}, nil
}

//用户信息修改

func UpdateProfile(db *gorm.DB, userID uint, rep *dto.UpdateProfileRequest) (*dto.UpdateProfileResponse, error) {
	updates := make(map[string]interface{})

	//用户名修改
	if rep.Username != "" {
		var exists bool
		if err := db.Model(&model.User{}).
			Where("username = ? AND id != ?", rep.Username, userID).
			Select("count(*) > 0").
			Find(&exists).Error; err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("username already taken")
		}
		updates["username"] = rep.Username
	}

	//用户邮箱修改
	if rep.Email != "" {
		var exists bool
		if err := db.Model(&model.User{}).
			Where("email = ? AND id != ?", rep.Email, userID).
			Select("count(*) > 0").
			Find(&exists).Error; err != nil {
			return nil, err
		}
		if exists {
			return nil, errors.New("email already taken")
		}
		updates["email"] = rep.Email
	}

	//用户简介修改
	if rep.Bio != "" {
		updates["bio"] = rep.Bio
	}

	//如果没有要修改的信息，则直接返回
	if len(updates) == 0 {
		return GetUserProfile(db, userID)
	}

	//更新信息
	if err := dao.UpdateProfile(db, userID, updates); err != nil {
		return nil, err
	}

	//返回更新后的用户信息
	return GetUserProfile(db, userID)
}
