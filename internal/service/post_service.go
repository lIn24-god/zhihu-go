package service

import (
	"zhihu-go/internal/model"

	"gorm.io/gorm"

	"zhihu-go/internal/dao"

	"zhihu-go/internal/dto"
)

// CreatePost 创建文章
func CreatePost(db *gorm.DB, rep *dto.PostRequest, authorID uint) error {
	post := &model.Post{
		Title:    rep.Title,
		Content:  rep.Content,
		AuthorID: authorID,
	}

	err := dao.CreatePost(db, post)

	return err
}

// GetPostByID 获取文章
func GetPostByID(db *gorm.DB, authorID uint) (*dto.PostResponse, error) {
	post, err := dao.GetPostByID(db, authorID)
	if err != nil {
		return nil, err
	}

	result := &dto.PostResponse{
		Title:    post.Title,
		AuthorID: post.AuthorID,
		Content:  post.Content,
	}

	return result, err
}
