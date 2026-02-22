package handler

import (
	"errors"
	"zhihu-go/internal/service"

	"zhihu-go/internal/dto"

	"net/http"

	"github.com/gin-gonic/gin"
)

// CommentHandler 结构体定义
type CommentHandler struct {
	commentService service.CommentService
	userService    service.UserService
}

// NewCommentHandler 构造函数
func NewCommentHandler(commentService service.CommentService, userService service.UserService) *CommentHandler {
	return &CommentHandler{
		commentService: commentService,
		userService:    userService,
	}
}

func (h *CommentHandler) CreateComment(c *gin.Context) {
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

	//检查是否被禁言
	if err := h.userService.CheckMuted(uintUserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	response, err := h.commentService.CreateComment(&request, uintUserID)
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
