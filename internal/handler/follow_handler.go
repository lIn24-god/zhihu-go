package handler

import (
	"net/http"
	"zhihu-go/internal/service"
	"zhihu-go/pkg/response"

	"zhihu-go/internal/dto"

	"github.com/gin-gonic/gin"
)

// FollowHandler 结构体定义
type FollowHandler struct {
	followService service.FollowService
}

// NewFollowHandler 构造函数
func NewFollowHandler(followService service.FollowService) *FollowHandler {
	return &FollowHandler{followService: followService}
}

// Follow 关注
func (h *FollowHandler) Follow(c *gin.Context) {
	var request = dto.FollowRequest{}

	// 获取已存储的 user_id（登录时已设置）
	followerID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// 类型断言，确保 followerID 是 uint 类型
	followerIDUint, ok := followerID.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	ctx := c.Request.Context()

	err := h.followService.FollowUser(ctx, request.FolloweeID, followerIDUint)
	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Followed successfully"})
}

// Unfollow 取关
func (h *FollowHandler) Unfollow(c *gin.Context) {
	var request = dto.FollowRequest{}

	// 获取已存储的 user_id（登录时已设置）
	followerID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// 类型断言，确保 followerID 是 uint 类型
	followerIDUint, ok := followerID.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	ctx := c.Request.Context()

	if err := h.followService.UnfollowUser(ctx, request.FolloweeID, followerIDUint); err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "Unfollowed successfully"})
}

// GetFollowers 获取用户粉丝列表
func (h *FollowHandler) GetFollowers(c *gin.Context) {

	// 获取已存储的 user_id（登录时已设置）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// 类型断言，确保 userID 是 uint 类型
	userIDUint, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	ctx := c.Request.Context()

	followers, err := h.followService.GetFollowers(ctx, userIDUint)
	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"followers": followers,
		"total":     len(followers),
	})
}

// GetFollowees 获取用户关注列表
func (h *FollowHandler) GetFollowees(c *gin.Context) {
	// 获取已存储的 user_id（登录时已设置）
	userID, exists := c.Get("user_id")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// 类型断言，确保 userID 是 uint 类型
	userIDUint, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	ctx := c.Request.Context()

	followees, err := h.followService.GetFollowees(ctx, userIDUint)
	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, gin.H{
		"followees": followees,
		"total":     len(followees),
	})
}
