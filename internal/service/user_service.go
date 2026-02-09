package service

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"zhihu-go/internal/dao"
	"zhihu-go/internal/model"
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

//关注用户

func FollowUser(db *gorm.DB, followeeID, followerID uint) error {
	return dao.FollowUser(db, followerID, followeeID)
}
