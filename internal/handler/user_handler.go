package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"zhihu-go/internal/service"
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

	c.JSON(http.StatusOK, user)
}

//关注

func Follow(c *gin.Context) {
	var request struct {
		FolloweeID uint `json:"followee_id"`
	}

	followerID := c.MustGet("user_id").(uint)

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)
	if err := service.FollowUser(db, request.FolloweeID, followerID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to follow user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Followed successfully"})
}
