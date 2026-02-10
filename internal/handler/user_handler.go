package handler

import (
	"net/http"
	"zhihu-go/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//用户注册

func Register(c *gin.Context) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	user, err := service.RegisterUser(db, request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register a user"})
		return
	}

	c.JSON(http.StatusOK, user)
}

//用户登录

func Login(c *gin.Context) {
	var request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	user, err := service.LoginUser(db, request.Username, request.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "failed to login"})
		return
	}

	c.Set("user_id", user.ID)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

//关注

func Follow(c *gin.Context) {
	var request struct {
		FolloweeID uint `json:"followee_id"`
	}

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
