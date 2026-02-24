package service

import (
	"context"
	"log"
	"zhihu-go/internal/model"

	"golang.org/x/sync/singleflight"

	"zhihu-go/internal/cache"

	"zhihu-go/internal/dao"

	"zhihu-go/internal/dto"
	"zhihu-go/pkg/bloom"
)

// PostService 定义文章相关的数据访问接口
type PostService interface {
	CreatePost(ctx context.Context, userID uint, req *dto.PostRequest) (*model.Post, error)
	GetPost(authorID uint, status string) ([]dto.PostResponse, error)
	SearchPosts(keyword string, page, pageSize int) ([]dto.PostResponse, int64, error)
	SoftDeletePost(postID, userID uint) error
	RestorePost(postID, userID uint) error
	GetUserTrash(userID uint) ([]model.Post, error)
	UpdatePost(ctx context.Context, userID, postID uint, req dto.UpdatePostRequest) error
	GetPostByID(ctx context.Context, postID uint) (*model.Post, error)
}

// 结构体定义
type postService struct {
	postDAO   dao.PostDAO
	postCache cache.PostCache
	bloom     bloom.Filter
	sfGroup   singleflight.Group
}

// NewPostService 构造函数
func NewPostService(postDAO dao.PostDAO, postCache cache.PostCache, bloom bloom.Filter) PostService {
	return &postService{
		postDAO:   postDAO,
		postCache: postCache,
		bloom:     bloom,
	}
}

// CreatePost 创建文章，并将新文章 ID 加入布隆过滤器
func (s *postService) CreatePost(ctx context.Context, userID uint, req *dto.PostRequest) (*model.Post, error) {

	//如果前端未传入status，则默认为draft
	if req.Status != "draft" && req.Status != "published" {
		req.Status = "draft"
	}

	// 构建文章模型（省略细节）
	post := &model.Post{
		AuthorID: userID,
		Title:    req.Title,
		Content:  req.Content,
		Status:   req.Status, // 草稿等
	}

	// 存入数据库
	if err := s.postDAO.CreatePost(ctx, post); err != nil {
		return nil, err
	}

	// 将新文章 ID 加入布隆过滤器
	if err := s.bloom.Add(ctx, "bloom:post", post.ID); err != nil {
		// 布隆过滤器添加失败，记录日志，不影响文章创建
		log.Printf("bloom add error: %v", err)
	}

	// 可以主动缓存
	if err := s.postCache.Set(ctx, post); err != nil {
		log.Printf("cache set error: %v", err)
	}

	return post, nil
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
func (s *postService) UpdatePost(ctx context.Context, userID, postID uint, req dto.UpdatePostRequest) error {
	post, err := s.postDAO.GetPostByID(ctx, postID)
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

// GetPostByID 靠ID获取文章(带缓存和singleFlight保护)
func (s *postService) GetPostByID(ctx context.Context, postID uint) (*model.Post, error) {

	// 布隆过滤器检查
	exists, err := s.bloom.Exists(ctx, "bloom:post", postID)
	if err != nil {
		// 布隆过滤器出错，降级处理：记录日志，继续查询（可能造成穿透）
		log.Printf("bloom exists error: %v", err)
	} else if !exists {
		return nil, ErrPostNotFound
	}

	// 尝试从缓存中获取
	post, err := s.postCache.Get(ctx, postID)
	if err != nil {
		log.Printf("cache get error: %v", err) // 记录日志，继续查询
	}
	if post != nil {
		return post, nil //缓存命中
	}

	//缓存未命中，使用singleflight合并请求，防止击穿
	key := string(rune(postID))
	v, err, _ := s.sfGroup.Do(key, func() (interface{}, error) {
		//这里执行真正的数据库查询
		post, err := s.postDAO.GetPostByID(ctx, postID)
		if err != nil {
			return nil, err
		}

		if post == nil {
			// 数据不存在
			return nil, ErrPostNotFound
		}

		//存入缓存
		if err := s.postCache.Set(ctx, post); err != nil {
			log.Printf("cache set error: %v", err) // 缓存设置失败不影响返回
		}

		return post, nil
	})
	if err != nil {
		return nil, err
	}

	return v.(*model.Post), err
}
