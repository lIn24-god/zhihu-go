package handler

import (
	"errors"
	"zhihu-go/internal/service"

	"zhihu-go/internal/dto"

	"net/http"

	"github.com/redis/go-redis/v9"
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
	rdb := c.MustGet("rdb").(*redis.Client)

	/*//检查是否被禁言
	if err := service.CheckMuted(db, uintUserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}*/

	response, err := service.CreateComment(db, rdb, &request, uintUserID)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTooFrequent):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create comment"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "create comment successfully",
		"comment": response,
	})
}
