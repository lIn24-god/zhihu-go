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

// SoftDeletePost 软删除文章
func SoftDeletePost(db *gorm.DB, postID, userID uint) error {

	//检查文章是否属于该用户
	post, err := dao.GetPostByID(db.Unscoped(), postID)
	if err != nil {
		return err
	}
	if post.AuthorID != userID {
		return ErrUnauthorized
	}

	return dao.SoftDeletePost(db, postID)
}

// RestorePost 恢复已删除的文章
func RestorePost(db *gorm.DB, postID, userID uint) error {

	//检查文章是否属于该用户
	post, err := dao.GetPostByID(db.Unscoped(), postID)
	if err != nil {
		return err
	}
	if post.AuthorID != userID {
		return ErrUnauthorized
	}

	return dao.RestorePost(db, postID)
}

// GetUserTrash 获取用户回收站里的文章
func GetUserTrash(db *gorm.DB, userID uint) ([]model.Post, error) {
	posts, err := dao.GetUserDeletedPosts(db, userID)
	if err != nil {
		return nil, err
	}

	return posts, err
}

// UpdatePost 更新文章
func UpdatePost(db *gorm.DB, userID, postID uint, req dto.UpdatePostRequest) error {
	post, err := dao.GetPostByID(db, postID)
	if err != nil {
		return err
	}

	if userID != post.AuthorID {
		return ErrUnauthorized
	}

	//只更新非空字段
	if req.Status != "" {
		post.Status = req.Status
	}

	if req.Title != "" {
		post.Title = req.Title
	}

	if req.Content != "" {
		post.Content = req.Content
	}

	return dao.UpdatePost(db, post)
}
