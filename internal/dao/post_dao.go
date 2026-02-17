package dao

import (
	"zhihu-go/internal/model"

	"gorm.io/gorm"
)

//创建文章

func CreatePost(db *gorm.DB, post *model.Post) error {
	return db.Create(post).Error
}

//获取文章详细信息

func GetPostByID(db *gorm.DB, postID uint) (*model.Post, error) {
	var post model.Post
	err := db.Where("id = ?", postID).First(&post).Error
	return &post, err
}

//获取文章的评论

func GetCommentsByPostID(db *gorm.DB, postID uint) ([]model.Comment, error) {
	var comments []model.Comment
	err := db.Where("post_id = ?", postID).Find(&comments).Error
	return comments, err
}
