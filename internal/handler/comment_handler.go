package handler

import (
	"zhihu-go/internal/service"

	"zhihu-go/internal/dto"

	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	var request dto.CommentRequest

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user id"})
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	response, err := service.CreateComment(db, &request, uintUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "create comment successfully",
		"comment": response,
	})
}
