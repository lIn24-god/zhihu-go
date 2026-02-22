package handler

import (
	"errors"
	"zhihu-go/internal/service"

	"zhihu-go/internal/dto"

	"net/http"

	"github.com/gin-gonic/gin"
)

// LikeHandler 结构体定义
type LikeHandler struct {
	likeService service.LikeService
}

// NewLikeHandler 构造函数
func NewLikeHandler(likeService service.LikeService) *LikeHandler {
	return &LikeHandler{likeService: likeService}
}

func (h *LikeHandler) CreateLike(c *gin.Context) {
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

	//防刷机制
	if err := h.likeService.CreateLike(request, uintUserID); err != nil {
		switch {
		case errors.Is(err, service.ErrTooFrequent):
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to like"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "create like successfully"})
}
