package dao

import (
	"zhihu-go/internal/model"

	"gorm.io/gorm"
)

//创建文章

func CreatePost(db *gorm.DB, post *model.Post) error {
	return db.Create(post).Error
}

//获取用户文章

func GetPostByID(db *gorm.DB, authorID uint) (*model.Post, error) {
	var post model.Post
	err := db.Where("author_id = ?", authorID).First(&post).Error
	return &post, err
}

// SearchPost 使用全文索引搜索文章，并返回文章列表和总数
func SearchPost(db *gorm.DB, keyword string, page, pageSize int) ([]model.Post, int64, error) {
	var posts []model.Post
	var total int64

	offSet := (page - 1) * pageSize

	//用全文索引构建查询
	query := db.Where("MATCH(title, content) AGAINST(? IN NATURAL LANGUAGE MODE)", keyword)
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
