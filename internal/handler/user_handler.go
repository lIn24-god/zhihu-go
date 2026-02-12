package handler

import (
	"net/http"
	"zhihu-go/internal/service"
	"zhihu-go/internal/utils"

	"zhihu-go/internal/dto"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

//用户注册

func Register(c *gin.Context) {
	var request = dto.UserRequest{}

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
	var request = dto.UserRequest{}

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
