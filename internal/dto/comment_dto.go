package dto

type CommentRequest struct {
	PostID  uint   `json:"post_id"`
	Content string `json:"content"`
}

type CommentResponse struct {
	Content string `json:"content"`
}
