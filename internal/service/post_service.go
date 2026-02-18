package service

import (
	"zhihu-go/internal/model"

	"gorm.io/gorm"

	"zhihu-go/internal/dao"

	"zhihu-go/internal/dto"
)

// CreatePost 创建文章
func CreatePost(db *gorm.DB, rep *dto.PostRequest, authorID uint) error {
	//如果前端未传入status，则默认为draft
	if rep.Status != "draft" && rep.Status != "published" {
		rep.Status = "draft"
	}

	post := &model.Post{
		Title:    rep.Title,
		Content:  rep.Content,
		AuthorID: authorID,
		Status:   rep.Status,
	}

	err := dao.CreatePost(db, post)

	return err
}

// GetPost 获取文章
func GetPost(db *gorm.DB, authorID uint, status string) ([]dto.PostResponse, error) {
	posts, err := dao.GetPost(db, authorID, status)
	if err != nil {
		return nil, err
	}

	var result []dto.PostResponse
	for _, f := range posts {
		result = append(result, dto.PostResponse{
			Title:    f.Title,
			AuthorID: f.AuthorID,
			Content:  f.Content,
			Status:   f.Status,
		})
	}

	return result, err
}

// SearchPosts 搜索文章
func SearchPosts(db *gorm.DB, keyword string, page, pageSize int) ([]dto.PostResponse, int64, error) {
	posts, total, err := dao.SearchPost(db, keyword, page, pageSize)
	if err != nil {
		return nil, 0, err
	}

	result := make([]dto.PostResponse, len(posts))
	for i, p := range posts {
		result[i] = dto.PostResponse{
			Title:    p.Title,
			AuthorID: p.AuthorID,
			Content:  p.Content,
		}
	}

	return result, total, nil
}
