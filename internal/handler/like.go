package handler

import (
	"zhihu-go/internal/service"

	"zhihu-go/internal/dto"

	"net/http"

	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
)

func CreateLike(c *gin.Context) {
	var request dto.LikeRequest

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

	if err := service.CreateLike(db, request, uintUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "create like successfully"})
}
