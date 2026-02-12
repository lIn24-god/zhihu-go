package dto

type FollowRequest struct {
	FolloweeID uint `json:"followee_id"`
}

type FollowUserInfo struct {
	Username string `json:"username"`
	ID       uint   `json:"id"`
}
