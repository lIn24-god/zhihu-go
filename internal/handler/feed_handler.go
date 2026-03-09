package handler

import (
	"net/http"
	"strconv"
	"zhihu-go/internal/service"
	"zhihu-go/pkg/response"

	"github.com/gin-gonic/gin"
)

type FeedHandler struct {
	feedService service.FeedService
}

func NewFeedHandler(feedService service.FeedService) *FeedHandler {
	return &FeedHandler{feedService: feedService}
}

// GetFeed 获取当前用户的关注动态
// GET /api/feed?page=1&size=20
func (h *FeedHandler) GetFeed(c *gin.Context) {
	userID := c.GetUint("userID") // 从 JWT 中间件获取
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")

	page, _ := strconv.Atoi(pageStr)
	if page < 1 {
		page = 1
	}
	size, _ := strconv.Atoi(sizeStr)
	if size < 1 || size > 50 {
		size = 20
	}

	ctx := c.Request.Context()
	items, total, err := h.feedService.GetUserFeed(ctx, userID, page, size)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取动态失败")
		return
	}

	response.Success(c, gin.H{
		"list":  items,
		"total": total,
		"page":  page,
		"size":  size,
	})
}
