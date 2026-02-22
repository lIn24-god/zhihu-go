package handler

import (
	"errors"
	"net/http"
	"zhihu-go/internal/service"

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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 类型断言，确保 followerID 是 uint 类型
	followerIDUint, ok := followerID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	err := h.followService.FollowUser(request.FolloweeID, followerIDUint)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUserNotFound):
			c.JSON(http.StatusNotFound, gin.H{"errors": err.Error()})
		case errors.Is(err, service.ErrAlreadyFollowed):
			c.JSON(http.StatusConflict, gin.H{"errors": err.Error()})
		case errors.Is(err, service.ErrCannotFollowSelf):
			c.JSON(http.StatusBadRequest, gin.H{"errors": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"errors": "failed to follow"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Followed successfully"})
}

// Unfollow 取关
func (h *FollowHandler) Unfollow(c *gin.Context) {
	var request = dto.FollowRequest{}

	// 获取已存储的 user_id（登录时已设置）
	followerID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 类型断言，确保 followerID 是 uint 类型
	followerIDUint, ok := followerID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if err := h.followService.UnfollowUser(request.FolloweeID, followerIDUint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to unfollow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "unfollowed successfully"})
}

// GetFollowers 获取用户粉丝列表
func (h *FollowHandler) GetFollowers(c *gin.Context) {

	// 获取已存储的 user_id（登录时已设置）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 类型断言，确保 userID 是 uint 类型
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	followers, err := h.followService.GetFollowers(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get followers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followers": followers,
		"total":     len(followers),
	})
}

// GetFollowees 获取用户关注列表
func (h *FollowHandler) GetFollowees(c *gin.Context) {
	// 获取已存储的 user_id（登录时已设置）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 类型断言，确保 userID 是 uint 类型
	userIDUint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	followees, err := h.followService.GetFollowees(userIDUint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get followees"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followees": followees,
		"total":     len(followees),
	})
}
