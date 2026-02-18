package dto

type PostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Status  string `json:"status"` //文章状态，前端可选‘draft‘或‘published’
}

type PostResponse struct {
	Title    string `json:"title"`
	AuthorID uint   `json:"authorID"`
	Content  string `json:"content"`
	Status   string `json:"status"`
}
