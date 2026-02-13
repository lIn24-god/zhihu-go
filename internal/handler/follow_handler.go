package handler

import (
	"net/http"
	"zhihu-go/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zhihu-go/internal/dto"
)

//关注

func Follow(c *gin.Context) {
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

	db := c.MustGet("db").(*gorm.DB)
	if err := service.FollowUser(db, request.FolloweeID, followerIDUint); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to follow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Followed successfully"})
}

// 获取用户粉丝列表

func GetFollowers(c *gin.Context) {

	// 获取已存储的 user_id（登录时已设置）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 类型断言，确保 userID 是 uint 类型
	userIDuint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	followers, err := service.GetFollowers(db, userIDuint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to get followers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followers": followers,
		"total":     len(followers),
	})
}

//获取用户关注列表

func GetFollowees(c *gin.Context) {
	// 获取已存储的 user_id（登录时已设置）
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 类型断言，确保 userID 是 uint 类型
	userIDuint, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	followees, err := service.GetFollowees(db, userIDuint)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get followees"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"followees": followees,
		"total":     len(followees),
	})
}
