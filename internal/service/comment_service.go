package service

import (
	"zhihu-go/internal/model"

	"gorm.io/gorm"

	"zhihu-go/internal/dto"

	"zhihu-go/internal/dao"
)

// CreateComment 创建评论
func CreateComment(db *gorm.DB, rep *dto.CommentRequest, authorID uint) (*dto.CommentResponse, error) {
	comment := &model.Comment{
		PostID:   rep.PostID,
		AuthorID: authorID,
		Content:  rep.Content,
	}

	err := dao.CreateComment(db, comment)

	response := &dto.CommentResponse{Content: rep.Content}

	return response, err
}
