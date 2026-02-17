package dto

type PostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostResponse struct {
	Title    string `json:"title"`
	AuthorID uint   `json:"authorID"`
	Content  string `json:"content"`
}

type GetPostRequest struct {
	PostID uint `json:"post_id"`
}
