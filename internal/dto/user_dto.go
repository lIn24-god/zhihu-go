package dto

//登录，注册请求

type UserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
