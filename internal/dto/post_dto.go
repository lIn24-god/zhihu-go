package dto

type PostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostResponse struct {
	Title    string `json:"title"`
	AuthorID string `json:"authorID"`
	Content  string `json:"content"`
}
