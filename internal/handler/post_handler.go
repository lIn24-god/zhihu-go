package handler

import (
	"net/http"
	"zhihu-go/internal/service"

	"strconv"

	"github.com/gin-gonic/gin"

	"zhihu-go/internal/dto"
)

// PostHandler 结构体定义
type PostHandler struct {
	postService service.PostService
}

// NewPostHandler 构造函数
func NewPostHandler(postService service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

// CreatePost 文章发布
func (h *PostHandler) CreatePost(c *gin.Context) {
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

	/*//检查是否被禁言
	if err := service.CheckMuted(db, uintUserID); err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}*/

	if err := h.postService.CreatePost(&request, uintUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create a post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"massage": "Create post successfully"})
}

// GetDraft 获取用户草稿
func (h *PostHandler) GetDraft(c *gin.Context) {

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

	result, err := h.postService.GetPost(uintUserID, "draft")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get draft"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": result})
}

// GetPublishedPost 获取用户已发布文章
func (h *PostHandler) GetPublishedPost(c *gin.Context) {

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

	result, err := h.postService.GetPost(uintUserID, "published")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get published post"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": result})
}

// SearchPosts 处理文章搜索请求
func (h *PostHandler) SearchPosts(c *gin.Context) {
	//获取查询参数
	keyword := c.Query("q")
	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing keyword"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("pageSize", "10"))

	//参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	results, total, err := h.postService.SearchPosts(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to search posts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":     results,
		"page":     page,
		"pageSize": pageSize,
		"total":    total,
	})
}

// DeletePost 删除文章
func (h *PostHandler) DeletePost(c *gin.Context) {
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

	//从url中获取postID
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	if err := h.postService.SoftDeletePost(uint(postID), uintUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "delete successfully"})
}

// RestorePost 恢复文章
func (h *PostHandler) RestorePost(c *gin.Context) {

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

	//从url中获取postID
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	if err := h.postService.RestorePost(uint(postID), uintUserID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "restore successfully"})
}

// GetTrash 获取用户回收站中的文章
func (h *PostHandler) GetTrash(c *gin.Context) {
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

	trash, err := h.postService.GetUserTrash(uintUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get trash"})
		return
	}

	var result []dto.PostResponse
	for _, f := range trash {
		result = append(result, dto.PostResponse{
			Title:    f.Title,
			AuthorID: f.AuthorID,
			Content:  f.Content,
		})
	}
	c.JSON(http.StatusOK, gin.H{"trash": result})
}

// UpdatePost 修改文章信息
func (h *PostHandler) UpdatePost(c *gin.Context) {

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

	//从url中获取postID
	postIDStr := c.Param("id")
	postID, err := strconv.ParseUint(postIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid post id"})
		return
	}

	var req dto.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.postService.UpdatePost(uintUserID, uint(postID), req); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "update successfully"})
}
