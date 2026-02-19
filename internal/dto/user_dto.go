package dto

//登录，注册请求

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UpdateProfileRequest 修改用户资料请求
type UpdateProfileRequest struct {
	Username string `json:"username" binding:"omitempty,min=2,max=20"`
	Email    string `json:"email" binding:"omitempty,email"`
	Bio      string `json:"bio" binding:"omitempty,max=200"`
}

// UpdateProfileResponse 修改用户资料回应
type UpdateProfileResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Bio      string `json:"bio"`
}

// MuteRequest 禁言请求
type MuteRequest struct {
	UserID uint `json:"user_id" binding:"required"`
	Hours  int  `json:"hours"` // 禁言时长，小于零则为解除禁言
}
