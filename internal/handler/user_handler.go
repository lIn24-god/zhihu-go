package handler

import (
	"net/http"
	"zhihu-go/internal/service"
	"zhihu-go/internal/utils"

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
	var request = dto.UserRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.userService.RegisterUser(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register a user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var request = dto.UserRequest{}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	user, err := h.userService.LoginUser(request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to login"})
		return
	}

	// 生成JWT token
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"token":   token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// UpdateProfile 更新用户信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	// 获取已存储的 user_id（登录时已设置）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 类型断言，确保 userID 是 uint 类型
	uintUserID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	var request dto.UpdateProfileRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.userService.UpdateProfile(uintUserID, &request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Profile updated successfully",
		"user":    updatedUser,
	})
}

// MuteUser 禁言/解禁用户
func (h *UserHandler) MuteUser(c *gin.Context) {
	var req dto.MuteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//检查目标用户是否存在
	target, err := h.userService.GetUserProfile(req.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
		return
	}

	//检查目标是否为管理员
	if target.IsAdmin {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no admin"})
		return
	}

	if err := h.userService.MuteUser(target.ID, req.Hours); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to mute"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "mute successfully"})
}
