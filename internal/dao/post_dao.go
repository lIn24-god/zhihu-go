package dao

import (
	"zhihu-go/internal/model"

	"gorm.io/gorm"
)

// PostDAO 定义文章数据访问接口
type PostDAO interface {
	CreatePost(post *model.Post) error
	GetPost(authorID uint, status string) ([]model.Post, error)
	GetPostByID(postID uint) (*model.Post, error)
	GetPostByIDWithDeleted(postID uint) (*model.Post, error)
	SearchPost(keyword string, page, pageSize int) ([]model.Post, int64, error)
	SoftDeletePost(postID uint) error
	RestorePost(postID uint) error
	GetUserDeletedPosts(userID uint) ([]model.Post, error)
	UpdatePost(post *model.Post) error
}

// 结构体定义
type postDAO struct {
	db *gorm.DB
}

// NewPostDAO 构造函数
func NewPostDAO(db *gorm.DB) PostDAO { return &postDAO{db: db} }

// CreatePost 创建文章
func (u *postDAO) CreatePost(post *model.Post) error {
	return u.db.Create(post).Error
}

// GetPost 获取用户文章
func (u *postDAO) GetPost(authorID uint, status string) ([]model.Post, error) {
	var posts []model.Post
	err := u.db.Where("author_id = ? AND status = ?", authorID, status).Find(&posts).Error
	return posts, err
}

// GetPostByID 通过postID获取文章
func (u *postDAO) GetPostByID(postID uint) (*model.Post, error) {
	var post model.Post
	err := u.db.Where("id = ?", postID).First(&post).Error
	return &post, err
}

// GetPostByIDWithDeleted 通过postID获取文章(包括已删除的)
func (u *postDAO) GetPostByIDWithDeleted(postID uint) (*model.Post, error) {
	var post model.Post
	err := u.db.Unscoped().First(&post, postID).Error
	return &post, err
}

// SearchPost 使用全文索引搜索文章，并返回文章列表和总数
func (u *postDAO) SearchPost(keyword string, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	offSet := (page - 1) * pageSize

	//用全文索引构建查询
	query := u.db.Where("MATCH(title, content) AGAINST(? IN NATURAL LANGUAGE MODE)", keyword)
	if err := query.Model(&model.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	//查询数据，按相关性排序
	err := query.
		Order(gorm.Expr("MATCH(title, content) AGAINST(?) DESC", keyword)). //用gorm.Expr将表达式包装为一个参数
		Offset(offSet).
		Limit(pageSize).
		Find(&posts).Error
	if err != nil {
		return nil, 0, err
	}

	return posts, total, nil

}

// SoftDeletePost 软删除文章
func (u *postDAO) SoftDeletePost(postID uint) error {
	return u.db.Delete(&model.Post{}, postID).Error
}

// RestorePost 恢复软删除的文章
func (u *postDAO) RestorePost(postID uint) error {
	return u.db.Model(&model.Post{}).Unscoped().Where("id = ?", postID).
		Update("deleted_at", nil).Error
}

// GetUserDeletedPosts 获取用户已删除的文章
func (u *postDAO) GetUserDeletedPosts(userID uint) ([]model.Post, error) {
	var posts []model.Post
	err := u.db.Unscoped().Where("author_id = ? AND deleted_at IS NOT NULL", userID).Find(&posts).Error
	return posts, err
}

// UpdatePost 更新文章信息
func (u *postDAO) UpdatePost(post *model.Post) error {
	return u.db.Save(&post).Error
}
