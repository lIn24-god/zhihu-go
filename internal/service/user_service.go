package service

import (
	"zhihu-go/internal/dao"
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
