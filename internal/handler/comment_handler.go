package handler

import (
	"zhihu-go/internal/service"
	"zhihu-go/pkg/response"

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
		response.Error(c, http.StatusUnauthorized, "User not authenticated")
		return
	}

	uintUserID, ok := userID.(uint)
	if !ok {
		response.Error(c, http.StatusInternalServerError, "Invalid user ID format")
		return
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid input")
		return
	}

	ctx := c.Request.Context()

	resp, err := h.commentService.CreateComment(ctx, &request, uintUserID)
	if err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, resp)
}
