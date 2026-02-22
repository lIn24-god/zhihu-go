package service

import (
	"zhihu-go/internal/model"

	"zhihu-go/internal/dao"

	"zhihu-go/internal/dto"
)

// PostService 定义文章相关的数据访问接口
type PostService interface {
	CreatePost(rep *dto.PostRequest, authorID uint) error
	GetPost(authorID uint, status string) ([]dto.PostResponse, error)
	SearchPosts(keyword string, page, pageSize int) ([]dto.PostResponse, int64, error)
	SoftDeletePost(postID, userID uint) error
	RestorePost(postID, userID uint) error
	GetUserTrash(userID uint) ([]model.Post, error)
	UpdatePost(userID, postID uint, req dto.UpdatePostRequest) error
}

// 结构体定义
type postService struct {
	postDAO dao.PostDAO
}

// NewPostService 构造函数
func NewPostService(postDAO dao.PostDAO) PostService { return &postService{postDAO: postDAO} }

// CreatePost 创建文章
func (s *postService) CreatePost(rep *dto.PostRequest, authorID uint) error {
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

	err := s.postDAO.CreatePost(post)

	return err
}

// GetPost 获取文章
func (s *postService) GetPost(authorID uint, status string) ([]dto.PostResponse, error) {
	posts, err := s.postDAO.GetPost(authorID, status)
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
func (s *postService) SearchPosts(keyword string, page, pageSize int) ([]dto.PostResponse, int64, error) {
	posts, total, err := s.postDAO.SearchPost(keyword, page, pageSize)
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
func (s *postService) SoftDeletePost(postID, userID uint) error {

	//检查文章是否属于该用户
	post, err := s.postDAO.GetPostByIDWithDeleted(postID)
	if err != nil {
		return err
	}
	if post.AuthorID != userID {
		return ErrUnauthorized
	}

	return s.postDAO.SoftDeletePost(postID)
}

// RestorePost 恢复已删除的文章
func (s *postService) RestorePost(postID, userID uint) error {

	//检查文章是否属于该用户
	post, err := s.postDAO.GetPostByIDWithDeleted(postID)
	if err != nil {
		return err
	}
	if post.AuthorID != userID {
		return ErrUnauthorized
	}

	return s.postDAO.RestorePost(postID)
}

// GetUserTrash 获取用户回收站里的文章
func (s *postService) GetUserTrash(userID uint) ([]model.Post, error) {
	posts, err := s.postDAO.GetUserDeletedPosts(userID)
	if err != nil {
		return nil, err
	}

	return posts, err
}

// UpdatePost 更新文章
func (s *postService) UpdatePost(userID, postID uint, req dto.UpdatePostRequest) error {
	post, err := s.postDAO.GetPostByID(postID)
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

	return s.postDAO.UpdatePost(post)
}
