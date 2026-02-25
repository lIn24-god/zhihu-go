package handler

import (
	"zhihu-go/internal/service"
	"zhihu-go/pkg/response"

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

	//防刷机制
	if err := h.likeService.CreateLike(ctx, request, uintUserID); err != nil {
		HandleError(c, err)
		return
	}

	response.Success(c, gin.H{"message": "create like successfully"})
}
