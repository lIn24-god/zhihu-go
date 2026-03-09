package dto

import "time"

type FeedItem struct {
	ID         uint      `json:"id"`
	PostID     uint      `json:"post_id"`
	AuthorID   uint      `json:"author_id"`
	AuthorName string    `json:"author_name"`
	Title      string    `json:"title"` // 文章标题
	CreatedAt  time.Time `json:"created_at"`
}
