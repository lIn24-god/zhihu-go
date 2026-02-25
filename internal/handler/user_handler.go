package handler

import (
	"net/http"
	"zhihu-go/internal/service"
	"zhihu-go/pkg/response"

	"zhihu-go/internal/dto"

	"github.com/gin-gonic/gin"
)

// UserHandler 结构体：处理用户相关的 HTTP 请求
type UserHandler struct {
	userService service.UserService // 依赖 Service 接口
}

// NewUserHandler 构造函数
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var request = dto.LoginRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	ctx := c.Request.Context()

	user, err := h.userService.RegisterUser(ctx, request.Username, request.Password)
	if err != nil {
		// 使用全局错误处理函数
		HandleError(c, err)
		return
	}

	resp := dto.RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
	}

	response.Success(c, resp)
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var request = dto.LoginRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	ctx := c.Request.Context()
	token, user, err := h.userService.LoginUser(ctx, request.Username, request.Password)
	if err != nil {
		HandleError(c, err)
		return
	}

	resp := dto.LoginResponse{
		Token: token,
		User: dto.UserBrief{
			ID:       user.ID,
			Username: user.Username,
		},
	}

	response.Success(c, resp)
}

// UpdateProfile 更新用户信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// 获取已存储的 user_id（登录时已设置）
	userID, exists := c.Get("user_id")
	if !exists {
		// 未认证，直接返回401，无需经过业务错误处理
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// 类型断言，确保 userID 是 uint 类型
	uintUserID, ok := userID.(uint)
	if !ok {
		// 类型断言失败，说明中间件设置的数据类型不正确，属于系统内部错误
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	var request dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	ctx := c.Request.Context()

	updatedUser, err := h.userService.UpdateProfile(ctx, uintUserID, &request)
	if err != nil {
		HandleError(c, err)
		return
	}

	resp := dto.UpdateProfileResponse{
		ID:       updatedUser.ID,
		Username: updatedUser.Username,
		Email:    updatedUser.Email,
		Bio:      updatedUser.Bio,
	}

	response.Success(c, resp)
}

// MuteUser 禁言/解禁用户
func (h *UserHandler) MuteUser(c *gin.Context) {
	var req dto.MuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	ctx := c.Request.Context()

	if err := h.userService.MuteUser(ctx, req.UserID, req.Hours); err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "User muted successfully"})
}
