package handler

import (
	"net/http"
	"zhihu-go/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"zhihu-go/internal/dto"
)

func CreatePost(c *gin.Context) {
	var request dto.PostRequest

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	if err := service.CreatePost(db, &request, uintUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create a post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"massage": "Create post successfully"})
}
